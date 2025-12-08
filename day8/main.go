package day8

import (
	"fmt"
	"math"
	"slices"

	"github.com/mbark/aoc2025/maps"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689
`

func Run(input string, isTest bool) {
	connections := 1000
	if isTest {
		input = testInput
		connections = 10
	}

	var boxes []maps.Coordinate3D
	for _, l := range util.ReadInput(input, "\n") {
		boxes = append(boxes, maps.NewCoordinate3D(l))
	}

	fmt.Printf("first: %d\n", first(boxes, connections))
	fmt.Printf("second: %d\n", second(boxes))
}

type distance struct {
	b1, b2 maps.Coordinate3D
	dist   float64
}

func first(boxes []maps.Coordinate3D, connections int) int {
	connected := make(map[maps.Coordinate3D]map[maps.Coordinate3D]bool)
	for _, b := range boxes {
		connected[b] = make(map[maps.Coordinate3D]bool)
	}
	var distances []distance
	for i, b1 := range boxes {
		for j, b2 := range boxes {
			if i <= j {
				continue
			}
			distances = append(distances, distance{b1, b2, b1.EuclideanDistance(b2)})
		}
	}
	slices.SortFunc(distances, func(a, b distance) int { return int(math.Round(a.dist - b.dist)) })

	for _, d := range distances[:connections] {
		connected[d.b1][d.b2] = true
		connected[d.b2][d.b1] = true
	}
	visited := make(map[maps.Coordinate3D]bool)
	var sizes []int
	for _, b := range boxes {
		if visited[b] {
			continue
		}

		sizes = append(sizes, bfs(b, connected, visited))
	}
	slices.Sort(sizes)
	slices.Reverse(sizes)
	return sizes[0] * sizes[1] * sizes[2]
}

func second(boxes []maps.Coordinate3D) int {
	connected := make(map[maps.Coordinate3D]map[maps.Coordinate3D]bool)
	for _, b := range boxes {
		connected[b] = make(map[maps.Coordinate3D]bool)
	}
	var distances []distance
	for i, b1 := range boxes {
		for j, b2 := range boxes {
			if i <= j {
				continue
			}
			distances = append(distances, distance{b1, b2, b1.EuclideanDistance(b2)})
		}
	}
	slices.SortFunc(distances, func(a, b distance) int { return int(math.Round(a.dist - b.dist)) })

	for _, d := range distances {
		connected[d.b1][d.b2] = true
		connected[d.b2][d.b1] = true
		if bfs(d.b1, connected, make(map[maps.Coordinate3D]bool)) == len(boxes) {
			return d.b1.X * d.b2.X
		}
	}
	return 0
}

func bfs(start maps.Coordinate3D, connected map[maps.Coordinate3D]map[maps.Coordinate3D]bool, visited map[maps.Coordinate3D]bool) int {
	size := 1
	nexts := []maps.Coordinate3D{start}
	visited[start] = true
	for len(nexts) > 0 {
		next := nexts[0]
		nexts = nexts[1:]
		for n := range connected[next] {
			if visited[n] {
				continue
			}
			visited[n] = true
			size++
			nexts = append(nexts, n)
		}
	}

	return size
}

func findShortest(boxes []maps.Coordinate3D, connected map[maps.Coordinate3D]map[maps.Coordinate3D]bool) (maps.Coordinate3D, maps.Coordinate3D) {
	var minDist float64
	var s1, s2 maps.Coordinate3D
	for i, b1 := range boxes {
		for j, b2 := range boxes {
			if i <= j {
				continue
			}
			if connected[b1][b2] {
				continue
			}

			dist := b1.EuclideanDistance(b2)
			if dist < minDist || minDist == 0 {
				minDist = dist
				s1, s2 = b1, b2
			}
		}
	}

	return s1, s2
}
