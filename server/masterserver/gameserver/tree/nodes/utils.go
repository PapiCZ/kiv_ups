package nodes

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v Vector) Normalized() Vector {
	length := v.Length()

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

func (v Vector) Multiplied(multiplier float64) Vector {
	v.X *= multiplier
	v.Y *= multiplier

	return v
}

// Angle returns angle between 2 vectors
func (v Vector) Angle(v2 Vector) float64 {
	product := v.X*v2.X + v.Y*v2.Y
	vLength := v.Length()
	v2Length := v2.Length()

	return math.Acos(product / vLength * v2Length)
}

func DegToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

type Circle struct {
	X      float64
	Y      float64
	Radius float64
}

func (c Circle) IsPointInside(x float64, y float64) bool {
	return c.Radius > math.Sqrt(math.Pow(c.X-x, 2)+math.Pow(c.Y-y, 2))
}
