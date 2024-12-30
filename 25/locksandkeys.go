package main

import (
	"bufio"
	"log"
	"os"
)

func readInput(filename string) ([][]int, [][]int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	current := make([]int, 5)
	readingWhat := ""
	locks := make([][]int, 0)
	keys := make([][]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(readingWhat) == 0 {
			if line == "....." {
				readingWhat = "key"
			} else if line == "#####" {
				readingWhat = "lock"
			}
		}

		if len(line) == 0 {
			switch readingWhat {
			case "key":
				keys = append(keys, current)
			case "lock":
				locks = append(locks, current)
			}
			current = make([]int, 5)
			readingWhat = ""
		}
		for idx, chr := range line {
			if chr == '#' {
				current[idx]++
			}
		}
	}

	switch readingWhat {
	case "key":
		keys = append(keys, current)
	case "lock":
		locks = append(locks, current)
	}

	return locks, keys
}

func main() {
	locks, keys := readInput("input.txt")
	log.Println("Locks:", locks)
	log.Println("Keys:", keys)

	fitting := 0
	// no overlaps mean sums <= 7 (shamelessly using the fact that input is uniform, 5x7 size rectangles)
	for _, key := range keys {
		for _, lock := range locks {
			fits := true
			for i := 0; i < 5; i++ {
				if key[i]+lock[i] > 7 {
					fits = false
					break
				}
			}
			if fits {
				fitting++
				log.Println("Match:", lock, key)
			}
		}
	}

	log.Println("Part 1:", fitting)
}
