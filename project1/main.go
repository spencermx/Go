package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spencermx/go/project1/video"
)

func main() {
	id := uuid.New()
	fmt.Println("Generated UUID:", id)
	video.SomePackinator()
    
    circle := video.Circle{
        Radius: 19.0,
    }
     
    video.PrintArea(circle)
    fmt.Println("circle: ", circle)
}
