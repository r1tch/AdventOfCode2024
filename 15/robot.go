package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

var UP = direction{0, -1}
var DOWN = direction{0, 1}
var LEFT = direction{-1, 0}
var RIGHT = direction{1, 0}

func (self point) add(direction direction) point {
	return point{self.x + direction.x, self.y + direction.y}
}
func (self point) subtract(direction direction) point {
	return point{self.x - direction.x, self.y - direction.y}
}

func (self point) pos(width int) int {
	return self.x + self.y*width
}

func posToPoint(pos int, width int) point {
	return point{pos % width, pos / width}
}

func printMap(field []rune, width int, height int, robot point) {

	field2 := make([]rune, len(field))
	copy(field2, field)
	field2[robot.pos(width)] = '@'
	for pos := 0; pos < width*height; pos += width {
		fmt.Println(string(field2[pos : pos+width]))
	}
}

func readMap(filename string) (int, int, []rune, []rune, point) {
	field := make([]rune, 0)
	moves := make([]rune, 0)
	width := 0
	height := 0

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}

	robot := point{-1, -1}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if width == 0 {
			width = len(line)
		}

		if len(line) > 0 && line[0] == '#' {
			if i := strings.IndexRune(line, '@'); i != -1 {
				robot = point{i, height}
			}
			field = append(field, []rune(line)...)

			height++
		} else {
			moves = append(moves, []rune(line)...)
		}
	}

	return width, height, field, moves, robot
}

func readMapPart2(filename string) (int, int, []rune, []rune, point) {
	field := make([]rune, 0)
	moves := make([]rune, 0)
	width := 0
	height := 0

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}

	robot := point{-1, -1}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if width == 0 {
			width = len(line)
		}

		if len(line) > 0 && line[0] == '#' {
			for i, chr := range line {
				switch chr {
				case '#':
					field = append(field, []rune{'#', '#'}...)
				case 'O':
					field = append(field, []rune{'[', ']'}...)
				case '.':
					field = append(field, []rune{'.', '.'}...)
				case '@':
					field = append(field, []rune{'.', '.'}...)
					robot = point{i * 2, height}
				default:
					log.Fatalln("invalid char on map", string(chr))
				}
			}

			height++
		} else {
			moves = append(moves, []rune(line)...)
		}
	}

	return width * 2, height, field, moves, robot
}

func chrToDirection(chr rune) direction {
	switch chr {

	case '^':
		return UP
	case 'v':
		return DOWN
	case '<':
		return LEFT
	case '>':
		return RIGHT
	}

	log.Fatalln("Invalid direction char:", string(chr))
	return direction{-1, -1}
}

func pushBoxes(field []rune, width int, fromPos point, direction direction) bool {

	pos := fromPos
	for field[pos.pos(width)] == 'O' {
		pos = pos.add(direction)
	}

	if field[pos.pos(width)] == '#' {
		return false
	}

	if field[pos.pos(width)] == '.' {

		field[fromPos.pos(width)] = '.'
		field[pos.pos(width)] = 'O'
		return true
	}

	log.Fatalln("unexpected value on map:", string(field[pos.pos(width)]), pos)
	return false
}

func pushBigBoxVertically(field []rune, width int, fromPos point, direction direction, testRun bool) bool {
	newPos := fromPos.add(direction)
	if field[newPos.pos(width)] == '#' || field[newPos.pos(width)+1] == '#' { // we use this for bounds checking too, yay!
		return false
	}

	pushable := true

	if field[newPos.pos(width)] == '[' && field[newPos.pos(width)+1] == ']' {
		pushable = pushBigBoxVertically(field, width, newPos, direction, testRun)
	}

	if field[newPos.pos(width)] == ']' {
		leftBox := point{newPos.x - 1, newPos.y}
		pushable = pushable && pushBigBoxVertically(field, width, leftBox, direction, testRun)
	}

	if field[newPos.pos(width)+1] == '[' {
		rightBox := point{newPos.x + 1, newPos.y}
		pushable = pushable && pushBigBoxVertically(field, width, rightBox, direction, testRun)
	}
	if !pushable {
		return false
	}

	if !testRun && field[newPos.pos(width)] != '.' {
		log.Fatalln("unexpected char at", newPos, string(field[newPos.pos(width)]))
	}
	if !testRun && field[newPos.pos(width)+1] != '.' {
		log.Fatalln("unexpected char at", newPos, string(field[newPos.pos(width)+1]))
	}
	if !testRun {
		field[newPos.pos(width)] = '['
		field[newPos.pos(width)+1] = ']'
		field[fromPos.pos(width)] = '.'
		field[fromPos.pos(width)+1] = '.'
	}

	return true
}
func pushBoxesPart2Horizontal(field []rune, width int, fromPos point, direction direction) bool {

	pos := fromPos

	// jump over column of boxes
	for field[pos.pos(width)] == '[' || field[pos.pos(width)] == ']' {
		pos = pos.add(direction)
	}

	// hit wall --> no move
	if field[pos.pos(width)] == '#' {
		return false
	}

	// free space --> move to pos!
	if field[pos.pos(width)] == '.' {
		//fmt.Println("__", pos)

		for i := pos; i != fromPos.subtract(direction); {
			field[i.pos(width)] = field[i.subtract(direction).pos(width)]
			//fmt.Println("XX", i, i.subtract(direction))
			//printMap(field, width, 10, point{0,0})

			i = i.subtract(direction)
		}
		field[fromPos.pos(width)] = '.'
		return true
	}

	log.Fatalln("unexpected value on map:", string(field[pos.pos(width)]), pos)
	return false
}

func moveBotAndBoxes(field []rune, width int, robot point, directionChr rune) point {
	direction := chrToDirection(directionChr)
	newPos := robot.add(direction)
	newPosChr := field[newPos.pos(width)]
	if newPosChr == '.' {
		return newPos
	}
	if newPosChr == '#' {
		return robot
	}
	if newPosChr == 'O' && pushBoxes(field, width, newPos, direction) {
		return newPos
	}

	return robot
}

func moveBotAndBoxesPart2(field []rune, width int, robot point, directionChr rune) point {
	direction := chrToDirection(directionChr)
	newPos := robot.add(direction)
	newPosChr := field[newPos.pos(width)]
	if newPosChr == '.' {
		return newPos
	}
	if newPosChr == '#' {
		return robot
	}

	if direction.y == 0 {
		if (newPosChr == '[' || newPosChr == ']') && pushBoxesPart2Horizontal(field, width, newPos, direction) {
			return newPos
		}
		return robot
	}

	if newPosChr == '[' {
		// log.Println("newPos", newPos, string(newPosChr), "robot", robot)
		if pushBigBoxVertically(field, width, newPos, direction, true) {
			pushBigBoxVertically(field, width, newPos, direction, false)
			return newPos
		}
	}
	if newPosChr == ']' {
		// log.Println("newPos", newPos, string(newPosChr))
		if pushBigBoxVertically(field, width, point{newPos.x - 1, newPos.y}, direction, true) {
			pushBigBoxVertically(field, width, point{newPos.x - 1, newPos.y}, direction, false)
			return newPos
		}
	}

	return robot
}

func countScore(pos int, width int, ) int {
	p := posToPoint(pos, width)
	score := p.x + p.y*100
	return score
}

func main() {
	width, height, field, moves, robot := readMap("input.txt")
	//moves = []rune("<<<vvv")
	field[robot.pos(width)] = '.'
	printMap(field, width, height, robot)
	log.Println(string(moves))
	log.Println(robot)
	for _, direction := range moves {
		robot = moveBotAndBoxes(field, width, robot, direction)
		//printMap(field, width, height, robot)
	}

	sum := 0
	for pos, chr := range field {
		if chr == 'O' {
			sum += countScore(pos, width)
		}
	}

	log.Println(sum)

	width, height, field, moves, robot = readMapPart2("input.txt")
	// moves = []rune("<v<^^^^")
	// moves = []rune("<<<^^^^^^^>>>>")
	// moves = []rune("v<^^<<<^^^^^^^>>>>")
	// moves = []rune("^^>>>>>>>>>>vvvvvvvv<vv>>^^>>vvvv<^^<<<^^^^^^^>>>>")

	printMap(field, width, height, robot)
	for _, direction := range moves {
		robot = moveBotAndBoxesPart2(field, width, robot, direction)
		//fmt.Println("RUN")
		//printMap(field, width, height, robot)
	}

	sum = 0
	for pos, chr := range field {
		if chr == '[' {
			sum += countScore(pos, width)
		}
	}
	log.Println(sum)
}
