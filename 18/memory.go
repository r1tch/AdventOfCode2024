package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func (p point) add(other point) point {
	return point{p.x + other.x, p.y + other.y}
}

func (p point) pos() uint {
	return uint(p.y*WIDTH + p.x)
}

type field struct {
	data []rune
	width int
	height int
}

func newField(width int, height int) field {
	return field{make([]rune, width*height), width, height}
}

func (f field) get(point point) rune {
	return f.data[point.y*f.width+point.x]
}

func (f field) set(point point, r rune) {
	f.data[point.y*f.width+point.x] = r
}

func (f field) clone() field {
	data := make([]rune, len(f.data))
	copy(data, f.data)
	return field{data, f.width, f.height}
}

func (f field) print() {
	for y := 0; y < f.height; y++ {
		fmt.Println(string(f.data[y*f.width:(y+1)*f.width]))
	}
}

func readField(filename string, width int, height int, limit int) (field, point) {
	var last point
	field := newField(width, height)
	for i := range field.data {
		field.data[i] = '.'
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && count < limit {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		field.set(point{x, y}, '#')
		last = point{x, y}
		count++
	}
	return field, last
}

type crawler struct {
	pos point
	cost int
}

func cloneIntSlice(slice []int) []int {
	cloned := make([]int, len(slice))
	copy(cloned, slice)
	return cloned
}
func (c crawler) clone() crawler {
	return crawler{c.pos, c.cost}
}

var directions = []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

// const INPUT = "input2.txt"
// const WIDTH = 7
// const HEIGHT = 7
// const MAXBYTES = 12
const INPUT = "input.txt"
const WIDTH = 71
const HEIGHT = 71
const MAXBYTES = 1024

const MAX_INT = int(^uint(0) >> 1) // max int value
func testUntil(until int) (int, point) {
	field, last := readField(INPUT, WIDTH, HEIGHT, until)
	// field.print()

	start := point{0, 0}
	end := point{WIDTH-1, HEIGHT-1}
	crawlers := []crawler{{point{start.x, start.y}, 0}}
	minCosts := make([]int, WIDTH*HEIGHT, WIDTH*HEIGHT)
	for i := range minCosts {
		minCosts[i] = MAX_INT
	}

	successfulCrawlers := make([]crawler, 0)
	for len(crawlers) > 0 {
		// log.Println("Crawlers:", len(crawlers))
		newCrawlers := make([]crawler, 0)
		for _, c := range crawlers {
			if c.pos == end {
				successfulCrawlers = append(successfulCrawlers, c)
				continue
			}

			if c.cost >= minCosts[c.pos.pos()] {
				continue
			}
			minCosts[c.pos.pos()] = c.cost
			
			for _, dir := range directions {
				newPos := c.pos.add(dir)
				if newPos.x < 0 || newPos.x >= WIDTH || newPos.y < 0 || newPos.y >= HEIGHT {
					continue
				}
				if field.get(newPos) == '#' {
					continue
				}
				newCrawlers = append(newCrawlers, crawler{newPos, c.cost+1})
			}
		}
		crawlers = newCrawlers
	}

	// finding minimum cost
	minCost := MAX_INT
	for _, c := range successfulCrawlers {
		if c.cost < minCost {
			minCost = c.cost
		}
	}

	
	return minCost, last
}

func main() {
	log.Println("Part 1")
	minCost, _ := testUntil(MAXBYTES)
	log.Println("Minimum cost is", minCost)
	
	log.Println("Part 2")

	for i:= MAXBYTES; i < 3451; i++ {
		minCost, last := testUntil(i)

		if minCost == MAX_INT {
			log.Println("Last point is", last)
			break
		}
	}

}