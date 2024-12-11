package main

import (
	"log"
	"math"
)
type stones map[int]int

func (stones *stones) clone() stones {
    cloned := make(map[int]int)

    for key, value := range *stones {
        cloned[key] = value
    }

    return cloned
}

func NewStones() stones {
	return make(map[int]int)
}

func (stones *stones) add(stone int, count int) {
	(*stones)[stone] += count
}

func (stones *stones) remove(stone int) {
	if _, exists := (*stones)[stone]; !exists {
		log.Fatal("Removing non-existing stone ", stone)
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

func iterate(stones stones) stones {
	newStones := NewStones()

	for num, count := range stones {
		if count == 0 {
			continue
		} 
		if num == 0 {
			newStones.add(1, count)
		} else {
			digits := int(math.Log10(float64(num))) + 1
			if digits % 2 == 0 {
				first := num / int(math.Pow10(digits/2))
				second := num % int(math.Pow10(digits/2))
				newStones.add(first, count)
				newStones.add(second, count)
			} else {
				newStones.add(num * 2024, count)
			}
		}
	}
	return newStones
}

func (stones *stones) print() {
	log.Println(*stones)
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
	stones = make(map[int]int)
	for _, v := range []int{9694820, 93, 54276, 1304, 314, 664481, 0, 4} {
	//for _, v := range []int{125,17} {
		stones.add(v, 1)
	}
	stones.print()
	for i := 0; i < 75 ; i++ {
		stones = iterate(stones)
		log.Println("iter", i, "len", len(stones))
		//stones.print()
	}
	log.Println(stones.count())
}