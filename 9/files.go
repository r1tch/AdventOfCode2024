package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func isFile(inputPosition int) bool {
	return inputPosition%2 == 0
}

func doPart1(inputRaw []byte) {
	countPerId := make(map[int]int)
	// Convert the bytes to integers
	var input []int
	for idTimes2, b := range inputRaw {
		if b >= '0' && b <= '9' { // Ensure the byte is a digit
			count := int(b - '0')
			input = append(input, count) // Convert byte to int
			if isFile(idTimes2) {
				countPerId[idTimes2/2] = count
			}
		}
	}

	result := make([]int, 0)

	// example filesRaw: 23331
	// expanded: 00...111...2
	//           ^          ^
	//           front      back
	// we step forward with front, backwards with back
	// position is the position within expanded part
	// id is file ID (0, 1 and 2 above)
	position := 0
	sum := 0
	inputFront := 0
	inputBack := len(input) - 1
	if !isFile(inputBack) {
		inputBack--
	}
	toBackFill := input[inputBack]
	for inputFront < inputBack {
		id := inputFront / 2

		// adding files defined
		for i := input[inputFront]; i > 0; i-- {
			sum += position * id
			fmt.Printf("+%d,", id)
			result = append(result, id)
			position++
		}
		inputFront++

		// filling up spaces from the back
		for i := input[inputFront]; i > 0 && inputFront < inputBack; i-- {
			backId := inputBack / 2
			sum += position * backId
			fmt.Printf("_%d,", backId)
			result = append(result, backId)
			toBackFill--
			if toBackFill == 0 {
				inputBack -= 2
				if inputBack >= 0 {
					toBackFill = input[inputBack]
					//log.Println("backfilling at", inputBack, toBackFill)
				}
			}

			position++
		}
		inputFront++

		// the rest...
		if inputFront == inputBack {
			backId := inputBack / 2
			//log.Println("Filling rest", toBackFill, backId )
			for i := 0; i < toBackFill; i++ {
				sum += position * backId
				fmt.Printf("r%d,", backId)
				result = append(result, backId)
				position++
			}
		}
	}

	countPerIdCheck := make(map[int]int)
	for _, id := range result {
		countPerIdCheck[id]++
	}

	for k, v := range countPerId {
		if countPerIdCheck[k] != v {
			fmt.Println("mismatch, id:", k, "count:", countPerIdCheck[k], "should be:", v)
		}
	}

	fmt.Printf("\nPart 1:%d\n", sum)
}

type idLength struct {
	id     int
	length int
}

func doPart2(inputRaw []byte) {
	var input []int
	blocks := make([]idLength, len(input)/2) // -1 means space
	for index, b := range inputRaw {
		if b >= '0' && b <= '9' { // Ensure the byte is a digit
			count := int(b - '0')
			input = append(input, count) // Convert byte to int
			if isFile(index) {
				blocks = append(blocks, idLength{index / 2, count})
				if count == 0 {
					log.Fatalln("zero len file found at", index)
				}
			} else {
				blocks = append(blocks, idLength{-1, count})
			}
		}
	}

	for filePos := len(blocks) - 1; filePos > 0; filePos-- {
		if blocks[filePos].id == -1 {
			continue
		}
		for spacePos := 0; spacePos < filePos && filePos < len(blocks); spacePos++ {
			/*
			fmt.Print("filePos:", filePos, ",spacePos:", spacePos, " ")
			for i := 0; i < len(blocks); i++ {
				fmt.Print(blocks[i], ",")
			}
			fmt.Println("")
			*/
			
			if blocks[spacePos].id == -1 && blocks[spacePos].length >= blocks[filePos].length {
				spaceRemaining := blocks[spacePos].length - blocks[filePos].length
				blocks[spacePos].id = blocks[filePos].id
				blocks[spacePos].length = blocks[filePos].length
				blocks[filePos].id = -1
				if spaceRemaining > 0 {
					blocks = append(blocks, idLength{})
					copy(blocks[spacePos+1:], blocks[spacePos:])
					blocks[spacePos+1] = idLength{-1, spaceRemaining}
					filePos++
				}
				spacePos = filePos
			}
		}
	}

	sum := 0
	calcPos := 0
	for pos := 0; pos < len(blocks); pos++ {
		// print file, print space, calc sum
		for i := 0; i < blocks[pos].length; i++ {
			if blocks[pos].id != -1 {
				sum += calcPos * blocks[pos].id
			}
			calcPos++
			fmt.Printf("%d,", blocks[pos].id)
		}
	}

	fmt.Printf("\npart 2:%d\n", sum)
}

func main() {

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("can't open file", err)
	}
	defer file.Close()

	inputRaw, _ := io.ReadAll(file)
	doPart1(inputRaw)
	doPart2(inputRaw)

}
