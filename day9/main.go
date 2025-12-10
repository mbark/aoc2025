package day9

import (
	"fmt"
	"slices"

	"github.com/mbark/aoc2025/maps"
	"github.com/mbark/aoc2025/maths"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var coords []maps.Coordinate
	for _, l := range util.ReadInput(input, "\n") {
		coords = append(coords, maps.CoordinateFromString(l))
	}

	fmt.Printf("first: %d\n", first(coords))
	fmt.Printf("second: %d\n", second(coords))
}

func first(coords []maps.Coordinate) int {
	var largest int
	var l1, l2 maps.Coordinate
	for i, c1 := range coords {
		for j, c2 := range coords {
			if i <= j {
				continue
			}

			dist := c1.ManhattanDistance(c2)
			if dist > largest {
				largest = dist
				l1, l2 = c1, c2
			}
		}
	}

	return (1 + maths.AbsInt(l1.X-l2.X)) * (1 + maths.AbsInt(l1.Y-l2.Y))
}

type Polygon []maps.Coordinate

type size struct {
	c1, c2 maps.Coordinate
	size   int
}

func second(coords []maps.Coordinate) int {
	polygon := getPolygon(coords)

	var sizes []size
	for i, c1 := range coords {
		for j, c2 := range coords {
			if i <= j {
				continue
			}

			sizes = append(sizes, size{c1: c1, c2: c2, size: (1 + maths.AbsInt(c1.X-c2.X)) * (1 + maths.AbsInt(c1.Y-c2.Y))})
		}
	}
	slices.SortFunc(sizes, func(a, b size) int { return b.size - a.size })

	for _, s := range sizes {
		if !checkBox(polygon, s.c1, s.c2) {
			continue
		}

		return s.size
	}

	return 0
}

type Box struct {
	xMin, yMin, xMax, yMax int
}

func checkBox(polygon Polygon, c1, c2 maps.Coordinate) bool {
	box := Box{
		xMin: maths.MinInt(c1.X, c2.X),
		yMin: maths.MinInt(c1.Y, c2.Y),
		xMax: maths.MaxInt(c1.X, c2.X),
		yMax: maths.MaxInt(c1.Y, c2.Y),
	}
	nw := maps.Coordinate{X: box.xMin, Y: box.yMin}
	ne := maps.Coordinate{X: box.xMax, Y: box.yMin}
	se := maps.Coordinate{X: box.xMax, Y: box.yMax}
	sw := maps.Coordinate{X: box.xMin, Y: box.yMax}

	for _, c := range []maps.Coordinate{nw, ne, se, sw} {
		if !inPolygon(c, polygon) {
			return false
		}
	}

	for _, v := range polygon {
		if v.X > box.xMin && v.X < box.xMax && v.Y > box.yMin && v.Y < box.yMax {
			return false
		}
	}

	return true
}

func inPolygon(c maps.Coordinate, polygon Polygon) bool {
	for _, p := range polygon {
		if p == c {
			return true
		}
	}

	inside := false
	for i := 0; i < len(polygon); i++ {
		j := (i + 1) % len(polygon)
		xi, yi := polygon[i].X, polygon[i].Y
		xj, yj := polygon[j].X, polygon[j].Y

		// Check if ray crosses edge
		if ((yi > c.Y) != (yj > c.Y)) &&
			(c.X < (xj-xi)*(c.Y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
	}
	return inside
}

func getPolygon(coords []maps.Coordinate) Polygon {
	var border []maps.Coordinate
	for i := 0; i < len(coords); i++ {
		next := i + 1
		if next == len(coords) {
			next = 0
		}
		border = append(border, getBorder(coords[i], coords[next])...)
	}
	return border
}

func getBorder(c1, c2 maps.Coordinate) []maps.Coordinate {
	var border []maps.Coordinate
	diff := c2.Sub(c1)
	var add maps.Direction
	switch {
	case diff.X > 0:
		add = maps.Right
	case diff.X < 0:
		add = maps.Left
	case diff.Y > 0:
		add = maps.Down
	case diff.Y < 0:
		add = maps.Up
	default:
		return nil
	}

	for at := c1; at != c2; at = add.Apply(at) {
		border = append(border, at)
	}

	return border
}
