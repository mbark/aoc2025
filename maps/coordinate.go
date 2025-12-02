package maps

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2025/maths"
	"github.com/mbark/aoc2025/util"
)

var CZero = Coordinate{}

type Coordinate struct {
	X int
	Y int
}

func C(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

func CoordinateFromString(s string) Coordinate {
	split := strings.Split(s, ",")
	return Coordinate{X: util.ParseInt[int](split[0]), Y: util.ParseInt[int](split[1])}
}

func (c Coordinate) Up() Coordinate {
	return Coordinate{X: c.X, Y: c.Y - 1}
}

func (c Coordinate) Right() Coordinate {
	return Coordinate{X: c.X + 1, Y: c.Y}
}

func (c Coordinate) Down() Coordinate {
	return Coordinate{X: c.X, Y: c.Y + 1}
}

func (c Coordinate) Left() Coordinate {
	return Coordinate{X: c.X - 1, Y: c.Y}
}

func (c Coordinate) Add(co Coordinate) Coordinate {
	return Coordinate{X: c.X + co.X, Y: c.Y + co.Y}
}

func (c Coordinate) Sub(co Coordinate) Coordinate {
	return Coordinate{X: c.X - co.X, Y: c.Y - co.Y}
}

func (c Coordinate) Neg() Coordinate {
	return Coordinate{X: -c.X, Y: -c.Y}
}

func (c Coordinate) Adjacent() []Coordinate {
	return []Coordinate{
		{X: c.X, Y: c.Y + 1}, // up
		{X: c.X + 1, Y: c.Y}, // right
		{X: c.X, Y: c.Y - 1}, // down
		{X: c.X - 1, Y: c.Y}, // left
	}
}

func (c Coordinate) Apply(dirs ...Direction) Coordinate {
	for _, d := range dirs {
		c = d.Apply(c)
	}

	return c
}

func (c Coordinate) Surrounding() []Coordinate {
	return []Coordinate{
		{X: c.X, Y: c.Y - 1},     // N
		{X: c.X, Y: c.Y + 1},     // S
		{X: c.X + 1, Y: c.Y},     // W
		{X: c.X - 1, Y: c.Y},     // E
		{X: c.X + 1, Y: c.Y - 1}, // NE
		{X: c.X + 1, Y: c.Y + 1}, // SE
		{X: c.X - 1, Y: c.Y + 1}, // SW
		{X: c.X - 1, Y: c.Y - 1}, // NW
	}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(x=%d,y=%d)", c.X, c.Y)
}

func (c Coordinate) IsZero() bool {
	return c.X == 0 && c.Y == 00
}

func (c Coordinate) ManhattanDistance(co Coordinate) int {
	return maths.AbsInt(c.X-co.X) + maths.AbsInt(c.Y-co.Y)
}

type Direction struct{ X, Y int }

var NoDirection = Direction{}

var (
	Up    = Direction{Y: -1}
	Right = Direction{X: 1}
	Down  = Direction{Y: 1}
	Left  = Direction{X: -1}

	North     = Up
	East      = Right
	South     = Down
	West      = Left
	NorthEast = Direction{Y: -1, X: 1}
	NorthWest = Direction{Y: -1, X: -1}
	SouthEast = Direction{Y: 1, X: 1}
	SouthWest = Direction{Y: 1, X: -1}
)

func (d Direction) Rotate(direction Direction) Direction {
	order := []Direction{Up, Right, Down, Left}
	index := map[Direction]int{Up: 0, Right: 1, Down: 2, Left: 3}

	switch direction {
	case Right:
		return order[(index[d]+1)%len(index)]
	case Left:
		return order[(len(index)+index[d]-1)%len(index)]
	default:
		return d
	}
}

func (d Direction) Opposite() Direction {
	switch d {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	case NorthEast:
		return SouthWest
	case NorthWest:
		return NorthEast
	case SouthEast:
		return NorthWest
	case SouthWest:
		return NorthEast
	default:
		panic(fmt.Sprintf("unknown direction: %s", d))
	}
}

func (d Direction) Apply(c Coordinate) Coordinate {
	return Coordinate{X: c.X + d.X, Y: c.Y + d.Y}
}

func (d Direction) ApplyN(c Coordinate, n int) Coordinate {
	return Coordinate{X: c.X + n*d.X, Y: c.Y + n*d.Y}
}

func (d Direction) String() string {
	switch d {
	case Left:
		return "<"
	case Right:
		return ">"
	case Up:
		return "^"
	case Down:
		return "v"
	case North:
		return "N"
	case East:
		return "E"
	case West:
		return "W"
	case South:
		return "S"
	case NorthEast:
		return "NE"
	case NorthWest:
		return "NW"
	case SouthEast:
		return "SE"
	case SouthWest:
		return "SW"
	}

	return ""
}

func DirectionFromString(s string) Direction {
	switch s {
	case "<":
		return Left
	case ">":
		return Right
	case "^":
		return Up
	case "v":
		return Down
	case "N":
		return North
	case "E":
		return East
	case "W":
		return West
	case "S":
		return South
	case "NE":
		return NorthEast
	case "NW":
		return NorthWest
	case "SE":
		return SouthEast
	case "SW":
		return SouthWest
	}

	panic(fmt.Sprintf("unknown direction: '%s'", s))
}

type CoordinateArray struct {
	Coordinates []Coordinate

	Y int
	X int
}

func NewCoordinateArray(coords []Coordinate) CoordinateArray {
	var x, y int
	for _, c := range coords {
		if c.X > x {
			x = c.X
		}
		if c.Y > y {
			y = c.Y
		}
	}

	return CoordinateArray{Coordinates: coords, X: x, Y: y}
}

func (arr CoordinateArray) Size() int {
	return arr.Y * arr.X
}

func (arr CoordinateArray) Index(c Coordinate) int {
	return c.Y*arr.Y + c.X
}

func (arr CoordinateArray) Coordinate(i int) Coordinate {
	y := i / 13
	x := i % 13

	return Coordinate{Y: y, X: x}
}
