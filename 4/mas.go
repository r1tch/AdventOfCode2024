package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var rowLen int
var numRows int

func hasMas(row int, col int, lines []string) bool {
	log.Println("row", row, "col", col)
	if row > len(lines)-3 || col > len(lines[0])-3 {
		return false
	}

	if lines[row][col] == 'M' &&
		lines[row][col+2] == 'M' &&
		lines[row+1][col+1] == 'A' &&
		lines[row+2][col] == 'S' &&
		lines[row+2][col+2] == 'S' {
		return true
	}
	if lines[row][col] == 'S' &&
		lines[row][col+2] == 'M' &&
		lines[row+1][col+1] == 'A' &&
		lines[row+2][col] == 'S' &&
		lines[row+2][col+2] == 'M' {
		return true
	}
	if lines[row][col] == 'S' &&
		lines[row][col+2] == 'S' &&
		lines[row+1][col+1] == 'A' &&
		lines[row+2][col] == 'M' &&
		lines[row+2][col+2] == 'M' {
		return true
	}
	if lines[row][col] == 'M' &&
		lines[row][col+2] == 'S' &&
		lines[row+1][col+1] == 'A' &&
		lines[row+2][col] == 'M' &&
		lines[row+2][col+2] == 'S' {
		return true
	}
	return false
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
			if hasMas(i, j, lines) {
				sum++
			}
		}
	}
	log.Println(sum)
}
