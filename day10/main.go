package day10

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"

	"github.com/mbark/aoc2025/fns"
	"github.com/mbark/aoc2025/queue"
	"github.com/mbark/aoc2025/util"
)

var testInput = `
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	var diagrams []Diagram
	for _, l := range util.ReadInput(input, "\n") {
		var d Diagram
		lightsEnd := strings.Index(l, "]")
		for _, b := range l[1:lightsEnd] {
			if b == '#' {
				d.Lights = append(d.Lights, true)
			} else {
				d.Lights = append(d.Lights, false)
			}
		}

		joltageStart := strings.Index(l, "{")
		for _, s := range strings.Split(l[lightsEnd+2:joltageStart], " ") {
			if s == "" {
				continue
			}
			var button Button
			for _, i := range strings.Split(s[1:len(s)-1], ",") {
				button = append(button, util.Str2Int(i))
			}
			d.Buttons = append(d.Buttons, button)
		}

		for _, j := range strings.Split(l[joltageStart+1:len(l)-1], ",") {
			d.Joltage = append(d.Joltage, util.Str2Int(j))
		}
		diagrams = append(diagrams, d)
	}

	fmt.Printf("first: %d\n", first(diagrams))
	fmt.Printf("second: %d\n", second(diagrams))
}
func first(diagrams []Diagram) int {
	var sum int
	for _, d := range diagrams {
		count := bfs(d)
		sum += count
	}
	return sum
}

func second(diagrams []Diagram) int {
	var sum int
	for i, d := range diagrams {
		fmt.Println("diagram", i, "/", len(diagrams), ":", d)
		count := bfs2(d)
		sum += count
	}
	return sum
}

func bfs(d Diagram) int {
	initial := make(Lights, len(d.Lights))
	queue := []clicks{{lights: initial, count: 0}}
	visited := map[string]bool{initial.String(): true}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		for _, b := range d.Buttons {
			pressed := b.Press(next.lights)
			if visited[pressed.String()] {
				continue
			}
			if d.IsDone(pressed) {
				return next.count + 1
			}

			visited[pressed.String()] = true
			queue = append(queue, clicks{lights: pressed, count: next.count + 1})
		}
	}

	return 0
}

func bfs2(d Diagram) int {
	initial := make(Joltage, len(d.Joltage))
	pq := queue.PriorityQueue[clicksJ]{}
	heap.Init(&pq)
	heap.Push(&pq, &queue.Item[clicksJ]{Value: clicksJ{joltage: initial, count: 0}, Priority: 0})
	visited := map[string]bool{initial.String(): true}
	var best int

	for len(pq) > 0 {
		q := heap.Pop(&pq).(*queue.Item[clicksJ])
		next := q.Value
		fmt.Println("states", len(pq), "visits", len(visited), "next", next.count, next.joltage, d.Diff(next.joltage))

		for _, b := range d.Buttons {
			pressed := b.PressJ(next.joltage)
			if best > 0 && d.Diff(pressed) >= best {
				continue
			}
			if visited[pressed.String()] {
				continue
			}
			if d.IsDoneJ(pressed) {
				best = next.count + 1
			}
			if d.IsImpossible(pressed) {
				continue
			}

			visited[pressed.String()] = true
			heap.Push(&pq, &queue.Item[clicksJ]{Value: clicksJ{joltage: pressed, count: next.count + 1}, Priority: d.Diff(pressed)})
		}
	}

	return best
}

func bfs2(d Diagram) int {
	initial := make(Joltage, len(d.Joltage))
	pq := queue.PriorityQueue[clicksJ]{}
	heap.Init(&pq)
	heap.Push(&pq, &queue.Item[clicksJ]{Value: clicksJ{joltage: initial, count: 0}, Priority: 0})
	visited := map[string]bool{initial.String(): true}
	var best int

	for len(pq) > 0 {
		q := heap.Pop(&pq).(*queue.Item[clicksJ])
		next := q.Value
		fmt.Println("states", len(pq), "visits", len(visited), "next", next.count, next.joltage, d.Diff(next.joltage))

		for _, b := range d.Buttons {
			pressed := b.PressJ(next.joltage)
			if best > 0 && d.Diff(pressed) >= best {
				continue
			}
			if visited[pressed.String()] {
				continue
			}
			if d.IsDoneJ(pressed) {
				best = next.count + 1
			}
			if d.IsImpossible(pressed) {
				continue
			}

			visited[pressed.String()] = true
			heap.Push(&pq, &queue.Item[clicksJ]{Value: clicksJ{joltage: pressed, count: next.count + 1}, Priority: d.Diff(pressed)})
		}
	}

	return best
}

type clicks struct {
	lights Lights
	count  int
}

type clicksJ struct {
	joltage Joltage
	count   int
}

type Diagram struct {
	Lights  Lights
	Buttons []Button
	Joltage Joltage
}

func (d Diagram) IsDone(l Lights) bool {
	for i := range l {
		if l[i] != d.Lights[i] {
			return false
		}
	}

	return true
}

func (d Diagram) IsDoneJ(j Joltage) bool {
	for i := range j {
		if j[i] != d.Joltage[i] {
			return false
		}
	}

	return true
}

func (d Diagram) Diff(j Joltage) int {
	var diff int
	for i := range j {
		diff += d.Joltage[i] - j[i]
	}
	return diff
}

func (d Diagram) IsImpossible(j Joltage) bool {
	for i := range j {
		if j[i] > d.Joltage[i] {
			return true
		}
	}

	return false
}

func (d Diagram) String() string {
	return fmt.Sprintf("[%s] %s {%s}", d.Lights.String(),
		strings.Join(fns.Map(d.Buttons, func(b Button) string {
			return "(" + strings.Join(fns.Map(b, func(i int) string { return fmt.Sprintf("%d", i) }), ",") + ")"
		}), " "),
		strings.Join(fns.Map(d.Joltage, func(j int) string { return fmt.Sprintf("%d", j) }), ","),
	)
}

type Lights []bool

func (l Lights) String() string {
	var sb strings.Builder
	for _, b := range l {
		if b {
			sb.WriteString("#")
		} else {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

type Button []int

func (b Button) String() string {
	var sb []string
	for _, i := range b {
		sb = append(sb, strconv.Itoa(i))
	}
	return fmt.Sprintf("(%s)", strings.Join(sb, ","))
}

func (b Button) Press(original Lights) Lights {
	copied := make(Lights, len(original))
	copy(copied, original)

	for _, i := range b {
		copied[i] = !original[i]
	}
	return copied
}

func (b Button) PressJ(original Joltage) Joltage {
	copied := make(Joltage, len(original))
	copy(copied, original)

	for _, i := range b {
		copied[i] += 1
	}
	return copied
}

type Joltage []int

func (j Joltage) String() string {
	var sb []string
	for _, i := range j {
		sb = append(sb, strconv.Itoa(i))
	}
	return fmt.Sprintf("%s", strings.Join(sb, ","))
}
