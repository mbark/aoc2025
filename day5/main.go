package day5

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mbark/aoc2025/fns"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
3-5
10-14
16-20
12-18

1
5
8
11
17
32
`

var testInput2 = `
1-7
4-4
6-8

1
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput2
	}

	var ranges []Range
	var ingredients []int

	splits := strings.Split(input, "\n\n")
	for _, line := range util.ReadInput(splits[0], "\n") {
		split := strings.Split(line, "-")
		ranges = append(ranges, Range{util.ParseInt[int](split[0]), util.ParseInt[int](split[1])})
	}
	for _, line := range util.ReadInput(splits[1], "\n") {
		ingredients = append(ingredients, util.ParseInt[int](line))
	}

	fmt.Printf("first: %d\n", first(ranges, ingredients))
	fmt.Printf("second: %d\n", second(ranges))
}

func first(ranges []Range, ingredients []int) int {
	var fresh int
	for _, i := range ingredients {
		if fns.Some(ranges, func(ir Range) bool { return ir.Contains(i) }) {
			fresh++
		}
	}
	return fresh
}

func second(ranges []Range) int {
	var sum int
	for _, r := range AggregateRanges(ranges) {
		fmt.Printf("range: %s\n", r)
		sum += r.End - r.Start + 1
	}
	return sum
}

func (r Range) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

func (r Range) Contains(i int) bool {
	return i >= r.Start && i <= r.End
}

type Range struct {
	Start int
	End   int
}
type Multirange []Range

func AggregateRanges(ranges Multirange) (aggregated Multirange) {
	if len(ranges) == 0 {
		return ranges
	}

	SortRanges(ranges)
	active := ranges[0]
	for _, r := range ranges[1:] {
		switch {
		case active.End < r.Start-1: // gap between ranges
			aggregated = append(aggregated, active)
			active = r
		case active.End >= r.End: // active is a superset of r, ignore
			continue
		default: // overlap between the ranges
			active.End = r.End
		}
	}

	if len(aggregated) == 0 || active.Start > aggregated[len(aggregated)-1].End {
		aggregated = append(aggregated, active)
	}
	return aggregated
}

func SortRanges(ranges []Range) {
	sort.SliceStable(ranges, func(i, j int) bool { return ranges[i].End < ranges[j].End })
	sort.SliceStable(ranges, func(i, j int) bool { return ranges[i].Start < ranges[j].Start })
}
