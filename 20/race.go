package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type point struct {
	x int
	y int
}

type direction struct {
	x int
	y int
}

var directions = []direction{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

func (self direction) turnLeft() direction {
	return direction{self.y, -self.x}
}

func (self direction) turnRight() direction {
	return direction{-self.y, self.x}
}

func (self point) pos(width int) int {
	return self.x + self.y*width
}

func (self point) add(dir direction) point {
	return point{self.x + dir.x, self.y + dir.y}
}

func posToPoint(pos int, width int) point {
	return point{pos % width, pos / width}
}

func readMap(filename string) (field, point, point) {
	fieldData := make([]rune, 0)
	width := 0
	height := 0

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}

	start := point{-1, -1}
	end := point{-1, -1}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if width == 0 {
			width = len(line)
		}

		if len(line) > 0 && line[0] == '#' {
			if i := strings.IndexRune(line, 'S'); i != -1 {
				start = point{i, height}
			}
			if i := strings.IndexRune(line, 'E'); i != -1 {
				end = point{i, height}
			}
			fieldData = append(fieldData, []rune(line)...)

			height++
		}
	}
	picoSecs := make([]int, len(fieldData))
	for i := range picoSecs {
		picoSecs[i] = -1
	}

	return field{fieldData, picoSecs, width, height}, start, end
}

func printMap(field field) {
	for pos := 0; pos < field.width*field.height; pos += field.width {
		fmt.Println(string(field.data[pos : pos+field.width]))
	}
}

type field struct {
	data     []rune
	picoSecs []int
	width    int
	height   int
}

func (field field) get(pos point) rune {
	return field.data[pos.x+pos.y*field.width]
}
func (field field) getWithDir(pos point, dir direction) rune {
	pos = point{pos.x + dir.x, pos.y + dir.y}
	return field.data[pos.x+pos.y*field.width]
}

func cloneMap(original map[point]bool) map[point]bool {
	clone := make(map[point]bool)
	for key, value := range original {
		clone[key] = value
	}
	return clone
}

type startEnd struct {
	start point
	end point
}

func recordPath(field field, start point, end point) {
	pos := start
	field.picoSecs[pos.pos(field.width)] = 0
	picoSec := 1
	madeStep := false
	for pos != end {
		for _, dir := range directions {
			newPos := pos.add(dir)

			if field.get(newPos) == '.' || field.get(newPos) == 'E' {
				pos = newPos
				field.data[pos.pos(field.width)] = 'o'
				field.picoSecs[pos.pos(field.width)] = picoSec
				picoSec++
				madeStep = true
				break
			}
		}
		if !madeStep {
			log.Println("No way out!", pos)
			break
		}
		madeStep = false
	}
	log.Println("Path recorded in", picoSec, "picoSecs")
}

func getCheats(field field, minPicoSecs int) int {
	numCheats := 0
	for idx, pico := range field.picoSecs {
		if pico < 0 {
			continue
		}
		pos := posToPoint(idx, field.width)
		for _, dir := range directions {
			wallPos := pos.add(dir)
			if field.get(wallPos) != '#' {
				continue
			}
			newPos := wallPos.add(dir)
			if newPos.pos(field.width) < 0 || newPos.pos(field.width) >= len(field.data) {
				continue
			}
			saved := field.picoSecs[newPos.pos(field.width)] - pico - 2 // 2 picoSecs to move through the wall
			if saved >= minPicoSecs {
				// log.Println("Found cheat at", pos, "to", newPos, "saved", saved, "picoSecs")
				numCheats++
			}
		}
	}
	return numCheats

}

func searchCheatFrom(field field, pos point, startPoint point, usedPicoSecs int, cheats map[startEnd]int, visited map[point]int) {
	// log.Println("Searching from", pos, "to", startPoint, "used", usedPicoSecs)
	if usedPicoSecs > 20 {
		return
	}
	if posVal, exists := visited[pos]; exists && posVal <= usedPicoSecs {
		// log.Println("Already visited", pos, "in", visited[pos], "picoSecs")
		return
	}
	// if pos.x <= 0 || pos.x >= field.width-1 || pos.y <= 0 || pos.y >= field.height-1 {  // don't allow going in outer walls
	// 	return
	// }
	visited[pos] = usedPicoSecs
	
	
	for _, dir := range directions {
		newPos := pos.add(dir)
		if newPos.x < 0 || newPos.x >= field.width || newPos.y < 0 || newPos.y >= field.height {
			continue
		}
		if field.get(newPos) == '#' {
			searchCheatFrom(field, newPos, startPoint, usedPicoSecs+1, cheats, visited)
		} else if usedPicoSecs > 0 && field.picoSecs[newPos.pos(field.width)] > 0 {
			startPicoSecs := field.picoSecs[startPoint.pos(field.width)]
			saved := field.picoSecs[newPos.pos(field.width)] - startPicoSecs - usedPicoSecs
			if saved > cheats[startEnd{startPoint, newPos}] {
				cheats[startEnd{startPoint, newPos}] = saved
				log.Println("Found cheat at", startPoint, "to", newPos, "saved", saved, "picoSecs")
			}
		}
	}
}

const INT_MAX = int(^uint(0) >> 1)

func get20PsCheats(field field, minPicoSecs int) int {
	// go into all directions --> bfs in the walls until cheat found or 20ps reached
	cheats := make(map[startEnd]int)
	
	for idx, pico := range field.picoSecs {
		if pico < 0 {
			continue
		}
		startPoint := posToPoint(idx, field.width)
		visited := make(map[point]int)
		searchCheatFrom(field, startPoint, startPoint, 0, cheats, visited)
	}

	numCheats := 0

	for _, saved := range cheats {
		if saved >= minPicoSecs {
			//log.Println("Found cheat with", saved, "picoSecs", cheat)
			numCheats++
		}
	}

	savedCounts := make(map[int]int)
	for _, saved := range cheats {
		savedCounts[saved]++
	}

	savedList := make([]int, 0, len(savedCounts))
	for saved := range savedCounts {
		savedList = append(savedList, saved)
	}
	sort.Ints(savedList)
	for _, saved := range savedList {
		log.Printf("Saved %d picoSecs: %d cheats\n", saved, savedCounts[saved])
	}

	return numCheats
}

func main() {
	field, start, end := readMap("input2.txt")
	log.Println(start, end)
	recordPath(field, start, end)
	printMap(field)
	log.Println("Part 1:", getCheats(field, 100))
	log.Println("Part 2:", get20PsCheats(field, 50))
}
