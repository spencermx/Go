package video_tests

import (
    "testing"
    "github.com/spencermx/go/project1/video" 
)

func TestAdd(t *testing.T) {
	result := video.Add(2, 3)
    expected := 5 
	if result != expected {
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
}
