package main

import (
	"log"
	"math"
)

type stones []int

func (stones *stones) add(stone int) {
	if len(*stones) <= stone {
		*stones = append(*stones, make([]int, stone-len(*stones)+1)...)
	}
	(*stones)[stone]++
}

func (stones *stones) remove(stone int) {
	if len(*stones) <= stone {
		log.Fatal("Removing non-existing stone", stone)
	}
	(*stones)[stone]--
}


func (stones *stones) count() int {
	sum := 0
	for _, count := range *stones {
		sum += int(count)
	}

	return sum
}

func (stones *stones) iterate() {
	for num, count := range *stones {
		if count == 0 {
			continue
		} 
		stones.remove(0)
		if num == 0 {
			stones.add(1)
		} else {
			digits := int(math.Log10(float64(num))) + 1
			if digits % 2 == 0 {
				first := num / int(math.Pow10(digits/2))
				second := num % int(math.Pow10(digits/2))
				stones.add(first)
				stones.add(second)				
			} else {
				stones.add(num * 2024)
			}
		}
	}
}

// naive won't work with part 2, will be >16TB... 2,792,105,593,192 stones, each 8 bytes...

func iterateOld(stones []int) []int {

	newStones := make([]int, 0, len(stones)*2)

	for _, num := range stones {
		if num == 0 {
			newStones = append(newStones, 1)
		} else {
			digits := int(math.Log10(float64(num))) + 1
			if digits % 2 == 0 {
				first := num / int(math.Pow10(digits/2))
				second := num % int(math.Pow10(digits/2))
				newStones = append(newStones, first)
				newStones = append(newStones, second)
			} else {
				newStones = append(newStones, num * 2024)
			}
		}
	}
	return newStones
}

func main() {
	//stones := []int{9694820, 93, 54276, 1304, 314, 664481, 0, 4}
	// stones := []int{9694820}
	//stones := []int{125, 17}


	var stones stones
	//for _, v := range []int{9694820, 93, 54276, 1304, 314, 664481, 0, 4} {
	for _, v := range []int{125,17} {
		stones.add(v)
	}
	log.Println("added")
	for i := 0; i < 25 ; i++ {
		stones.iterate()
		log.Println("iter", i, "len", len(stones))
	}
	log.Println(stones.count())
}