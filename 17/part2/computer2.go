package main

import (
	"fmt"
	"math"
)

func main() {
	// Program instructions as a sequence of integers
	program := []int{2, 4, 1, 1, 7, 5, 1, 5, 0, 3, 4, 4, 5, 5, 3, 0}

	// Find the lowest A that makes the program output itself
	lowestA := findLowestA(program)
	fmt.Println("Lowest positive A:", lowestA)
}

func findLowestA(program []int) int {
	possibleA := map[int]bool{}

	// Start backward and populate possible A values
	for ip := len(program) - 2; ip >= 0; ip -= 2 {
		opcode := program[ip]
		operand := program[ip+1]

		newCandidates := map[int]bool{}

		for candidate := range possibleA {
			switch opcode {
			case 0: // adv: Reverse A division
				for i := 0; i < 8; i++ { // Allow multiple reverse powers of 2
					newCandidates[candidate*int(math.Pow(2, float64(i)))] = true
				}
			case 5: // out: Reverse output operation
				if candidate%8 == operand {
					newCandidates[candidate] = true
				}
			}
		}

		if len(possibleA) == 0 { // First out instruction
			for i := 0; i < 8; i++ { // Seed A values
				newCandidates[operand] = true
			}
		}

		possibleA = newCandidates
	}

	// Test candidates forward
	for candidate := range possibleA {
		if verifyA(candidate, program) {
			return candidate
		}
	}

	return -1
}

func verifyA(A int, program []int) bool {
	B, C := 0, 0
	ip := 0
	output := []int{}

	for ip < len(program) {
		opcode := program[ip]
		operand := program[ip+1]
		ip += 2

		switch opcode {
		case 0:
			A = A / powerOfTwo(getComboValue(operand, A, B, C))
		case 1:
			B ^= operand
		case 2:
			B = getComboValue(operand, A, B, C) % 8
		case 3:
			if A != 0 {
				ip = operand
			}
		case 4:
			B ^= C
		case 5:
			output = append(output, getComboValue(operand, A, B, C)%8)
		case 6:
			B = A / powerOfTwo(getComboValue(operand, A, B, C))
		case 7:
			C = A / powerOfTwo(getComboValue(operand, A, B, C))
		}
	}

	for i, v := range output {
		if v != program[i] {
			return false
		}
	}

	return true
}

func getComboValue(operand, A, B, C int) int {
	switch operand {
	case 4:
		return A
	case 5:
		return B
	case 6:
		return C
	default:
		return operand
	}
}

func powerOfTwo(value int) int {
	return 1 << value
}
