package maps

import "C"
import (
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/mbark/aoc2025/util"
)

type Map[T any] struct {
	Columns int
	Rows    int
	Cells   [][]T
}

func NewEmpty[T any](columns, rows int) Map[T] {
	cells := make([][]T, rows)
	for i := range cells {
		cells[i] = make([]T, columns)
	}

	return Map[T]{Columns: columns, Rows: rows, Cells: cells}
}

func (m Map[T]) ArraySize() int {
	return (m.Rows + 1) * (m.Columns + 1)
}

func NewIntMap(definition string) Map[int] {
	var cells [][]int

	var rows, cols int
	for y, l := range util.ReadInput(definition, "\n") {
		rows = y
		var row []int
		for x, n := range util.NumberList(l, "") {
			cols = x
			row = append(row, n)
		}

		cells = append(cells, row)
	}

	return Map[int]{Columns: cols + 1, Rows: rows + 1, Cells: cells}
}

func NewByte(definition string) Map[byte] {
	return New[byte](definition, func(x, y int, b byte) byte { return b })
}

func New[T any](definition string, fn func(x, y int, b byte) T) Map[T] {
	var cells [][]T

	var rows, cols int
	for y, l := range util.ReadInput(definition, "\n") {
		rows = y
		var row []T
		for x, n := range l {
			cols = x
			row = append(row, fn(x, y, byte(n)))
		}

		cells = append(cells, row)
	}

	return Map[T]{Columns: cols + 1, Rows: rows + 1, Cells: cells}
}

func MapFromCoordinates[T any](coordinates map[Coordinate]T) Map[T] {
	var rows, cols int
	for c := range coordinates {
		if c.Y > rows {
			rows = c.Y
		}
		if c.X > cols {
			cols = c.X
		}
	}

	rows, cols = rows+1, cols+1

	cells := make([][]T, rows)
	for i := range cells {
		cells[i] = make([]T, cols)
	}

	for c, val := range coordinates {
		cells[c.Y][c.X] = val
	}

	return Map[T]{Columns: cols, Rows: rows, Cells: cells}
}

func (m Map[T]) IterHorizontal() iter.Seq[Coordinate] {
	return func(yield func(coordinate Coordinate) bool) {
		for y := 0; y < m.Rows; y++ {
			for x := 0; x < m.Columns; x++ {
				if !yield(Coordinate{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

func (m Map[T]) IterVertical() iter.Seq[Coordinate] {
	return func(yield func(coordinate Coordinate) bool) {
		for x := 0; x < m.Columns; x++ {
			for y := 0; y < m.Rows; y++ {
				if !yield(Coordinate{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

func (m Map[T]) WithPadding(n, e, s, w int) Map[T] {
	newm := Map[T]{
		Columns: e + m.Columns + w,
		Rows:    n + m.Rows + s,
	}

	cells := make([][]T, newm.Rows)
	for i := range cells {
		cells[i] = make([]T, newm.Columns)
	}

	for _, c := range m.Coordinates() {
		cells[n+c.Y][e+c.X] = m.At(c)
	}

	newm.Cells = cells
	return newm
}

func (m Map[T]) At(c Coordinate) T {
	return m.Cells[c.Y][c.X]
}

func (m Map[T]) AtSafe(c Coordinate) T {
	var t T
	if !m.Exists(c) {
		return t
	}

	return m.Cells[c.Y][c.X]
}

// Get is just an alias for at.
func (m Map[T]) Get(c Coordinate) T {
	return m.At(c)
}

func (m Map[T]) ArrPos(c Coordinate) int {
	return c.Y*m.Rows + c.X
}

func (m Map[T]) Coordinates() []Coordinate {
	coordinates := make([]Coordinate, m.Length())
	for y, row := range m.Cells {
		for x := range row {
			coordinates[y*m.Columns+x] = Coordinate{Y: y, X: x}
		}
	}

	return coordinates
}

func (m Map[T]) CopyWith(fn func(c Coordinate, val T) T) Map[T] {
	cells := make([][]T, len(m.Cells))

	for i := range m.Cells {
		row := make([]T, len(m.Cells[i]))
		for j, val := range m.Cells[i] {
			row[j] = fn(Coordinate{X: j, Y: i}, val)
		}

		cells[i] = row
	}

	return Map[T]{Columns: m.Columns, Rows: m.Rows, Cells: cells}
}

func Merged[T any](maps [][]Map[T]) Map[T] {
	var cells [][]T
	var columns, rows int

	for _, row := range maps {
		rows += row[0].Rows
	}
	for _, col := range maps[0] {
		columns += col.Columns
	}

	// for each map in the row
	for _, mapRow := range maps {
		// for each row in the map
		for i := 0; i < mapRow[0].Rows; i++ {
			var row []T
			for _, mapCol := range mapRow {
				row = append(row, mapCol.Cells[i]...)
			}

			cells = append(cells, row)
		}
	}

	return Map[T]{Columns: columns, Rows: rows, Cells: cells}
}

func (m *Map[T]) Set(c Coordinate, val T) {
	m.Cells[c.Y][c.X] = val
}

func (m *Map[T]) Move(from, to Coordinate, d T) {
	m.Cells[to.Y][to.X], m.Cells[from.Y][from.X] = m.Cells[from.Y][from.X], d
}

func (m Map[T]) Exists(c Coordinate) bool {
	return c.X >= 0 && c.X < m.Columns &&
		c.Y >= 0 && c.Y < m.Rows
}

func (m Map[T]) WrapCoordinate(c Coordinate) Coordinate {
	return Coordinate{
		X: (c.X + m.Columns) % m.Columns,
		Y: (c.Y + m.Rows) % m.Rows,
	}
}

func (m Map[T]) filterNonExistent(coords []Coordinate) []Coordinate {
	var cs []Coordinate
	for _, c := range coords {
		if m.Exists(c) {
			cs = append(cs, c)
		}
	}

	return cs
}

func (m Map[T]) Adjacent(c Coordinate) []Coordinate {
	coordinates := make([]Coordinate, 4)
	var at int
	for _, x := range []int{-1, 1} {
		c := Coordinate{X: c.X + x, Y: c.Y}
		if m.Exists(c) {
			coordinates[at] = c
			at += 1
		}
	}
	for _, y := range []int{-1, 1} {
		c := Coordinate{X: c.X, Y: c.Y + y}
		if m.Exists(c) {
			coordinates[at] = c
			at += 1
		}
	}

	return coordinates[:at]
}

func (m Map[T]) IterAdjacent(c Coordinate) iter.Seq2[Coordinate, Direction] {
	return func(yield func(Coordinate, Direction) bool) {
		for _, x := range []int{-1, 1} {
			c := Coordinate{X: c.X + x, Y: c.Y}
			if m.Exists(c) {
				yield(c, Direction{X: x, Y: 0})
			}
		}
		for _, y := range []int{-1, 1} {
			c := Coordinate{X: c.X, Y: c.Y + y}
			if m.Exists(c) {
				yield(c, Direction{X: 0, Y: y})
			}
		}
	}
}

func (m Map[T]) Surrounding(c Coordinate) []Coordinate {
	var coordinates []Coordinate
	for _, x := range []int{-1, 0, 1} {
		for _, y := range []int{-1, 0, 1} {
			if x == 0 && y == 0 {
				continue
			}

			c := Coordinate{X: c.X + x, Y: c.Y + y}
			if m.Exists(c) {
				coordinates = append(coordinates, c)
			}
		}
	}

	return coordinates
}

func (m Map[T]) String() string {
	var sb strings.Builder
	for _, row := range m.Cells {
		for _, cell := range row {
			switch t := any(cell).(type) {
			case byte:
				sb.WriteByte(t)
			case int:
				sb.WriteString(strconv.Itoa(t))
			default:
				sb.WriteString(fmt.Sprintf("%s", cell))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Map[T]) Stringf(sprintf func(c Coordinate, val T) string) string {
	var sb strings.Builder
	for y, row := range m.Cells {
		for x, cell := range row {
			sb.WriteString(sprintf(Coordinate{X: x, Y: y}, cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Map[T]) Length() int {
	return m.Rows * m.Columns
}
