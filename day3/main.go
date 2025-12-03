package day3

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2025/fns"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
987654321111111
811111111111119
234234234234278
818181911112111
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var batteries []battery
	for _, l := range util.ReadInput(input, "\n") {
		joltage := fns.Map(strings.Split(l, ""), func(s string) int { return util.ParseInt[int](s) })
		batteries = append(batteries, battery{joltage})
	}

	fmt.Printf("first: %d\n", first(batteries))
	fmt.Printf("second: %d\n", second(batteries, 12))
}

type battery struct {
	joltage []int
}

func first(batteries []battery) int {
	var sum int
	for _, b := range batteries {
		var maxJ int
		var maxIdx int
		for idx := 0; idx < len(b.joltage)-1; idx++ {
			if b.joltage[idx] > maxJ {
				maxJ = b.joltage[idx]
				maxIdx = idx
			}
		}
		var maxN int
		for idx := maxIdx + 1; idx < len(b.joltage); idx++ {
			if b.joltage[idx] > maxN {
				maxN = b.joltage[idx]
			}
		}

		val := util.ParseInt[int](fmt.Sprintf("%d%d", maxJ, maxN))
		sum += val
	}

	return sum
}

func second(batteries []battery, turnOn int) int {
	var sum int
	for _, b := range batteries {
		var numbers []int
		var atIdx int
		for i := turnOn; i > 0; i-- {
			var maxJ int
			for idx := atIdx; idx < len(b.joltage)-(i-1); idx++ {
				if b.joltage[idx] > maxJ {
					maxJ = b.joltage[idx]
					atIdx = idx + 1
				}
			}

			numbers = append(numbers, maxJ)
		}

		val := util.ParseInt[int](strings.Join(fns.Map(numbers, func(n int) string { return fmt.Sprintf("%d", n) }), ""))
		sum += val
	}

	return sum
}
