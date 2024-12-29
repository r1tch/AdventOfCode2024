package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// surmise: we need to find the shortest solutions on each level, and choose from those
// strategy: prefill keypad movement combinations with shortest paths (eg, A->2 is ^< or <^)

type coord struct {
	x int
	y int
}

type direction struct {
	x int
	y int
}

var charToDir = map[rune]direction{
	'^': {0, -1},
	'v': {0, 1},
	'<': {-1, 0},
	'>': {1, 0},
}

var dirToChar = map[direction]rune{
	{0, -1}: '^',
	{0, 1}:  'v',
	{-1, 0}: '<',
	{1, 0}:  '>',
}

var keypad = map[rune]coord{
	'A': {2, 3},
	'0': {1, 3},
	'1': {0, 2},
	'2': {1, 2},
	'3': {2, 2},
	'4': {0, 1},
	'5': {1, 1},
	'6': {2, 1},
	'7': {0, 0},
	'8': {1, 0},
	'9': {2, 0},
}

var invKeypad = map[coord]rune{
	{2, 3}: 'A',
	{1, 3}: '0',
	{0, 2}: '1',
	{1, 2}: '2',
	{2, 2}: '3',
	{0, 1}: '4',
	{1, 1}: '5',
	{2, 1}: '6',
	{0, 0}: '7',
	{1, 0}: '8',
	{2, 0}: '9',
}

var dirpad = map[rune]coord{
	'^': {1, 0},
	'A': {2, 0},
	'<': {0, 1},
	'v': {1, 1},
	'>': {2, 1},
}

var invDirpad = map[coord]rune{
	{1, 0}: '^',
	{2, 0}: 'A',
	{0, 1}: '<',
	{1, 1}: 'v',
	{2, 1}: '>',
}

func (self coord) add(dir direction) coord {
	return coord{self.x + dir.x, self.y + dir.y}
}

func numRepeatingCharacters(s string) int {
	count := 1
	current := 1
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			current++
		} else {
			current = 1
		}
		if current > count {
			count = current
		}
	}
	return count
}

func findShortestPaths(start coord, end coord, invMap map[coord]rune) []string {
	// log.Println("Finding path from", start, "to", end)
	if start == end {
		return []string{"A"}
	}

	dx := end.x - start.x
	dy := end.y - start.y

	var paths []string
	var generatePaths func(path string, x, y int)
	generatePaths = func(path string, x, y int) {
		// log.Println("Checking", path, x, y, end.x, end.y)
		if x == end.x && y == end.y {
			paths = append(paths, path+"A")
			return
		}
		if _, exists := invMap[coord{x, y}]; !exists { // cannot step on non-keys
			// log.Println("Invalid path", path, x, y)
			return
		}

		if dx > 0 {
			generatePaths(path+">", x+1, y)
		}
		if dx < 0 {
			generatePaths(path+"<", x-1, y)
		}
		if dy > 0 {
			generatePaths(path+"v", x, y+1)
		}
		if dy < 0 {
			generatePaths(path+"^", x, y-1)
		}
	}

	generatePaths("", start.x, start.y)

	//TODO bug --> ^^^vA and v^^^A - does not keep both... but why?
	maxRepeatingCharacters := 0
	filteredPaths := make([]string, 0)
	for _, path := range paths {
		repeatingCharacters := numRepeatingCharacters(path)
		if repeatingCharacters > maxRepeatingCharacters {
			filteredPaths = []string{path}
			maxRepeatingCharacters = repeatingCharacters
		} else if repeatingCharacters == maxRepeatingCharacters {
			filteredPaths = append(filteredPaths, path)
		}
	}
	// log.Println("Filtered paths", paths, filteredPaths)
	return filteredPaths
}

type fromTo struct {
	from rune
	to   rune
}

var shortestKeypadPaths map[fromTo][]string
var shortestDirpadPaths map[fromTo][]string

func fillShortestPaths(keymap map[rune]coord, invMap map[coord]rune) map[fromTo][]string {
	pathMap := make(map[fromTo][]string)
	for start, startCoord := range keymap {
		for end, endCoord := range keymap {
			pathMap[fromTo{start, end}] = findShortestPaths(startCoord, endCoord, invMap)
		}
	}
	return pathMap
}

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// a d g
// b e h
// c f i
// paths[0] = [a,b,c]
// --> adg, adh, adi, aeg, aeh, aei, afg, afh, afi, bdg, bdh, bdi, beg, beh, bei, bfg, bfh, bfi, cdg, cdh, cdi, ceg, ceh, cei, cfg, cfh, cfi
// for each paths[0], call ourselves with paths[1:] and append the results
func permutatePaths(paths [][]string) []string {
	//log.Println("permutating paths:", paths)
	generated := make([]string, 0)
	if len(paths) == 1 {
		return paths[0]
	}

	// ...in each case, we only look at the shortest way
	minLength := len(paths[0][0])
	for _, path := range paths[0] {
		if len(path) < minLength {
			minLength = len(path)
		}
	}

	filteredPaths := make([]string, 0)
	for _, path := range paths[0] {
		if len(path) == minLength {
			filteredPaths = append(filteredPaths, path)
		}
	}
	count := 0
	for _, path := range filteredPaths { // a, b, c from above
		for _, nextPath := range permutatePaths(paths[1:]) { // dg, dh, di, eg, eh, ei, fg, fh, fi from above
			generated = append(generated, path+nextPath)
			count++
			// log.Println("Generated", count, "paths")
		}
	}
	return generated
}

func calculateMinLength(paths [][]string) int {
	sum := 0
	for _, subpaths := range paths {
		minLength := INT_MAX
		for _, path := range subpaths {
			if len(path) < minLength {
				minLength = len(path)
			}
		}
		sum += minLength
	}
	return sum
}

const INT_MAX = int(^uint(0) >> 1)

func solvePart1() {
	tasks := readFile("input.txt")
	tasksValues := make(map[string]int)
	for _, task := range tasks {
		numValue, err := strconv.Atoi(task[:len(task)-1])
		if err != nil {
			log.Println("Error converting task to integer:", err)
		}
		// log.Println(task, numValue)
		tasksValues[task] = numValue
	}

	// just some manual testing
	// log.Println("shortestkeypadpaths", shortestKeypadPaths)
	// log.Println("shortestdirpadpaths", shortestDirpadPaths)
	// log.Println("Shortest keypad path from 5 to 9:", findShortestPaths(keypad['5'], keypad['9'], invKeypad))
	// log.Println("Shortest keypad path from 8 to 9:", shortestKeypadPaths[fromTo{'8', '9'}])
	// log.Println("Shortest keypad path from 7 to A:", shortestKeypadPaths[fromTo{'7', 'A'}])
	// log.Println("Shortest dirpad path from ^ to A:", shortestDirpadPaths[fromTo{'^', 'A'}])
	// log.Println("Shortest dirpad path from ^ to ^:", shortestDirpadPaths[fromTo{'^', '^'}])
	// log.Println("Minimum of paths", calculateMinLength([][]string{[]string{"^^", "vvv", "<"}, []string{"A", "A", "AA"}}))

	sum := 0
	for task, value := range tasksValues {
		minLength := INT_MAX
		log.Println("Solving ", task, " (", value, ")")
		shortestPathsPerChar := make([][]string, 0)
		prevChar := 'A'
		for _, char := range task {
			shortestPathsPerChar = append(shortestPathsPerChar, shortestKeypadPaths[fromTo{prevChar, char}])
			prevChar = char
		}
		shortestKeypadPaths := permutatePaths(shortestPathsPerChar)
		// shortestLength := calculateMinLength(shortestPathsPerChar)
		// log.Println("Shortest keypad paths for ", task, ":", shortestKeypadPaths, shortestLength)

		for _, dirpadPath1 := range shortestKeypadPaths {
			shortestPathsPerChar = make([][]string, 0)
			prevChar = 'A'
			for _, char := range dirpadPath1 {
				shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
				//log.Println("  Finding path from", string(prevChar), "to", string(char), ":", shortestPathsPerChar)
				prevChar = char
			}
			shortestDirpadPaths1 := permutatePaths(shortestPathsPerChar)
			// shortestLength := calculateMinLength(shortestPathsPerChar)

			// log.Println("Shortest dirpad1 paths for ", dirpadPath1, ":", shortestPathsPerChar, shortestDirpadPaths1, shortestLength)

			for _, dirpadPath2 := range shortestDirpadPaths1 {
				shortestPathsPerChar = make([][]string, 0)
				prevChar = 'A'
				for _, char := range dirpadPath2 {
					shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
					//log.Println("    Finding dirpad2 path from", string(prevChar), "to", string(char), ":", shortestPathsPerChar)
					prevChar = char
				}
				//				shortestDirpadPaths2 := permutatePaths(shortestPathsPerChar)
				shortestLength := calculateMinLength(shortestPathsPerChar)
				// log.Println("Shortest paths for (", task, dirpadPath1, ")", dirpadPath2, ":", shortestPathsPerChar, shortestLength)

				if shortestLength < minLength {
					log.Println("Shortest dirpad2 path length: ", shortestLength, task)
					minLength = shortestLength
				}
			}
		}
		sum += minLength * value
	}
	log.Println("Sum of all tasks:", sum)
}

// so this is super interesting and took lot of time.
// The optimal path from dirpad "v" to "A" is ^>A (as opposed to >^A), and this will
// only be found if we go at least 4 levels deep!
func optimizeShortestPaths(pathMap map[fromTo][]string) map[fromTo][]string {
	optimizedPathMap := make(map[fromTo][]string)
	for fromToChr, paths := range pathMap {
		// log.Println("Optimizing", string(fromToChr.from), string(fromToChr.to), paths)
		if len(paths) == 1 {
			optimizedPathMap[fromToChr] = paths
			continue
		}
		minLength := INT_MAX

		for _, dirpadPath1 := range paths {
			shortestPathsPerChar := make([][]string, 0)
			prevChar := 'A'
			for _, char := range dirpadPath1 {
				shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
				//log.Println("  Finding path from", string(prevChar), "to", string(char), ":", shortestPathsPerChar)
				prevChar = char
			}
			shortestDirpadPaths1 := permutatePaths(shortestPathsPerChar)
			// shortestLength := calculateMinLength(shortestPathsPerChar)

			// log.Println("Shortest dirpad1 paths for ", dirpadPath1, ":", shortestPathsPerChar, shortestDirpadPaths1, shortestLength)

			for _, dirpadPath2 := range shortestDirpadPaths1 {
				shortestPathsPerChar = make([][]string, 0)
				prevChar = 'A'
				for _, char := range dirpadPath2 {
					shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
					//log.Println("    Finding dirpad2 path from", string(prevChar), "to", string(char), ":", shortestPathsPerChar)
					prevChar = char
				}
				shortestDirpadPaths2 := permutatePaths(shortestPathsPerChar)
				// shortestLength := calculateMinLength(shortestPathsPerChar)
				// log.Println("Shortest paths for (", task, dirpadPath1, ")", dirpadPath2, ":", shortestPathsPerChar, shortestLength)

				/*
				for _, dirpadPath3 := range shortestDirpadPaths2 {
						shortestPathsPerChar = make([][]string, 0)
						prevChar = 'A'
						for _, char := range dirpadPath3 {
							shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
							prevChar = char
						}
						shortestLength := calculateMinLength(shortestPathsPerChar)
	
						if shortestLength < minLength {
							minLength = shortestLength
							optimizedPathMap[fromToChr] = []string{dirpadPath1}
							va := fromTo{'v', 'A'}
							if fromToChr == va {
								log.Println("vA:", dirpadPath3, shortestLength)
							}
						}
	
				}
				*/
				for _, dirpadPath3 := range shortestDirpadPaths2 {
					shortestPathsPerChar = make([][]string, 0)
					prevChar = 'A'
					for _, char := range dirpadPath3 {
						shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
						prevChar = char
					}
					shortestDirpadPaths3 := permutatePaths(shortestPathsPerChar)
					for _, dirpadPath4 := range shortestDirpadPaths3 {
						shortestPathsPerChar = make([][]string, 0)
						prevChar = 'A'
						for _, char := range dirpadPath4 {
							shortestPathsPerChar = append(shortestPathsPerChar, shortestDirpadPaths[fromTo{prevChar, char}])
							prevChar = char
						}
						shortestLength := calculateMinLength(shortestPathsPerChar)
	
						if shortestLength < minLength {
							minLength = shortestLength
							optimizedPathMap[fromToChr] = []string{dirpadPath1}
							va := fromTo{'v', 'A'}
							if fromToChr == va {
								shortestDirpadPaths4 := permutatePaths(shortestPathsPerChar)
								log.Println("vA:", dirpadPath4, shortestLength, shortestDirpadPaths4)
							}
						}
	
					}
	
				}
			}
		}

	}

	return optimizedPathMap
}

type pathCost struct {
	path string
	cost int
}

// keep the minimal cost path for each character
func optimizeKeypad(pathMap map[fromTo][]string, levels int) map[fromTo]pathCost {
	optimizedPathMap := make(map[fromTo]pathCost)
	for fromToChr, paths := range pathMap {
		log.Println("Optimizing", string(fromToChr.from), string(fromToChr.to), paths)

		minCost := INT_MAX
		var minPath string

		for _, path := range paths {
			cost := 0
			prevChar := 'A'
			for _, char := range path {
				s := calculateSize(false, fromTo{prevChar, char}, levels)
				cost += s
				log.Println("  calculating size for", string(prevChar), string(char), ":", s)
				prevChar = char
			}
			if cost < minCost {
				minCost = cost
				minPath = path
			}
		}
		log.Println("  optimized", string(fromToChr.from), string(fromToChr.to), minPath, minCost)
		optimizedPathMap[fromToChr] = pathCost{minPath, minCost}
	}
	return optimizedPathMap
}

type sizeCacheKey struct {
	fromToChr fromTo
	level     int
}

// to cache the size of the shortest path between two chars and on a certain level
var sizeCache = make(map[sizeCacheKey]int)

func calculateSize(isFirst bool, fromToChr fromTo, level int) int {

	if size, exists := sizeCache[sizeCacheKey{fromToChr, level}]; exists {
		return size
	}

	shortestPaths := shortestDirpadPaths

	if isFirst {
		shortestPaths = shortestKeypadPaths
	}
	shortestPath := shortestPaths[fromToChr][0]

	log.Println("  level", level, "calculating size for", string(fromToChr.from), string(fromToChr.to), ":", shortestPath)
	size := 0
	prevChar := 'A'
	for _, char := range shortestPath {
		if level == 0 {
			size += len(shortestPaths[fromTo{prevChar, char}][0])
		} else {
			s := calculateSize(false, fromTo{prevChar, char}, level-1)
			size += s
		}
		prevChar = char
	}
	log.Println("  level", level, "shortestPath", shortestPath, "size", size)
	sizeCache[sizeCacheKey{fromToChr, level}] = size
	return size
}

func solvePart2Calc2(optimizedKeypadPaths map[fromTo]pathCost) {
	tasks := readFile("input.txt")
	tasksValues := make(map[string]int)
	for _, task := range tasks {
		numValue, err := strconv.Atoi(task[:len(task)-1])
		if err != nil {
			log.Println("Error converting task to integer:", err)
		}
		// log.Println(task, numValue)
		tasksValues[task] = numValue
	}

	sum := 0
	for task, value := range tasksValues {
		log.Println("Solving ", task, " (", value, ")")
		prevChar := 'A'
		size := 0
		for _, char := range task {
			size += optimizedKeypadPaths[fromTo{prevChar, char}].cost
			prevChar = char
		}
		log.Println("Size of", task, ":", size)
		sum += size * value
	}
	log.Println("Sum of all tasks:", sum)

}
const numLevels = 23

func main() {
	// pre-fill shortest paths
	// TODO - dirpad needs no optimization.. (?), keypad does, full 25-level optimiz
	shortestKeypadPaths = fillShortestPaths(keypad, invKeypad)
	shortestDirpadPaths = fillShortestPaths(dirpad, invDirpad)
	shortestDirpadPaths = optimizeShortestPaths(shortestDirpadPaths)
	
	// log.Println("shortestDirpadPaths", shortestDirpadPaths)
	
	for key, path := range shortestDirpadPaths {
		log.Println("shortestDirpadPaths", string(key.from), string(key.to), path)
	}
	
	// for key, path := range shortestKeypadPaths {
	// 	log.Println("shortestKeypadPaths", string(key.from), string(key.to), path)
	// }
	// log.Println("shortestDirpadPaths", shortestDirpadPaths)
	// log.Println("Part 2")
	// log.Println("tests")
	// log.Println("size of A", calculateSize(false, fromTo{'A', 'A'}, 0))
	// log.Println("size of >>", calculateSize(false, fromTo{'>', '>'}, 0))
	// log.Println("size of >A", calculateSize(false, fromTo{'>', 'A'}, 0))
	// log.Println("size of A<", calculateSize(false, fromTo{'A', '<'}, 3))
	// log.Println("size of <^", calculateSize(false, fromTo{'<', '^'}, 3))
	// log.Println(sizeCache)
	optimizedKeypadPaths := optimizeKeypad(shortestKeypadPaths, numLevels)

	// Go has a nice way of print(sort(keys))...
	keys := make([]fromTo, 0, len(optimizedKeypadPaths))
	for key := range optimizedKeypadPaths {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].from == keys[j].from {
			return keys[i].to < keys[j].to
		}
		return keys[i].from < keys[j].from
	})
	for _, key := range keys {
		path := optimizedKeypadPaths[key]
		log.Println("optimizedKeypadPaths", string(key.from), string(key.to), path)
	}

	log.Println("-------Part 1--------")
	solvePart1()
	log.Println("-------Part 2--------")
	solvePart2Calc2(optimizedKeypadPaths)
}
