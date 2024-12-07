package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type report []int

func isSafePair(prev int, curr int, increasing bool, decreasing bool) bool {
	diff := math.Abs(float64(prev - curr))

	if diff < 1 || diff > 3 {
		return false
	}

	if decreasing && prev < curr {
		return false
	}

	if increasing && prev > curr {
		return false
	}

	return true
}

func isSafe(report report) bool {
	increasing := false
	decreasing := false
	for i := 1; i < len(report); i++ {
		prev := report[i-1]
		curr := report[i]

		if !isSafePair(prev, curr, increasing, decreasing) {
			return false
		}

		if prev < curr {
			increasing = true
		} else if prev > curr {
			decreasing = true
		}
	}
	return true
}

func main() {
	var sum int64
	
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)

		report := make(report, len(numbers))
		for i := 0; i < len(report); i++ {
			report[i], err = strconv.Atoi(numbers[i])
			if err != nil {
				log.Fatalf("Could not convert a number in line: %s", line)
			}
			
		}
		log.Println("got", report)
		if isSafe(report) {
			sum++
			log.Println("safe", report)
		} else {
			found := false
			for i := 0; i < len(report) && !found; i++ {
				shortSlice := make([]int, 0, len(report)-1)
				shortSlice = append(shortSlice, report[:i]...) // Add elements before N
				shortSlice = append(shortSlice, report[i+1:]...) // Add elements after N
			
				if isSafe(shortSlice) {
					sum++
					found = true
					log.Println("found extra", report, shortSlice)
				}
			}
			if found {
				log.Println("safe", report)
			} else {
				log.Println("unsafe", report)
			}
		}
	}

	log.Println(sum)

}