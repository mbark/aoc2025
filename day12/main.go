package day12

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2025/fns"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}
	splits := util.ReadInput(input, "\n\n")
	parsePoly := func(d string) Poly {
		var poly Poly
		for i, l := range util.ReadInput(d, "\n") {
			if i == 0 {
				continue
			}

			poly = append(poly, []bool{l[0] == '#', l[1] == '#', l[2] == '#'})
		}
		return poly
	}
	polys := []Poly{
		parsePoly(splits[0]),
		parsePoly(splits[1]),
		parsePoly(splits[2]),
		parsePoly(splits[3]),
		parsePoly(splits[4]),
		parsePoly(splits[5]),
	}

	var problems []Problem
	for _, l := range util.ReadInput(splits[6], "\n") {
		s := strings.Split(l, ": ")
		ss := strings.Split(s[0], "x")
		rows := util.ParseInt[int](ss[0])
		cols := util.ParseInt[int](ss[1])
		problems = append(problems, Problem{
			Rows:  rows,
			Cols:  cols,
			Polys: util.NumberList(s[1], " "),
		})
	}

	fmt.Printf("first: %d\n", first(polys, problems))
}

func first(polys []Poly, problems []Problem) int {
	covered := map[int]int{}
	for i, p := range polys {
		for _, l := range p {
			for _, k := range l {
				if k {
					covered[i]++
				}
			}
		}
	}

	var possible int
	for _, p := range problems {
		area := p.Rows * p.Cols
		var polyArea int
		for i, count := range p.Polys {
			polyArea += covered[i] * count
		}
		if polyArea >= area {
			continue
		}
		possible++
	}

	return possible
}

type Poly [][]bool

func (p Poly) String() string {
	var ss []string
	for _, row := range p {
		ss = append(ss, strings.Join(fns.Map(row, func(b bool) string {
			if b {
				return "#"
			} else {
				return "."
			}
		}), ""))
	}
	return strings.Join(ss, "\n")
}

type Problem struct {
	Rows, Cols int
	Polys      []int
}

func (p Problem) String() string {
	return fmt.Sprintf("%dx%d: %s", p.Rows, p.Cols, strings.Join(fns.Map(p.Polys, func(i int) string { return fmt.Sprintf("%d", i) }), " "))
}
