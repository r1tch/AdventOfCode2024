package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"math"
)

// Registers
var A, B, C int

func main() {
	// Read the input from input.txt
	data, err := ioutil.ReadFile("input2.txt")
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Parse the input
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	parseRegisters(lines[0], lines[1], lines[2])
	program := parseProgram(lines[4])

	// Execute the program and collect the output
	output := executeProgram(program)
	fmt.Println("Final Output:", strings.Join(output, ","))
}

// Parse the initial register values
func parseRegisters(aLine, bLine, cLine string) {
	A = parseRegisterValue(aLine)
	B = parseRegisterValue(bLine)
	C = parseRegisterValue(cLine)
}

func parseRegisterValue(line string) int {
	parts := strings.Split(line, ": ")
	value, _ := strconv.ParseInt(parts[1], 10, 0)
	return int(value)
}

// Parse the program as a slice of integers
func parseProgram(line string) []int {
	parts := strings.Split(strings.Split(line, ": ")[1], ",")
	program := make([]int, len(parts))
	for i, part := range parts {
		program[i], _ = strconv.Atoi(part)
	}
	return program
}

// Execute the program
func executeProgram(program []int) []string {
	ip := 0 // Instruction Pointer
	output := []string{}

	for ip < len(program) {
		opcode := program[ip]
		operand := program[ip+1]
		ip += 2 // Default step unless overridden
		fmt.Printf("o: %d %d A: 0%o B: 0%o C: 0%o\n", opcode, operand, A, B, C)
		

		switch opcode {
		case 0: // adv: Divide A by 2^operand (combo)
			A = A / powerOfTwo(getComboValue(operand))
		case 1: // bxl: B XOR literal operand
			B ^= operand
		case 2: // bst: B = combo operand % 8
			B = getComboValue(operand) % 8
		case 3: // jnz: Jump if A != 0
			if A != 0 {
				ip = operand
			}
		case 4: // bxc: B XOR C
			B ^= C
		case 5: // out: Output combo operand % 8
			outChr := strconv.Itoa(getComboValue(operand) % 8)
			output = append(output, outChr)
			fmt.Println("Output:", outChr)
		case 6: // bdv: Divide A by 2^operand, store in B
			B = A / powerOfTwo(getComboValue(operand))
		case 7: // cdv: Divide A by 2^operand, store in C
			C = A / powerOfTwo(getComboValue(operand))
		default:
			fmt.Println("Unknown opcode:", opcode)
			return output
		}
	}
	return output
}

// Helper: Calculate 2^value
func powerOfTwo(value int) int {
	return int(math.Pow(2, float64(value)))
}

// Helper: Get combo operand value
func getComboValue(operand int) int {
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
