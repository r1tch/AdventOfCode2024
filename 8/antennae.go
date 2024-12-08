package main

import (
	"bufio"
	"log"
	"os"
)

type point struct {
	x, y int
}

func (self point) subtract(other point) point {
	return point{self.x - other.x, self.y - other.y}
}

type gameMap [][]rune

func (gameMap gameMap) print() {
	log.Println(" ")
	for _, runeLine := range gameMap {
		line := string(runeLine)
		log.Println(line)
	}
}

var antinodesPart1 map[point]bool
var antinodesPart2 map[point]bool

// Go, it is a shame that we have to define this... :(
// even bigger shame, no built-in constraint defined for comparison
// even bigger shame, `comparable` only allows for == and !=
func diff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func (gameMap *gameMap) isInside(location point) bool {
	return location.x >= 0 && location.y >= 0 &&
		location.x < width && location.y < len(*gameMap) 
}

func (gameMap *gameMap) addOneAntinode(chr rune, location point, antinodes map[point]bool) {
	if !gameMap.isInside(location) {
		return
	}

	if (*gameMap)[location.y][location.x] != chr {
		antinodes[location] = true
	}
}

func (gameMap *gameMap) addAntinodes(chr rune, l1 point, l2 point) {
	gameMap.addOneAntinode(' ', l1, antinodesPart2)
	gameMap.addOneAntinode(' ', l2, antinodesPart2)
	//log.Println("addAntinodes", string(chr), l1, l2)
	vec1 := l2.subtract(l1)
	antiNode1 := l1.subtract(vec1)
	gameMap.addOneAntinode(chr, antiNode1, antinodesPart1)
	for gameMap.isInside(antiNode1) {
		gameMap.addOneAntinode(' ', antiNode1, antinodesPart2)
		antiNode1 = antiNode1.subtract(vec1)
	}

	vec2 := l1.subtract(l2)
	antiNode2 := l2.subtract(vec2)
	gameMap.addOneAntinode(chr, antiNode2, antinodesPart1)
	for gameMap.isInside(antiNode2) {
		gameMap.addOneAntinode('x', antiNode2, antinodesPart2)
		antiNode2 = antiNode2.subtract(vec2)
	}
}

var width int

func main() {
	antinodesPart1 = make(map[point]bool)
	antinodesPart2 = make(map[point]bool)
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("can't open input:", err)
	}

	// char --> slice of locations
	antennae := make(map[rune][]point)
	gameMap := make(gameMap, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if width == 0 {
			width = len(line)
		} else if width != len(line) {
			log.Fatalln("Error: line length mismatch, had", width, "got", len(line))
		}

		gameMap = append(gameMap, []rune(line))

		y := len(gameMap) - 1
		for x, chr := range line {
			if chr != '.' {
				antennae[chr] = append(antennae[chr], point{x, y})
				//log.Println("Adding", string(chr), "at", x,y)
			}
		}
	}

	for chr, locations := range antennae {
		log.Println("Looking at", string(chr))
		for i := 0; i < len(locations); i++ {
			for j := i + 1; j < len(locations); j++ {
				gameMap.addAntinodes(chr, locations[i], locations[j])
			}
		}
	}
	for antinode, _ := range antinodesPart2 {
		gameMap[antinode.y][antinode.x] = '#'
	}
	gameMap.print()
	log.Println(len(antinodesPart1))
	log.Println(len(antinodesPart2))
}
