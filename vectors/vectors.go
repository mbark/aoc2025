package vectors

import "fmt"

type Vector struct {
	X, Y, Z float64
}

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vector) ScalarMultiply(s float64) Vector {
	return Vector{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

func (a Vector) Divide(s float64) Vector {
	return Vector{
		X: a.X / s,
		Y: a.Y / s,
		Z: a.Z / s,
	}
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) String() string {
	return fmt.Sprintf("(%.2f,%.2f,%.2f)", a.X, a.Y, a.Z)
}
