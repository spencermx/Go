package video_tests

import (
    "fmt"
	"testing"
    "github.com/spencermx/go/project1/video" 
)

func TestEncodeVideoFile(t *testing.T){
    videoFilePath := "videoFilePath" // whats the type?

    var videoFile *video.VideoFile = video.NewVideoFile(videoFilePath)

     
    

    fmt.Println("", videoFile) 


//     video.SomePackinator()

//      
    
//    result := video.Add(2, 3)
//    expected := 6 
//    if result != expected {
//        t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
//    }
}



