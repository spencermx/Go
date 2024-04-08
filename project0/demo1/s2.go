package main

import (
	"fmt"
	"math"
)

// Shape is an interface for geometric shapes
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle represents a rectangle shape
type Rectangle struct {
	Width float64

	Height float64
}

// Area calculates the area of a rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of a rectangle
func (r Rectangle) Perimeter() float64 {
	return 2*r.Width + 2*r.Height
}

// Circle represents a circle shape
type Circle struct {
	Radius float64
}

// Area calculates the area of a circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of a circle
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	shapes := []Shape{
		Rectangle{Width: 5, Height: 3},
		Circle{Radius: 4},
	}

	for _, shape := range shapes {
		fmt.Printf("Shape: %T\n", shape)
		fmt.Printf("Area: %.2f\n", shape.Area())
		fmt.Printf("Perimeter: %.2f\n", shape.Perimeter())
		fmt.Println()
	}
}
