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

type Circle struct {
	X      float64
	Y      float64
	Radius float64
}

func (c Circle) IsPointInside(x float64, y float64) bool {
	return c.Radius > math.Sqrt(math.Pow(c.X-x, 2)+math.Pow(c.Y-y, 2))
}
