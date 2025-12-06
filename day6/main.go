package day6

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2025/util"
)

var testInput = `
123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var grid [][]int
	var operands []rune

	readInput := util.ReadInput(input, "\n")
	for i, line := range readInput {
		if len(readInput)-1 == i {
			for _, s := range strings.Split(line, " ") {
				if len(s) == 0 {
					continue
				}
				if len(s) > 1 {
					panic(fmt.Sprintf("invalid input: %s", line))
				}
				operands = append(operands, rune(s[0]))
			}
			continue
		}

		var gridLine []int
		s := strings.Split(line, " ")
		for _, v := range s {
			if len(v) == 0 {
				continue
			}
			gridLine = append(gridLine, util.Str2Int(v))
		}
		grid = append(grid, gridLine)
	}

	fmt.Printf("first: %d\n", first(grid, operands))
	fmt.Printf("second: %d\n", second(readInput[:len(readInput)-1], operands))
}

func first(grid [][]int, operands []rune) int {
	var totalSum int
	size := len(grid[0])
	for j := 0; j < size; j++ {
		sum := grid[0][j]
		for i := 1; i < len(grid); i++ {
			switch operands[j] {
			case '+':
				sum += grid[i][j]
			case '*':
				sum *= grid[i][j]
			}
		}
		totalSum += sum
	}
	return totalSum
}

func second(in []string, operands []rune) int {
	var problems [][]int
	var problem []int

	size := len(in[0])
	for j := size - 1; j >= 0; j-- {
		var numString []byte
		for i := 0; i < len(in); i++ {
			if in[i][j] != ' ' {
				numString = append(numString, in[i][j])
			}
		}

		if len(numString) == 0 {
			problems = append(problems, problem)
			problem = nil
			continue
		}
		problem = append(problem, util.Str2Int(string(numString)))
	}
	if len(problem) > 0 {
		problems = append(problems, problem)
	}

	var total int
	for idx, p := range problems {
		sum := p[0]
		operand := operands[len(operands)-1-idx]
		for i := 1; i < len(p); i++ {
			switch operand {
			case '+':
				sum += p[i]
			case '*':
				sum *= p[i]
			}
		}
		total += sum
	}
	return total
}
