package nodes

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Normalized() Vector {
	length := math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))

	v.X = 1 / length * v.X
	v.Y = 1 / length * v.Y

	return v
}

func (v Vector) Rotated(rad float64) Vector {
	tempX := v.X

	v.X = v.X*math.Cos(rad) - v.Y*math.Sin(rad)
	v.Y = tempX*math.Sin(rad) + v.Y*math.Cos(rad)

	return v
}
