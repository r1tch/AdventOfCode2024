package main

import (
	"log"
	"math"
)



func iterate(stones []int) []int {

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
	stones := []int{9694820}
	//stones := []int{125, 17}

	for i := 0; i < 75 ; i++ {
		oldLen := len(stones)
		stones = iterate(stones)
		//log.Println(stones)
		log.Println("iter", i, len(stones), float64(len(stones)) / float64(oldLen))
	}
	// will be >16TB... 2,792,105,593,192 stones, each 8 bytes...
	log.Println(len(stones))
}