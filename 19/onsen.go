package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func readInput(filename string) (map[string]bool, []string) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	scanner := bufio.NewScanner(file)
	towels := make(map[string]bool)
	designs := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(towels) == 0 {
			for _, towel := range strings.Split(line, ", ") {
				towels[towel] = true
			}
			continue
		}

		designs = append(designs, line)
	}
	
	return towels, designs
}

var numPossibilities map[string]int = make(map[string]int)

func getNumPossibilities(towels map[string]bool, design string, minLen int, maxLen int) int {
	if _, exists := numPossibilities[design]; exists {
		return numPossibilities[design]
	}

	// log.Println("Checking", design)

	sum := 0
	if _, exists := towels[design]; exists {
		sum += 1
	}

	for skip := minLen; skip < len(design) && skip <= maxLen; skip++ {
		// log.Println(" ...", design[:skip], design[skip:])
		if _, exists := towels[design[:skip]]; exists {
			sum += getNumPossibilities(towels, design[skip:], minLen, maxLen)
		}
	}
	// log.Println("  ", design, ":", sum)
	numPossibilities[design] = sum
	return sum
}

const MAX_INT = int(^uint(0) >> 1)
func main() {
	towels, designs := readInput("input.txt")
	// log.Println("Towels:", towels, "Designs:", designs)

	minLen := MAX_INT
	maxLen := 0
	for towel, _ := range towels {
		if len(towel) < minLen {
			minLen = len(towel)
		}
		if len(towel) > maxLen {
			maxLen = len(towel)
		}
	}

	possibleDesigns := 0
	sum := 0
	for _, design := range designs {
		poss :=  getNumPossibilities(towels, design, minLen, maxLen) 
		if poss > 0 {
			possibleDesigns++
		}
		sum += poss
		log.Println("Design", design, "has", poss, "possibilities")

	}
	log.Println("Part 1:", possibleDesigns)
	log.Println("Part 2:", sum)
}