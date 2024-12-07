package main

import (
	"bufio"
	"github.com/bits-and-blooms/bitset"
	"log"
	"os"
	"strings"
)

type point struct {
	x, y int
}

func (p point) add(p2 point) point {
	return point{p.x + p2.x, p.y + p2.y}
}

type direction int

const (
	UP direction = iota
	RIGHT
	DOWN
	LEFT
)

type pointDirection struct {
	point
	direction
}

func NewPointDirection(p point, direction point) pointDirection {
	up := point{0, -1}
	down := point{0, 1}
	left := point{-1, 0}
	right := point{1, 0}
	switch direction {
	case up:
		return pointDirection{p, UP}
	case down:
		return pointDirection{p, DOWN}
	case left:
		return pointDirection{p, LEFT}
	case right:
		return pointDirection{p, RIGHT}
	}

	log.Fatalln("bad direction:", direction)
	return pointDirection{p, -1}
}

func (pd pointDirection) encode() uint {
	return uint(pd.x)<<12 + uint(pd.y)<<4 + uint(pd.direction)
}

type loopDetector bitset.BitSet

func (ld *loopDetector) visit(pd pointDirection) {
	(*bitset.BitSet)(ld).Set(pd.encode())
}

func (ld *loopDetector) seen(pd pointDirection) bool {
	return (*bitset.BitSet)(ld).Test(pd.encode())
}

func NewLoopDetector() *loopDetector {
	return (*loopDetector)(bitset.New(8 * 8 * 4))
}

var guard point

func (p point) withinBounds(width int, height int) bool {
	return p.x >= 0 && p.y >= 0 && p.x < width && p.y < height
}

func printMap(gameMap [][]rune) {
	log.Println(" ")
	for _, runeLine := range gameMap {
		line := string(runeLine)
		log.Println(line)
	}
}

var width int

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// load map
	gameMap := make([][]rune, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if width == 0 {
			width = len(line)
		} else if width != len(line) {
			log.Fatalln("Error: line length mismatch, had", width, "got", len(line))
		}

		gameMap = append(gameMap, []rune(line))
		log.Println(len(gameMap), line)

		if x := strings.Index(line, "^"); x != -1 {
			guard = point{x, len(gameMap) - 1}
			log.Println("Found guard", guard)
		}
	}

	const placeObstacles = true
	visited := walker(gameMap, guard, point{0, -1}, placeObstacles)

	log.Println("visited", visited)
	log.Println("effectiveObstacles", effectiveObstacles)
}

var effectiveObstacles int

func cloneMap(gameMap [][]rune) [][]rune {
	clonedMap := make([][]rune, len(gameMap))
	for i := range gameMap {
		clonedMap[i] = make([]rune, len(gameMap[i]))
		copy(clonedMap[i], gameMap[i])
	}

	return clonedMap
}

func testWithObstacle(clonedMap [][]rune, guard point, direction point, obstacle point) {
	//log.Println("testing", guard, direction)

	clonedMap[obstacle.y][obstacle.x] = '#' // place obstacle

	if (walker(clonedMap, guard, direction, false) == -1) {
		//log.Println("obstacle added", obstacle)
		effectiveObstacles++
	}
	clonedMap[obstacle.y][obstacle.x] = '.' // remove obstacle for new tests
}

func walker(gameMap [][]rune, guard point, direction point, placeObstacles bool) int {
	loopDetector := NewLoopDetector()
	visited := 1                    // the starting point
	gameMap[guard.y][guard.x] = 'X' // mark the starting point
	clonedMap := cloneMap(gameMap)

	loopDetector.visit(NewPointDirection(guard, direction))
	for {
		// printMap(gameMap)
		nextField := guard.add(direction)
		if !nextField.withinBounds(width, len(gameMap)) {
			// log.Println("Done at", nextField)
			break
		}
		switch gameMap[nextField.y][nextField.x] {
		case '#':
			// turning right:
			direction = point{-direction.y, direction.x}
		case '.':
			if (placeObstacles) {
				testWithObstacle(clonedMap, guard, direction, nextField)
			}
			guard = nextField
			if loopDetector.seen(NewPointDirection(guard, direction)) {
				//log.Println("seen", guard, direction)
				return -1
			}
			loopDetector.visit(NewPointDirection(guard, direction))
			//log.Println("visited", guard, direction)
			visited++
			gameMap[guard.y][guard.x] = 'X'
		case 'X':
			guard = nextField
			if loopDetector.seen(NewPointDirection(guard, direction)) {
				return -1
			}
			loopDetector.visit(NewPointDirection(guard, direction))
		}
	}
	return visited
}
