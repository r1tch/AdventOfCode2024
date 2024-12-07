package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func calc(remainingNumbers []int, result int, expectedResult int, part int) bool {
	if len(remainingNumbers) == 1 {
		return remainingNumbers[0]*result == expectedResult ||
			remainingNumbers[0]+result == expectedResult ||
			(part == 2 && concat(result, remainingNumbers[0]) == expectedResult)
	}

	if part == 1 {
		return calc(remainingNumbers[1:], result*remainingNumbers[0], expectedResult, part) ||
			calc(remainingNumbers[1:], result+remainingNumbers[0], expectedResult, part)

	}

	return calc(remainingNumbers[1:], concat(result,remainingNumbers[0]), expectedResult, part) ||
		calc(remainingNumbers[1:], result*remainingNumbers[0], expectedResult, part) ||
		calc(remainingNumbers[1:], result+remainingNumbers[0], expectedResult, part)
}

func concat(x int, y int) int {
	factor := 1
	for y >= factor {
		factor *= 10
	}
	return x*factor + y
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	sumPart1 := 0
	sumPart2 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		value, _ := strconv.Atoi(parts[0])
		numbersStr := strings.Split(parts[1], " ")
		numbers := make([]int, len(numbersStr))
		for i, s := range numbersStr {
			num, _ := strconv.Atoi(s)
			numbers[i] = num
		}
		if calc(numbers, 0, value, 1) {
			sumPart1 += value
		}
		if calc(numbers, 0, value, 2) {
			sumPart2 += value
		}
	}
	log.Println(sumPart1)
	log.Println(sumPart2)

}
