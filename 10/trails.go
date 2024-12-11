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

func cloneMap[K comparable, V any](original map[K]V) map[K]V {
    cloned := make(map[K]V)

    for key, value := range original {
        cloned[key] = value
    }

    return cloned
}

func readGrid() [][]int {
	file, err := os.Open("input.txt")
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

func walk(grid [][]int, x int, y int, seen map[point]bool, reachable map[point]bool, total *int, isPart1 bool) int {
	log.Println("looking at", x, y, grid[y][x], reachable)
	if grid[y][x] == 0 {
		seen = make(map[point]bool)
		reachable = make(map[point]bool)
	} else if seen == nil {
		log.Panicf("Initial call on ", x, y, "which is not a trailhead.")
	}

	seen = cloneMap(seen)

	if _, seenPoint := seen[point{x, y}]; seenPoint {
		return len(reachable)
	}
	seen[point{x, y}] = true

	if x >= 1 && grid[y][x-1]-grid[y][x] == 1 {
		walk(grid, x-1, y, seen, reachable, total, isPart1)
	}
	if x < len(grid[0])-1 && grid[y][x+1]-grid[y][x] == 1 {
		walk(grid, x+1, y, seen, reachable, total, isPart1)

	}
	if y >= 1 && grid[y-1][x]-grid[y][x] == 1 {
		walk(grid, x, y-1, seen, reachable, total, isPart1)

	}
	if y < len(grid)-1 && grid[y+1][x]-grid[y][x] == 1 {
		walk(grid, x, y+1, seen, reachable, total, isPart1)

	}

	if grid[y][x] == 9 {
		// log.Println("reachable:", x, y)
		reachable[point{x, y}] = true
		*total++
		return len(reachable)
	}

	if grid[y][x] == 0 {
		log.Println("trailhead", x, y, "total", *total, "score", len(reachable), reachable)
	}
	return len(reachable)
}

func main() {
	grid := readGrid()

	sum := 0
	total := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == 0 {
				score := walk(grid, x, y, nil, nil, &total, true)
				sum += score
				//log.Println("Trailhead at", x, y, score)
			}
		}
	}
	log.Println("part1:", sum)
	log.Println("part2:", total)
}
