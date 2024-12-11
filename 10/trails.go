package main

import (
	"bufio"
	"log"
	"os"
)

type point struct {
	x int
	y int
}

func readGrid() [][]int {
	file, err := os.Open("input2.txt")
	defer file.Close()
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	grid := make([][]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		intLine := make([]int, 0, len(line))
		for _, chr := range line {
			if chr >= '0' && chr <= '9' {
				intLine = append(intLine, int(chr-'0'))
			} else {
				intLine = append(intLine, -1)

			}
		}
		grid = append(grid, intLine)
	}
	return grid
}

func walk(grid [][]int, x int, y int, seen map[point]bool, score int) int {
	log.Println("looking at", x, y, grid[y][x], score)
	if grid[y][x] == 0 {
		seen = make(map[point]bool)
	} else if seen == nil {
		log.Panicf("Initial call on ", x, y, "which is not a trailhead.")
	}
	if _, seenPoint := seen[point{x, y}]; seenPoint {
		return score
	}
	seen[point{x, y}] = true

	if x >= 1 && grid[y][x-1] >= grid[y][x] {
		score = walk(grid, x-1, y, seen, score)
	}
	if x < len(grid[0])-1 && grid[y][x+1] >= grid[y][x] {
		score = walk(grid, x+1, y, seen, score)

	}
	if y >= 1 && grid[y-1][x] >= grid[y][x] {
		score = walk(grid, x, y-1, seen, score)

	}
	if y < len(grid)-1 && grid[y+1][x] >= grid[y][x] {
		score = walk(grid, x, y+1, seen, score)

	}

	if grid[y][x] == 9 {
		return score + 1
	}

	if grid[y][x] == 0 {
		log.Println("trailhead", x, y, "score", score)

	}
	return score
}

func main() {
	grid := readGrid()

	sum := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == 0 {
				log.Println("Trailhead at", x, y)
				sum += walk(grid, x, y, nil, 0)
			}
		}
	}
	log.Println("part1:", sum)
}
