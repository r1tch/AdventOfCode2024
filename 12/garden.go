package main

import (
	"bufio"
	"log"
	"os"
)

type garden struct {
	plots  []rune
	width  int
	height int
}

func NewGarden() garden {
	return garden{
		plots:  make([]rune, 0),
		width:  0,
		height: 0,
	}
}

func (garden *garden) addLine(line string) {
	if garden.width == 0 {
		garden.width = len(line)
	} else if garden.width != len(line) {
		log.Fatalln("line width not uniform", garden.width, len(line), line)
	}
	garden.plots = append(garden.plots, []rune(line)...)
	garden.height++
}

func (garden *garden) get(x int, y int) rune {
	return garden.plots[y*garden.width+x]
}

func (garden *garden) set(x int, y int, chr rune) {
	garden.plots[y*garden.width+x] = chr
}

func (garden *garden) print() {
	for row := 0; row < garden.height; row++ {
		log.Println(string(garden.plots[row*garden.width : (row+1)*garden.width]))
	}
}

func (garden garden) withinBounds(x int, y int) bool {
	return x >= 0 && y >= 0 && x < garden.width && y < garden.height
}

func (garden garden) numFences(x int, y int) int {
	sum := 0
	chr := garden.get(x, y)
	if !garden.withinBounds(x-1, y) || chr != garden.get(x-1, y) {
		sum++
	}
	if !garden.withinBounds(x+1, y) || chr != garden.get(x+1, y) {
		sum++
	}
	if !garden.withinBounds(x, y-1) || chr != garden.get(x, y-1) {
		sum++
	}
	if !garden.withinBounds(x, y+1) || chr != garden.get(x, y+1) {
		sum++
	}
	return sum
}

/*
   +-+-+-

     --
    |AA|
     --


	not start if not a continuation
    left: is continuation if upper neighbor same plot and it has left fence
	right: is continuation if bottom neighbor is same plot and it has right fence
	top: is continuation if left is same plot and has top fence
	bottom: continuation if left neighbor is same plot and has bottom fence


*/

func (garden garden) numFencesPart2(x int, y int) int {
	sum := 0
	chr := garden.get(x, y)
	if (!garden.withinBounds(x-1, y) || chr != garden.get(x-1, y)) && // left
		!(garden.withinBounds(x, y-1) && chr == garden.get(x, y-1) &&
			(!garden.withinBounds(x-1, y-1) || chr != garden.get(x-1, y-1))) {
		log.Println("  left fence2", point{x, y}, string(chr))
		sum++
	}
	if (!garden.withinBounds(x+1, y) || chr != garden.get(x+1, y)) && // right
		!(garden.withinBounds(x, y-1) && chr == garden.get(x, y-1) &&
			(!garden.withinBounds(x+1, y-1) || chr != garden.get(x+1, y-1))) {
		log.Println("  right fence2", point{x, y}, string(chr))
		sum++
	}
	if (!garden.withinBounds(x, y+1) || chr != garden.get(x, y+1)) && // below
		!(garden.withinBounds(x-1, y) && chr == garden.get(x-1, y) &&
			(!garden.withinBounds(x-1, y+1) || chr != garden.get(x-1, y+1))) {
		log.Println("  below fence2", point{x, y}, string(chr))
		sum++
	}
	if (!garden.withinBounds(x, y-1) || chr != garden.get(x, y-1)) && // over
		!(garden.withinBounds(x-1, y) && chr == garden.get(x-1, y) &&
			(!garden.withinBounds(x-1, y-1) || chr != garden.get(x-1, y-1))) {
		log.Println("  over fence2", point{x, y}, string(chr))
		sum++
	}
	return sum
}

func readGarden() garden {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("can't open file", err)
	}
	scanner := bufio.NewScanner(file)
	garden := NewGarden()
	for scanner.Scan() {
		line := scanner.Text()
		garden.addLine(line)
	}
	return garden
}

type point struct {
	x int
	y int
}

type region struct {
	plots       []point
	fences      int
	fencesPart2 int
	plotChar    rune
}

func (region region) area() int {
	return len(region.plots)
}

func (region *region) add(plot point, fences int, fencesPart2 int) {
	region.plots = append(region.plots, plot)
	region.fences += fences
	region.fencesPart2 += fencesPart2
}

func (region *region) addRegion(x int, y int, garden garden, visitedMap map[point]bool) {
	if !garden.withinBounds(x, y) {
		return
	}

	toBeCleared := false
	if visitedMap == nil {
		region.plotChar = garden.get(x, y)
		visitedMap = make(map[point]bool)
		log.Println("START ", string(garden.get(x, y)))
		toBeCleared = true
	}
	fences := garden.numFences(x, y)
	fencesPart2 := garden.numFencesPart2(x, y)
	region.add(point{x, y}, fences, fencesPart2)
	log.Println("  ADD fences", fences, point{x, y}, string(garden.get(x, y)))

	visitedMap[point{x, y}] = true
	chr := garden.get(x, y)
	if _, visited := visitedMap[point{x + 1, y}]; !visited && garden.withinBounds(x+1, y) && garden.get(x+1, y) == chr {
		region.addRegion(x+1, y, garden, visitedMap)
	}
	if _, visited := visitedMap[point{x - 1, y}]; !visited && garden.withinBounds(x-1, y) && garden.get(x-1, y) == chr {
		region.addRegion(x-1, y, garden, visitedMap)
	}
	if _, visited := visitedMap[point{x, y + 1}]; !visited && garden.withinBounds(x, y+1) && garden.get(x, y+1) == chr {
		region.addRegion(x, y+1, garden, visitedMap)
	}
	if _, visited := visitedMap[point{x, y - 1}]; !visited && garden.withinBounds(x, y-1) && garden.get(x, y-1) == chr {
		region.addRegion(x, y-1, garden, visitedMap)
	}

	if toBeCleared {
		for _, xy := range region.plots {
			garden.set(xy.x, xy.y, '.')
		}
	}
}

func main() {
	garden := readGarden()
	garden.print()

	regions := make([]region, 0)
	for y := 0; y < garden.height; y++ {
		for x := 0; x < garden.width; x++ {
			plot := garden.get(x, y)
			if plot != '.' {
				var region region
				region.addRegion(x, y, garden, nil)
				regions = append(regions, region)
			}
		}
	}

	sum := 0
	sum2 := 0
	for _, region := range regions {
		log.Println(string(region.plotChar), ":", region.fences, region.area())
		sum += region.fences * region.area()
		sum2 += region.fencesPart2 * region.area()
	}
	log.Println("part1:", sum)
	log.Println("part2:", sum2)
}
