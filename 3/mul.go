package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)


func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}
	s := string(data)

	re := regexp.MustCompile(`(do\(\)|don't\(\)|mul\((\d+), *(\d+)\))`)
	matches := re.FindAllStringSubmatch(s, -1)

	enabled := true
	var sum int64
	for _, match := range matches {
		log.Println(match[0])
		if match[0] == "do()" {
			log.Println("enabling")

			enabled = true
		} else if match[0] == "don't()" {
			log.Println("disabling")
			enabled = false
		} else if len(match) == 4 && enabled {
			a, _ := strconv.Atoi(match[2])
			b, _ := strconv.Atoi(match[3])
			sum += int64(a * b)
			log.Println("adding", a, b, a*b)
		} else {
			log.Println("match", len(match), match)

		}
	}
	log.Println(sum)
}
