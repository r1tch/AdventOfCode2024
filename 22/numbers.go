package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

/* Chain of thought: 
   - input.txt has 1557 numbers -> sequence of 3000000 numbers in total
   - 4 changes, we need to store a gain to each
   - 20 values / change --> 20*20*20*20 = 160000 possible values
   - we can store the gain for each value in an array... 
   - we just need to calc all "secret numbers"

   Checking the example
   -2,1,-1,3 --> encoded: 
   7 * 19*19*19 + 10 * 19*19 + 8 * 19 + 12 = 51787
*/

func readNumbers(filename string) []int {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	numbers := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		num, _ := strconv.Atoi(line)
		numbers = append(numbers, num)
	}
	return numbers
}

func nextSecret(num int) int {
	num = num ^ (num << 6) % 0x1000000
	num = (num ^ (num >> 5)) % 0x1000000
	num = (num ^ (num * 2048)) % 0x1000000
	return num
}

func convertChanges(maxChanges int) string {
	strChanges := "]"
	for i := 0; i < 4; i++ {
		change := maxChanges % 19
		maxChanges /= 19
		strChanges = strconv.Itoa(change - 9) + " " + strChanges
	}

	return "[ " + strChanges
}

func main() {
	numbers := readNumbers("input.txt")

	bananas := make([]int, 19*19*19*19)
	part1Sum := 0
	for _, num := range numbers {

		log.Println("Num:", num)
		n := num
		prevLastDigit := -1
		fourChanges := 0 // four changes to be encoded here
		// we need to mark encountered four-changes, as we can only sell the information once
		encounteredFourChanges := make([]int, 19*19*19*19)
		for i := 0; i < 2000; i++ {
			lastDigit := n % 10
			change := lastDigit - prevLastDigit
			fourChanges = (fourChanges*19 + change + 9) % (19*19*19*19)
			if i > 3 && encounteredFourChanges[fourChanges] == 0 {
				bananas[fourChanges] += lastDigit
				encounteredFourChanges[fourChanges] = 1
			}
			prevLastDigit = lastDigit
			n = nextSecret(n)
		}
		// log.Println(num, ":", n)
		part1Sum += n
	}

	maxBananas := 0
	maxChanges := 0
	for changes, bananas := range bananas {
		if bananas > maxBananas {
			maxBananas = bananas
			maxChanges = changes
		}
	}
	log.Println("Part 1:", part1Sum)
	log.Println("Part 2:", maxBananas)

	log.Println("MaxChanges:", maxChanges, convertChanges(maxChanges))
}
