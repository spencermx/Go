package handlers

import (
    "path/filepath"
    "strings"
)


// File represents a file object
type File struct {
    Key string
}

// GetFileExtension returns the file extension of the file
func (f *File) GetFileExtension() string {
    return strings.ToLower(filepath.Ext(f.Key))
}

// IsImage checks if the file is an image based on its extension
func (f *File) IsImage() bool {
    imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
    extension := f.GetFileExtension()

    for _, imageExt := range imageExtensions {
        if extension == imageExt {
            return true
        }
    }
    return false
}

// IsImage checks if the file is an image based on its extension
func (f *File) IsVideo() bool {
    imageExtensions := []string{ ".mp4" }
    extension := f.GetFileExtension()

    for _, imageExt := range imageExtensions {
        if extension == imageExt {
            return true
        }
    }
    return false
}

// hasSuffix checks if a string ends with a specific suffix
// func hasSuffix(s, suffix string) bool {
//     return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
// }
