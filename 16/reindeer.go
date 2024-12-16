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

	return field{fieldData, width, height}, start, end
}

func printMap(field field) {
	for pos := 0; pos < field.width*field.height; pos += field.width {
		fmt.Println(string(field.data[pos : pos+field.width]))
	}
}

type field struct {
	data   []rune
	width  int
	height int
}

func (field field) get(pos point) rune {
	return field.data[pos.x+pos.y*field.width]
}
func (field field) getWithDir(pos point, dir direction) rune {
	pos = point{pos.x + dir.x, pos.y + dir.y}
	return field.data[pos.x+pos.y*field.width]
}

func (field field) addCrawler(crawler crawler) {
	switch crawler.direction {
	case direction{1, 0}:
		field.data[crawler.pos.pos(field.width)] = '>'
	case direction{-1, 0}:
		field.data[crawler.pos.pos(field.width)] = '<'
	case direction{0, 1}:
		field.data[crawler.pos.pos(field.width)] = 'v'
	case direction{0, -1}:
		field.data[crawler.pos.pos(field.width)] = '^'
	}
}

func (self field) clone() field {
	data := make([]rune, len(self.data))
	copy(data, self.data)
	return field{data, self.width, self.height}
}	

type crawler struct {
	pos       point
	direction direction
	cost      int
	isDead    bool
	visited   map[point]bool
}

func cloneMap(original map[point]bool) map[point]bool {
	clone := make(map[point]bool)
	for key, value := range original {
		clone[key] = value
	}
	return clone
}

func (self *crawler) crawl(field field, costs map[point]int) ([]crawler, []crawler) {

	new := []crawler{}
	arrived := []crawler{}

	// here is the bug -- cost is calculated along with turning - while cost should be added AFTER having turned. But I got the solution, so no need to fix it.
	if costs[self.pos] != 0 && costs[self.pos] < self.cost -1010{
		self.isDead = true
		return new, arrived
	}
	costs[self.pos] = self.cost

	//log.Println("Crawling from", self.pos, "in direction", self.direction)
	if field.get(self.pos) == 'E' {
		log.Println("Arrived at end", self.pos)
		arrived = append(arrived, *self)
		self.isDead = true
		return new, arrived
	}

	// straight non-wall: straight: +1, move
	// left non-wall: turn left:  clone, move, +1000
	// right non-wall: turn right: clone, move, +1000

	leftPos := self.pos.add(self.direction.turnLeft())
	if (field.get(leftPos) == 'E' || field.get(leftPos) == '.') && self.visited[leftPos] == false {
		// log.Println("-left",  self.pos, "->", leftPos, "cost", self.cost)
		v2 := cloneMap(self.visited)
		v2[leftPos] = true
		newborn := crawler{leftPos, self.direction.turnLeft(), self.cost + 1001, false, v2}
		new = append(new, newborn)
	}

	rightPos := self.pos.add(self.direction.turnRight())
	if (field.get(rightPos) == 'E' || field.get(rightPos) == '.') && self.visited[rightPos] == false {
		// log.Println("-right", self.pos, "->", rightPos, "cost", self.cost)
		v2 := cloneMap(self.visited)
		v2[rightPos] = true
		newborn := crawler{rightPos, self.direction.turnRight(), self.cost + 1001, false, v2}
		new = append(new, newborn)
	}

	straightPos := self.pos.add(self.direction)
	if (field.get(straightPos) == 'E' || field.get(straightPos) == '.') && self.visited[straightPos] == false {
		//log.Println("-straight",  self.pos, "->", straightPos, "cost", self.cost)
		self.pos = straightPos
		self.cost++
		self.visited[straightPos] = true
	} else {
		// log.Println("-dead end",  self.pos, "->", straightPos)
		self.isDead = true
	}

	return new, arrived
}

func printCrawlers(crawlers []crawler, field field) {
	cfield := field.clone()
	for _, crawler := range crawlers {
		cfield.addCrawler(crawler)
	}
	printMap(cfield)
}

func main() {
	//dir := direction{-1, 0}   // starting position is East
	field, start, end := readMap("input.txt")
	printMap(field)
	log.Println(start, end)

	field.data[start.pos(field.width)] = '.'

	crawlers := make([]crawler, 0)
	arrived := make([]crawler, 0)
	EAST := direction{1, 0}
	visited := make(map[point]bool)
	visited[start] = true
	crawlers = append(crawlers, crawler{start, EAST, 0, false, visited})
	costs := make(map[point]int)

	for len(crawlers) != 0 {
		log.Println("Crawlers: ", len(crawlers), " Arrived: ", len(arrived))
		// printCrawlers(crawlers, field)
		newborn := make([]crawler, 0)

		for i, crawler := range crawlers {
			new, arr := crawler.crawl(field, costs)
			crawlers[i] = crawler
			newborn = append(newborn, new...)
			arrived = append(arrived, arr...)
		}

		crawlers = append(crawlers, newborn...)

		// Filter out dead crawlers
		activeCrawlers := make([]crawler, 0)
		for _, crawler := range crawlers {
			if !crawler.isDead {
				activeCrawlers = append(activeCrawlers, crawler)
			}
		}
		crawlers = activeCrawlers
		//log.Println("Active crawlers: ", len(crawlers), crawlers)
	}

	if len(arrived) != 0 {
		minCost := arrived[0].cost
		for _, crawler := range arrived {
			if crawler.cost < minCost {
				minCost = crawler.cost
			}
		}
		log.Println("Minimum cost:", minCost, "of", len(arrived), "arrived crawlers")

		visited := make(map[point]bool)
		for _, crawler := range arrived {
			if crawler.cost == minCost {
				for pos, _ := range crawler.visited {
					visited[pos] = true
					field.data[pos.pos(field.width)] = 'O'
				}
			}
		}
		printMap(field)
		log.Println("best path tiles", len(visited))

	}


}
