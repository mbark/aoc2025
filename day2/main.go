package day2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbark/aoc2025/util"
)

var testInput = `
11-22,95-115,998-1012,1188511880-1188511890,222220-222224,
1698522-1698528,446443-446449,38593856-38593862,565653-565659,
824824821-824824827,2121212118-2121212124
`

func Run(input string, isTest bool) {
	if isTest {
		input = strings.ReplaceAll(testInput, "\n", "")
	}

	var ranges []rangePair
	for _, l := range util.ReadInput(input, ",") {
		ranges = append(ranges, rangePair{
			min: util.ParseInt[int](strings.Split(l, "-")[0]),
			max: util.ParseInt[int](strings.Split(l, "-")[1]),
		})
	}

	fmt.Printf("first: %d\n", first(ranges))
	fmt.Printf("second: %d\n", second(ranges))
}

func first(ranges []rangePair) int {
	var ids int
	for _, r := range ranges {
		for i := r.min; i <= r.max; i++ {
			if isRepeat(strconv.Itoa(i), 2) {
				ids += i
			}
		}
	}
	return ids
}

func second(ranges []rangePair) int {
	var ids int
	for _, r := range ranges {
		for i := r.min; i <= r.max; i++ {
			if isRepeat(strconv.Itoa(i), 0) {
				ids += i
			}
		}
	}
	return ids
}

func isRepeat(s string, maxRepeat int) bool {
	for i := 1; i <= len(s)/2; i++ {
		if len(s)%i != 0 {
			continue
		}
		if maxRepeat > 0 && len(s)/i > maxRepeat {
			continue
		}

		part := s[:i]
		isPartRepeat := true
		for idx := 0; idx < len(s)/i; idx++ {
			if s[i*idx:i*(idx+1)] != part {
				isPartRepeat = false
				break
			}
		}

		if isPartRepeat {
			return true
		}
	}

	return false
}

type rangePair struct {
	min int
	max int
}
