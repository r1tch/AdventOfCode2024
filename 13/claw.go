package main

// naive search for part 2:
// 2024/12/14 16:13:19 {{55 84} {64 29} {10000000006049 10000000005045}}  a: -1  b: -1  cost: -1
// 2024/12/14 16:15:44 {{52 11} {67 67} {10000000002826 10000000001760}}  a: -1  b: -1  cost: -1
// 2024/12/14 16:17:04 {{78 56} {22 50} {10000000001222 10000000001322}}  a: -1  b: -1  cost: -1
// 2024/12/14 16:27:14 {{13 49} {55 12} {10000000007975 10000000007671}}  a: 169358015095  b: 141788105668  cost: 649862150953

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type point struct {
	x int
	y int
}

type count struct {
	a int
	b int
}

func (self *point) add(other point) {
	self.x += other.x
	self.y += other.y
}

type machine struct {
	buttonA point
	buttonB point
	prize   point
}

func readMachines(filename string, part2 bool) []machine {
	machines := make([]machine, 0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open ", filename, ": ", err)
	}
	scanner := bufio.NewScanner(file)
	var machine machine
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Button A: ") {
			fmt.Sscanf(line, "Button A: X+%d, Y+%d", &machine.buttonA.x, &machine.buttonA.y)
		} else if strings.HasPrefix(line, "Button B: ") {
			fmt.Sscanf(line, "Button B: X+%d, Y+%d", &machine.buttonB.x, &machine.buttonB.y)
		} else if strings.HasPrefix(line, "Prize: ") {
			fmt.Sscanf(line, "Prize: X=%d, Y=%d", &machine.prize.x, &machine.prize.y)
			if part2 {
				machine.prize.x += 10000000000000
				machine.prize.y += 10000000000000
			}
			machines = append(machines, machine)
		}
	}

	return machines
}

const COST_A = 3
const COST_B = 1

func getCostForNaive(machine machine) int {
	minCost := -1
	bestA := -1
	bestB := -1
	// A is more expensive, so try to use A the least amount of time
	for acount := 0; acount*machine.buttonA.x <= machine.prize.x; acount++ {
		bcount := (machine.prize.x - acount*machine.buttonA.x) / machine.buttonB.x
		if acount*machine.buttonA.x+bcount*machine.buttonB.x == machine.prize.x &&
			acount*machine.buttonA.y+bcount*machine.buttonB.y == machine.prize.y {
			cost := acount*COST_A + bcount*COST_B
			if minCost == -1 || cost < minCost {
				minCost = cost
				bestA = acount
				bestB = bcount
				//log.Println("tmp ", machine, cost, bestA, bestB)
			}
		}
	}
	log.Println(machine, " a:", bestA, " b:", bestB, " cost:", minCost)
	return minCost
}


func getSideSolution(prize int, astep int, bstep int) (int, int) {
	for acount := 0; acount < astep*bstep; acount++ {
		bcount := (prize - acount*astep) / bstep
		if acount*astep+bcount*bstep == prize {
			return acount, bcount
		}
	}

	return -1, -1
}

func getSideSolutions(prize int, astep int, bstep int, max int) []count {
	counts := make([]count, 0)

	for acount := 0; acount < astep*bstep*1111 && acount < prize && max > 0; acount++ {
		bcount := (prize - acount*astep) / bstep
		if acount*astep+bcount*bstep == prize {
			counts = append(counts, count{acount, bcount})
			max--
			// log.Println("found counts ", acount, bcount)
		}
	}

	return counts
}

/*
func getCostForOptimizedOldTry(machine machine) int {
	//log.Println("looking at ", machine)
	countsx := getSideSolutions(machine.prize.x, machine.buttonA.x, machine.buttonB.x)
	if len(countsx) == 0 {
		log.Println("No X solution for ", machine)
		return -1
	}
	countsy := getSideSolutions(machine.prize.y, machine.buttonA.y, machine.buttonB.y)
	if len(countsy) == 0 {
		log.Println("No Y solution for ", machine)
		return -1
	}

	//log.Println(countsx)

	return -1
}
	*/

func getCostForOptimized(machine machine) int {
	astepx := machine.buttonA.x
	bstepx := machine.buttonB.x
	prizex := machine.prize.x
	astepy := machine.buttonA.y
	bstepy := machine.buttonB.y
	prizey := machine.prize.y

	minCost := -1
	bestA := -1
	bestB := -1
	//log.Println("looking at ", machine)
	for acount := 0; acount < astepx*bstepx*1111111 && acount < prizex; acount++ { // enough to do it this far?
		bcount := (prizex - acount*astepx) / bstepx
		if acount*astepx+bcount*bstepx == prizex && acount*astepy+bcount*bstepy == prizey {
			cost := acount*COST_A + bcount*COST_B
			if minCost == -1 || cost < minCost {
				minCost = cost
				bestA = acount
				bestB = bcount
				//log.Println("tmp ", machine, cost, bestA, bestB)
			}
//			log.Println("found ", machine, "a:", acount, "b:", bcount, "cost:", acount*3+bcount)
//			return acount*3 + bcount
		}
	}

//	log.Println("No solution for", machine)
//	return -1
	log.Println(machine, " a:", bestA, " b:", bestB, " cost:", minCost)
	return minCost
}

func getCostForOptimized2(machine machine) int {
	astepx := machine.buttonA.x
	bstepx := machine.buttonB.x
	prizex := machine.prize.x
	astepy := machine.buttonA.y
	bstepy := machine.buttonB.y
	prizey := machine.prize.y

	// minCost := -1
	// bestA := -1
	// bestB := -1

	//log.Println("looking at", machine)
	two_x_solutions := getSideSolutions(prizex, astepx, bstepx, 4)
	two_y_solutions := getSideSolutions(prizey, astepy, bstepy, 4)
	if len(two_x_solutions) < 2 || len(two_y_solutions) < 2 {
		return 0
	}
	ax1 := two_x_solutions[0].a
	ax2 := two_x_solutions[1].a
	bx1 := two_x_solutions[0].b
	bx2 := two_x_solutions[1].b
	ay1 := two_y_solutions[0].a
	ay2 := two_y_solutions[1].a
	by1 := two_y_solutions[0].b
	by2 := two_y_solutions[1].b

	// example solutions for separate sides (axes), number of a and b button presses to get 
	// to the prize. a1, a2,... and b1, b2,... are button presses,
	// x [{a1 b1} {a2 b2} ...]
	// first four:
	// x [{9 312500000043} {41 312499999982} {73 312499999921} {105 312499999860}]
    // y [{19 181818181855} {74 181818181821} {129 181818181787} {184 181818181753}]
	// Let's use i as iteration (that is, first iteration for x: i = 1,  a1 = 9)
	// 
	// On above example, the x-axis solutions follow the formula:
	// ax = 9+i*32, bx = i*-61 + 312500000043
	// ...generalized:
	// ax = ax1 + i*(ax2-ax1)
	// bx = bx1 + i*(bx2-bx1)
	// 
	// we must use the same solution for a and b, that is, use the same iteration:
	// i = (ax-9)/32 = (bx-312500000043) / (-61)
	// i = (ax-ax1)/(ax2-ax1) = (bx-bx1) / (bx2-bx1)
	// 
	// y solutions similarly: 
	// (using j to denote the y-axis iteration, could be iy too, but j is shorter :)
	// ay = 19+j*55, by = j*-31 + 181818181855
	// ay = ay1+j*(ay2-ay1)
	// by = by1+j*(by2-by1)
	// 
	// a and b are the number of button presses, must be equal for x and y --> ax == ay, bx == by
	// ax = ay
	// 9+i*32 = 19+j*55
	// i = (10 + j*55)/32
	// ...with formula:
    // ax1+i*(ax2-ax1) = ay1+j*(ay2-ay1)
	// i = (ay1-ax1 + j*(ay2-ay1))/(ax2-ax1)
	// 
	// bx = by
	// i*-61 + 312500000043 = j*-31 + 181818181855
	// i = (j*-31 + 181818181855 - 312500000043) / (-61)
	// 
	// ...we have two variables (i, j) and two equations, getting the formulae of i:
	// 
	// (10 + j*55)/32 = (j*-31 + 181818181855 - 312500000043) / (-61)
	// (-61 * 10) + -61*55*j = j*-31*32 + 181818181855*32 - 312500000043*32
	// j = (181818181855*32 - 312500000043*32 - (-61 * 10))/(-61*55)
	// 
	// with formulae:
	// bx1 + i*(bx2-bx1) = by1+j*(by2-by1)
	// i = (by1-bx1 + j*(by2-by1))/(bx2-bx1)
	//
	// (ay1-ax1 + j*(ay2-ay1))/(ax2-ax1) = (by1-bx1 + j*(by2-by1))/(bx2-bx1)
	// ay1-ax1 + j*(ay2-ay1) = (by1-bx1 + j*(by2-by1)) * (ax2-ax1)/(bx2-bx1)
	// ay1-ax1 + j*(ay2-ay1) = (by1-bx1)*(ax2-ax1)/(bx2-bx1) + j*(by2-by1)*(ax2-ax1)/(bx2-bx1)
	// j*(ay2-ay1) - j*(by2-by1)*(ax2-ax1)/(bx2-bx1) = (by1-bx1)*(ax2-ax1)/(bx2-bx1)-ay1+ax1 
	// j*(ay2-ay1 - (by2-by1)*(ax2-ax1)/(bx2-bx1)) = (by1-bx1)*(ax2-ax1)/(bx2-bx1) - ay1+ax1
	// j = ((by1-bx1)*(ax2-ax1)/(bx2-bx1) - ay1+ax1) / (ay2-ay1 - (by2-by1)*(ax2-ax1)/(bx2-bx1))
	// ...nice.

	// let's check.
	// {{13 49} {55 12} {10000000007975 10000000007671}}  a: 169358015095  b: 141788105668  cost: 649862150953
    // first solutions:
	// x [{5 181818181962} {60 181818181949} {115 181818181936} {170 181818181923}]
    // y [{7 833333333944} {19 833333333895} {31 833333333846} {43 833333333797}]
	// --> ax = ax1 + i*(ax2-ax1)
	// --> ax = 5 + i*55
	// i = (169358015095-5)/55 = 3079236638
	// bx = 181818181962 + i*-13 --> i = (141788105668 - 181818181962) / -13 = 3079236638 --> OK!!
	// 
	// --> ay = 7 + j*12 --> j = (169358015095-7)/12 = 14113167924
	// --> by = 833333333944 + j*-49 --> j = (141788105668-833333333944) / -49 = 14113167924 --> OK!!

    // ax1+i*(ax2-ax1) = ay1+j*(ay2-ay1)
	// 5+3079236638*55 =? 7+14113167924*12  = 169358015095 OK
	// 
	// i = (ay1-ax1 + j*(ay2-ay1))/(ax2-ax1)
	// 3079236638 ==? (7-5 + 14113167924*12)//55
	// i = (7-5 + 14113167924*12/55) = 3079236638.0 OK!


	//log.Println("x", two_x_solutions)
	//log.Println("y", two_y_solutions)
	// check for divisability
	if (bx2==bx1) {
		log.Println("not divisable")
		return 0
	}

	/*
	j := 14113167924 // ---> this is THE solution for 2nd input3.txt 
	// (ay1-ax1 + j*(ay2-ay1))/(ax2-ax1) = (by1-bx1 + j*(by2-by1))/(bx2-bx1)
	log.Println((ay1-ax1 + j*(ay2-ay1))/(ax2-ax1), "=?", (by1-bx1 + j*(by2-by1))/(bx2-bx1))
	// ay1-ax1 + j*(ay2-ay1) = (by1-bx1 + j*(by2-by1)) * (ax2-ax1)/(bx2-bx1)
	log.Println(ay1-ax1 + j*(ay2-ay1),"=?", (by1-bx1 + j*(by2-by1)) * (ax2-ax1)/(bx2-bx1))
	// ay1-ax1 + j*(ay2-ay1) = (by1-bx1)*(ax2-ax1)/(bx2-bx1) + j*(by2-by1)*(ax2-ax1)/(bx2-bx1)
	log.Println(ay1-ax1 + j*(ay2-ay1), "=?", (by1-bx1)*(ax2-ax1)/(bx2-bx1) + j*(by2-by1)*(ax2-ax1)/(bx2-bx1))
	// j*(ay2-ay1) - j*(by2-by1)*(ax2-ax1)/(bx2-bx1) = (by1-bx1)*(ax2-ax1)/(bx2-bx1) - ay1 + ax1 
	log.Println(j*(ay2-ay1) - j*((by2-by1)*(ax2-ax1)/(bx2-bx1)), "=1?", (by1-bx1)*(ax2-ax1)/(bx2-bx1)-ay1+ax1)

	log.Println("XX", j, "*", (ay2-ay1), "-", j, "*", (by2-by1)*(ax2-ax1)/(bx2-bx1), "=", j*(ay2-ay1), "-", (by2-by1)*(ax2-ax1)/(bx2-bx1)*j, "=", j*(ay2-ay1) - (by2-by1)*(ax2-ax1)/(bx2-bx1)*j)
	log.Println("XX", j, "*", (ay2-ay1 - (by2-by1)*(ax2-ax1)/(bx2-bx1)))
	log.Println("YY", (by2-by1)*(ax2-ax1)/(bx2-bx1)*j)
	log.Println("YY", j*(by2-by1)*(ax2-ax1)/(bx2-bx1))
	log.Println("YY", (by2-by1)*(ax2-ax1))
	log.Println("ZZ", (by2-by1)*(ax2-ax1), (bx2-bx1), j)
	// j*(ay2-ay1 - (by2-by1)*(ax2-ax1)/(bx2-bx1)) = (by1-bx1)*(ax2-ax1)/(bx2-bx1) - ax1+ay1
	log.Println(j*(ay2-ay1 - (by2-by1)*(ax2-ax1)/(bx2-bx1)), "=2?", (by1-bx1)*(ax2-ax1)/(bx2-bx1) - ay1+ax1)
*/

	// trying naively...
    jf := ((float64((by1-bx1)*(ax2-ax1))/float64(bx2-bx1) - float64(ay1+ax1)) / (float64(ay2-ay1) - float64((by2-by1)*(ax2-ax1))/float64(bx2-bx1)))
	ifl := ((float64(by1-bx1) + jf*float64(by2-by1))/float64(bx2-bx1))

	//log.Printf("i %.2f j %.2f", ifl, jf)
	jtmp := int(math.Round(jf))
	itmp := int(math.Round(ifl))


	// okay this is me being desperate...
	for i := itmp-111; i <= itmp+111; i++ {
		for j:= jtmp-111; j <= jtmp+111; j++ {
			if i < 0 || j < 0 {
				continue
			}
			ax := ax1 + i*(ax2-ax1)
			bx := bx1 + i*(bx2-bx1)
		
			ay := ay1+j*(ay2-ay1)
			by := by1+j*(by2-by1)
			
			//log.Println(ax , ay,  bx , by,  ax * astepx + bx * bstepx == prizex , ay * astepy + by * bstepy == prizey)
			if ax == ay && bx == by && ax * astepx + bx * bstepx == prizex && ay * astepy + by * bstepy == prizey {
				log.Println("found", ax, bx)
				return ax * COST_A + bx * COST_B
			}
		}
	} 


	//log.Println("i", i, "j", j, "ax", ax, "ay", ay, "bx", bx, "by", by)

	return 0

/*
	for count, _ := range xsolutions {
		if _, exists := ysolutions[count]; exists {
			cost := count.a*COST_A + count.b*COST_B
			if minCost == -1 || cost < minCost {
				minCost = cost
				bestA = count.a
				bestB = count.b
				//log.Println("tmp ", machine, cost, bestA, bestB)
			}
		}
	}
	log.Println(machine, " a:", bestA, " b:", bestB, " cost:", minCost)

	return minCost
	*/
}


func main() {
	for i := 0; i < 2; i++ {
		sum := 0
		if i == 0 {
			log.Println("PART 1 ------------")
		} else {
			log.Println("PART 2 ------------")

		}
		machines := readMachines("input.txt", i == 1)
		for _, machine := range machines {
			cost := 0
			if i == 0 {
				//cost = getCostForNaive(machine)
				cost = getCostForOptimized2(machine)
			} else {
				cost = getCostForOptimized2(machine)
			}
			if cost != -1 {
				sum += cost
			}
		}
		log.Println(sum)
	}
}
