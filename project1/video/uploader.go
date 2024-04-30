package video

import "fmt"

// Define an interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Define a struct that doesn't implement the Shape interface
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return c.Radius * 2
}

func (c Circle) Perimeter() float64 {
    return c.Radius * 4
}



// Function that accepts a Shape interface parameter
func PrintArea(s Shape) {
    fmt.Println("Area:", s.Area())
}
