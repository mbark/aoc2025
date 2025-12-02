package day1

import (
	"fmt"

	"github.com/mbark/aoc2025/util"
)

var testInput = `
L68
L30
R48
L5
R60
L55
L1
L99
R14
L82
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var instructions []instruction
	for _, l := range util.ReadInput(input, "\n") {
		instructions = append(instructions, instruction{
			direction: l[0],
			distance:  util.ParseInt[int](l[1:]),
		})
	}

	fmt.Printf("first: %d\n", first(instructions))
	fmt.Printf("second: %d\n", second(instructions))
}

type instruction struct {
	direction byte
	distance  int
}

func first(instructions []instruction) int {
	direction := 50
	var count int
	for _, i := range instructions {
		switch i.direction {
		case 'L':
			direction -= i.distance
		case 'R':
			direction += i.distance
		}

		direction = (direction + 100) % 100
		if direction == 0 {
			count++
		}
	}
	return count
}

func second(instructions []instruction) int {
	direction := 50
	var count int
	for _, i := range instructions {
		for range i.distance {
			switch i.direction {
			case 'L':
				direction--
			case 'R':
				direction++
			}
		
			direction = (direction + 100) % 100
			if direction == 0 {
				count++
			}
		}
	}

	return count
}
