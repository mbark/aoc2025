package day4

import (
	"fmt"

	"github.com/mbark/aoc2025/maps"
)

var testInput = `
..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	m := maps.New(input, func(x, y int, b byte) bool { return b == '@' })

	fmt.Printf("first %d\n", first(m))
	fmt.Printf("second %d\n", second(m))
}

func first(m maps.Map[bool]) int {
	var accessible int
	for _, c := range m.Coordinates() {
		if !m.At(c) {
			continue
		}

		var tot int
		for _, a := range m.Surrounding(c) {
			if m.At(a) {
				tot++
			}
		}
		if tot < 4 {
			accessible++
		}
	}

	return accessible
}

func second(m maps.Map[bool]) int {
	var removed int
	for {
		var toRemove []maps.Coordinate
		for _, c := range m.Coordinates() {
			if !m.At(c) {
				continue
			}

			var tot int
			for _, a := range m.Surrounding(c) {
				if m.At(a) {
					tot++
				}
			}
			if tot < 4 {
				toRemove = append(toRemove, c)
			}
		}

		if len(toRemove) == 0 {
			break
		}
		removed += len(toRemove)
		for _, c := range toRemove {
			m.Set(c, false)
		}
	}
	return removed
}
