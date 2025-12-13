package day11

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2025/util"
)

var testInput = `
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out
`

var testInput2 = `
svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`

func Run(input string, isTest bool) {
	if isTest {
		input = testInput
	}

	graph := map[string][]string{}
	for _, l := range util.ReadInput(input, "\n") {
		parts := strings.Split(l, ": ")
		graph[parts[0]] = strings.Split(parts[1], " ")
	}
	fmt.Printf("first: %d\n", first(graph))

	if isTest {
		input = testInput2
		graph = map[string][]string{}
		for _, l := range util.ReadInput(input, "\n") {
			parts := strings.Split(l, ": ")
			graph[parts[0]] = strings.Split(parts[1], " ")
		}
	}
	fmt.Printf("second: %d\n", second(graph))

}

func first(graph map[string][]string) int {
	return dfs(graph, "you", "out")
}

func second(graph map[string][]string) int {
	p1 := countPaths(graph, "svr", "fft")
	p2 := countPaths(graph, "fft", "dac")
	p3 := countPaths(graph, "dac", "out")

	return p1 * p2 * p3
}

func dfs(graph map[string][]string, start, end string) int {
	if start == end {
		return 1
	}

	var count int
	for _, n := range graph[start] {
		count += dfs(graph, n, end)
	}
	return count
}

func topologicalSort(graph map[string][]string, start string) []string {
	var order []string
	visited := make(map[string]bool)

	var dfs func(string)
	dfs = func(node string) {
		if visited[node] {
			return
		}
		visited[node] = true
		for _, neighbor := range graph[node] {
			dfs(neighbor)
		}
		order = append(order, node)
	}

	dfs(start)
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}

func countPaths(graph map[string][]string, start, end string) int {
	order := topologicalSort(graph, start)

	visited := map[string]int{start: 1}
	for _, node := range order {
		for _, neighbor := range graph[node] {
			visited[neighbor] += visited[node]
		}
	}

	return visited[end]
}
