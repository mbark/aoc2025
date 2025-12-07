package day7

import (
	"fmt"

	"github.com/mbark/aoc2025/fns"
	"github.com/mbark/aoc2025/maps"
)

var testInput = `
.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var start maps.Coordinate
	m := maps.New(input, func(x, y int, b byte) byte {
		if b == 'S' {
			start = maps.C(x, y)
		}
		return b
	})

	fmt.Printf("first %d\n", first(m, start))
	fmt.Printf("second %d\n", second(m, start))
}

func first(m maps.Map[byte], start maps.Coordinate) int {
	var splits int
	beams := []maps.Coordinate{start}
	beamAt := map[maps.Coordinate]bool{start: true}
	for i := 0; len(beams) > 0; i++ {
		beamsNext := map[maps.Coordinate]bool{}
		for _, b := range beams {
			next := b.Down()
			if !m.Exists(next) {
				continue
			}

			switch m.At(next) {
			case '.':
				beamsNext[next] = true
			case '^':
				beamsNext[next.Left()] = true
				beamsNext[next.Right()] = true
				splits++
			}
		}

		beams = fns.Keys(beamsNext)
		for _, b := range beams {
			beamAt[b] = true
		}
	}

	return splits
}

func second(m maps.Map[byte], start maps.Coordinate) int {
	beams := map[maps.Coordinate]int{start: 1}
	ends := map[maps.Coordinate]int{}
	for i := 0; len(beams) > 0; i++ {
		beamsNext := map[maps.Coordinate]int{}
		for b, v := range beams {
			next := b.Down()
			if !m.Exists(next) {
				ends[b] = v
				continue
			}

			switch m.At(next) {
			case '.':
				beamsNext[next] += v
			case '^':
				beamsNext[next.Left()] += v
				beamsNext[next.Right()] += v
			}
		}

		beams = beamsNext
	}

	var sum int
	for _, v := range ends {
		sum += v
	}
	return sum
}
