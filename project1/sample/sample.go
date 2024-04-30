package sample

import (
    "fmt"
    "math"
)

// Struct definition
type Person struct {
    Name   string
    Age    int
    Height float64
}

// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Struct implementing the Shape interface
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2*r.Width + 2*r.Height
}

// Function with parameters and return type
func add(a, b int) int {
    return a + b
}

// Function with variadic parameters
func sumNumbers(numbers ...int) int {
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    return sum
}

// Function with named return value
func divide(a, b float64) (result float64) {
    if b != 0 {
        result = a / b
    }
    return
}

func main() {
    // Variable declarations
    var age int = 25
    name := "John"
    isStudent := true
    
    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("Is Student:", isStudent)
    
    // Pointer
    ptr := &age
    fmt.Println("Pointer value:", *ptr)
    
    // Struct instantiation
    person := Person{
        Name:   "Alice",
        Age:    30,
        Height: 165.5,
    }
    fmt.Println("Person:", person)
    
    // Array
    numbers := [5]int{1, 2, 3, 4, 5}
    fmt.Println("Array:", numbers)
    
    // Slice
    sliceNumbers := []int{1, 2, 3, 4, 5}
    fmt.Println("Slice:", sliceNumbers)
    
    // Appending to a slice
    sliceNumbers = append(sliceNumbers, 6, 7, 8)
    fmt.Println("Slice after appending:", sliceNumbers)
    
    // Map
    ages := map[string]int{
        "Alice": 30,
        "Bob":   35,
        "Charlie": 28,
    }
    fmt.Println("Map:", ages)
    
    // Accessing map value
    aliceAge := ages["Alice"]
    fmt.Println("Alice's age:", aliceAge)
    
    // For loop
    for i := 0; i < 5; i++ {
        fmt.Println("Loop iteration:", i)
    }
    
    // For-each loop
    for index, value := range sliceNumbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }
    
    // While loop (using for with a condition)
    count := 0
    for count < 3 {
        fmt.Println("Count:", count)
        count++
    }
    
    // Switch statement
    grade := "B"
    switch grade {
    case "A":
        fmt.Println("Excellent!")
    case "B", "C":
        fmt.Println("Well done!")
    default:
        fmt.Println("Keep working hard!")
    }
    
    // Calling a function
    result := add(5, 3)
    fmt.Println("Addition result:", result)
    
    // Calling a variadic function
    sum := sumNumbers(1, 2, 3, 4, 5)
    fmt.Println("Sum:", sum)
    
    // Calling a function with named return value
    quotient := divide(10, 3)
    fmt.Println("Division result:", quotient)
    
    // Working with interfaces
    var shape Shape = Rectangle{Width: 5, Height: 3}
    fmt.Println("Area:", shape.Area())
    fmt.Println("Perimeter:", shape.Perimeter())
    
    // Using the math package
    sqrt := math.Sqrt(16)
    fmt.Println("Square root:", sqrt)
}
