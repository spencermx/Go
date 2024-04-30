package video

import "github.com/giorgisio/goav" 




type VideoFile struct {
    VideoFilePath string
}

func NewVideoFile(videoFilePath string) *VideoFile {
   videoFile := VideoFile {
    VideoFilePath: videoFilePath,        
   }       
    
   return &videoFile
}


type Encoder struct {


}


func (e Encoder) EncodeVideoFileToAv1() {

}


    
