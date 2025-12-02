package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mbark/aoc2025/day1"
	"github.com/mbark/aoc2025/day10"
	"github.com/mbark/aoc2025/day11"
	"github.com/mbark/aoc2025/day12"
	"github.com/mbark/aoc2025/day2"
	"github.com/mbark/aoc2025/day3"
	"github.com/mbark/aoc2025/day4"
	"github.com/mbark/aoc2025/day5"
	"github.com/mbark/aoc2025/day6"
	"github.com/mbark/aoc2025/day7"
	"github.com/mbark/aoc2025/day8"
	"github.com/mbark/aoc2025/day9"
	"github.com/mbark/aoc2025/util"
)

func main() {
	var (
		flagDay    = flag.Int("day", 0, "use test input")
		flagTest   = flag.Bool("test", false, "use test input")
		cpuprofile = flag.Bool("profile", false, "write cpu profile to file")
	)
	flag.Parse()

	if *cpuprofile {
		fmt.Println("using cpu profile")
		fn := util.WithProfiling()
		defer fn()
	}

	var input string
	if !*flagTest {
		input = util.GetInput(*flagDay)
	}

	switch *flagDay {
	case 1:
		day1.Run(input, *flagTest)
	case 2:
		day2.Run(input, *flagTest)
	case 3:
		day3.Run(input, *flagTest)
	case 4:
		day4.Run(input, *flagTest)
	case 5:
		day5.Run(input, *flagTest)
	case 6:
		day6.Run(input, *flagTest)
	case 7:
		day7.Run(input, *flagTest)
	case 8:
		day8.Run(input, *flagTest)
	case 9:
		day9.Run(input, *flagTest)
	case 10:
		day10.Run(input, *flagTest)
	case 11:
		day11.Run(input, *flagTest)
	case 12:
		day12.Run(input, *flagTest)
	default:
		fmt.Println("not implemented")
		os.Exit(1)
	}
}
