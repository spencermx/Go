package common

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// File represents a file object
type BucketKey struct {
	Key string // File Key including the extension ex. "hello.png" "accounts/spencer/temp.png"
}

func (f *BucketKey) GetGuid() uuid.UUID {
	fileName := f.GetFileNameWithoutExtension()

	uuidString := fileName[:36]

	// Parse the UUID string
	parsedUUID, _ := uuid.Parse(uuidString)

	return parsedUUID
}

func (f *BucketKey) GetFileNameWithoutExtension() string {
	// Find the index of the last forward slash (/) in the Key
	slashIndex := strings.LastIndex(f.Key, "/")
	if slashIndex == -1 {
		// If no slash is found, consider the entire Key as the file name
		return removeExtension(f.Key)
	}

	// Extract the file name from the Key
	fileName := f.Key[slashIndex+1:]
	return removeExtension(fileName)
}

func (f *BucketKey) GetFileNameWithoutExtensionAndGuid() string {
	// Find the index of the last forward slash (/) in the Key
	slashIndex := strings.LastIndex(f.Key, "/")
	if slashIndex == -1 {
		// If no slash is found, consider the entire Key as the file name
        return removeExtension(f.Key[slashIndex+38:])
	}

	// Extract the file name from the Key
    fileName := f.Key[slashIndex+38:]
	return removeExtension(fileName)
}
// GetFileExtension returns the file extension of the file
func (f *BucketKey) GetFileExtension() string {
	return strings.ToLower(filepath.Ext(f.Key))
}

// IsImage checks if the file is an image based on its extension
func (f *BucketKey) IsImage() bool {
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
func (f *BucketKey) IsVideo() bool {
	imageExtensions := []string{".mp4"}
	extension := f.GetFileExtension()

	for _, imageExt := range imageExtensions {
		if extension == imageExt {
			return true
		}
	}
	return false
}

func (f *BucketKey) ToFileUrl() *FileUrl {
	fileUrl := NewFileUrlWithDefaultRootUrl(*f)

	return fileUrl
}

func removeExtension(fileName string) string {
	// Find the index of the last dot (.) in the file name
	dotIndex := strings.LastIndex(fileName, ".")
	if dotIndex == -1 {
		// If no dot is found, return the file name as is
		return fileName
	}

	// Trim the file extension from the file name
	fileNameWithoutExt := fileName[:dotIndex]
	return fileNameWithoutExt
}

// hasSuffix checks if a string ends with a specific suffix
// func hasSuffix(s, suffix string) bool {
//     return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
// }
