package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// const WIDTH = 11
// const HEIGHT = 7

const WIDTH = 101
const HEIGHT = 103

type robot struct {
	x  int
	y  int
	vx int
	vy int
}

func mod(x, y int) int {
	return (x%y + y) % y
}
func (robot *robot) move(times int) {
	robot.x = mod(robot.x+robot.vx*times, WIDTH)
	robot.y = mod(robot.y+robot.vy*times, HEIGHT)
}

func (robot *robot) quandrant() int {
	midx := WIDTH / 2
	midy := HEIGHT / 2
	quandrant := 0
	if robot.x == midx || robot.y == midy {
		return -1
	}

	if robot.x > midx {
		quandrant++
	}
	if robot.y > midy {
		quandrant += 2
	}
	return quandrant
}

func readRobots(filename string) []robot {
	robots := make([]robot, 0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var robot robot
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &robot.x, &robot.y, &robot.vx, &robot.vy)
		robots = append(robots, robot)
	}

	return robots
}

func printRobots(field []rune) {
	for pos:=0; pos<WIDTH*HEIGHT;pos+=WIDTH {
		fmt.Println(string(field[pos:pos+WIDTH]))
	}
}

func countGroups(field []rune) int {
	var prev rune
	count := 0
	for _, chr := range field {
		if prev != chr {
			prev = chr
			count++
		}
	}
	return count
}

func main() {
	robots := readRobots("input.txt")
	quadrants := [4]int{0, 0, 0, 0}

	for _, robot := range robots {
		robot.move(100)
		//log.Println(robot)
		quadrant := robot.quandrant()
		if quadrant != -1 {
			quadrants[quadrant]++
		}
	}

	log.Println(quadrants)
	log.Println("safety factor: ", quadrants[0]*quadrants[1]*quadrants[2]*quadrants[3])

	robots = readRobots("input.txt")
	emptyfield := make([]rune, WIDTH*HEIGHT)
	for i:=0;i<WIDTH*HEIGHT;i++ {
		emptyfield[i] = '.'
	}

	for secs := 1; secs < 100000; secs++ {
		field := make([]rune, WIDTH*HEIGHT)
		copy(field, emptyfield)

		for i, robot := range robots {
			robot.move(1)
			robots[i] = robot
			field[robot.x+robot.y*WIDTH] = '#'
		}
		count := countGroups(field)
		if count < 500 {
			fmt.Println("sec ", secs)
			//printRobots(field)
		}
	}
}
