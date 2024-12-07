package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	listOne := make([]int, 0)
	listTwo := make([]int, 0)
	listOneNumbers := make(map[int]bool)
	listTwoFrequency := make(map[int]int)

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)
		if len(numbers) != 2 {
			log.Fatalf("Invalid line: %s", line)
		}
		num1, err1 := strconv.Atoi(numbers[0])
		num2, err2 := strconv.Atoi(numbers[1])
		if err1 != nil || err2 != nil {
			log.Fatalf("Invalid number in line: %s", line)
		}
		listOne = append(listOne, num1)
		listTwo = append(listTwo, num2)
		listOneNumbers[num1] = true
		listTwoFrequency[num2]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(listOne)
	sort.Ints(listTwo)

	var sum int64
	var sumPart2 int64
	for i, val := range listOne {
		distance := int64(math.Abs(float64(val - listTwo[i])))
		sum += distance
	}

	for num := range listOneNumbers {
		sumPart2 += int64(num) * int64(listTwoFrequency[num])
	}
	log.Println(sum)
	log.Println(sumPart2)
}
