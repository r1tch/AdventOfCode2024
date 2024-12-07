package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var rowLen int
var numRows int

func hasHorizontal(word string, row int, col int, lines []string) bool {
	return col <= rowLen-len(word) &&
		lines[row][col:col+len(word)] == word
}

func hasVertical(word string, row int, col int, lines []string) bool {
	if row > len(lines)-len(word) {
		return false
	}
	for i := 0; i < len(word); i++ {
		if lines[row+i][col] != word[i] {
			return false
		}
	}
	return true
}

func hasDiagonal(word string, row int, col int, lines []string) bool {
	if row > len(lines)-len(word) || col > rowLen-len(word) {
		return false
	}
	for i := 0; i < len(word); i++ {
		if lines[row+i][col+i] != word[i] {
			return false
		}
	}
	return true
}

func hasDiagonal2(word string, row int, col int, lines []string) bool {
	if row < len(word)-1 || col > rowLen-len(word) {
		return false
	}
	for i := 0; i < len(word); i++ {
		if lines[row-i][col+i] != word[i] {
			return false
		}
	}
	return true
}

func xmasNum(row int, col int, lines []string) int {
	sum := 0
	if hasHorizontal("XMAS", row, col, lines) ||
		hasHorizontal("SAMX", row, col, lines) {
		//log.Println("found horizontal XMAS at", row, col)
		sum++
	}
	if hasVertical("XMAS", row, col, lines) ||
		hasVertical("SAMX", row, col, lines) {
		//log.Println("found vertical XMAS at", row, col)
		sum++
	}
	if hasDiagonal("XMAS", row, col, lines) ||
		hasDiagonal("SAMX", row, col, lines) {
		//log.Println("found diagonal XMAS at", row, col)
		sum++
	}
	if hasDiagonal2("XMAS", row, col, lines) ||
		hasDiagonal2("SAMX", row, col, lines) {
		//log.Println("found rev diagonal XMAS at", row, col)
		sum++
	}

	return sum
}

func main() {

	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	sum := 0
	rowLen = len(lines[0])
	numRows = len(lines)
	log.Println("rows", numRows, "cols", rowLen)
	for i := 0; i < numRows; i++ {
		for j := 0; j < rowLen; j++ {
			sum += xmasNum(i, j, lines)
		}
	}
	log.Println(sum)
}
