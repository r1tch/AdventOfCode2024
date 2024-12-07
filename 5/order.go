package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var violators = make(map[int][]int)

var rules = make(map[[2]int]bool)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var sum int
	var sumPart2 int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			before, err1 := strconv.Atoi(parts[0])
			after, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				violators[after] = append(violators[after], before)
				rules[[2]int{before, after}] = true
			}
		} else if parts := strings.Split(line, ","); len(parts) > 1 {
			var numbers []int
			for _, part := range parts {
				number, err := strconv.Atoi(part)
				if err == nil {
					numbers = append(numbers, number)
				}
			}
			sum += check(numbers)
			sumPart2 += check2(numbers)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	log.Println(sum)
	log.Println(sumPart2)
}

func check2(numbers []int) int {

	if (check(numbers) != 0) {  // ok that's lame :) 
		return 0
	}

	sorted := make([]int, len(numbers))
	copy(sorted, numbers)
	sort.Slice(sorted, func(i, j int) bool {
		first := sorted[i]
		second := sorted[j]
		return rules[[2]int{first, second}]
	})
	return sorted[len(numbers)/2]
}

func check(numbers []int) int {
	cantCome := make(map[int]bool)

	log.Println("checking", numbers)
	for _, v := range numbers {
		log.Println("checking", v)
		if cantCome[v] {
			log.Println(v, "not allowed")
			return 0
		}

		if _, exists := violators[v]; exists {
			for _, v := range violators[v] {
				cantCome[v] = true
				log.Println(v, "not allowed anymore")
			}
		}
	}
	log.Println(numbers, "adding", numbers[len(numbers)/2])
	return numbers[len(numbers)/2]
}
