package handlers

import (
    "bytes"
    "io"
    "net/http"
    "strings"
    "mime/multipart"
)

type MultipartFile struct {
    File multipart.File
}

func (m *MultipartFile) IsImage() bool {
    contentType, err := m.detectContentType()
    if err != nil {
        return false
    }
    return strings.HasPrefix(contentType, "image/")
}

func (m *MultipartFile) IsVideo() bool {
    contentType, err := m.detectContentType()
    if err != nil {
        return false
    }
    return strings.HasPrefix(contentType, "video/")
}

func (m *MultipartFile) detectContentType() (string, error) {
    var buf bytes.Buffer
    
    _, err := io.CopyN(&buf, m.File, 512)
    if err != nil {
        return "", err
    }

    _, err = m.File.Seek(0, 0)
    if err != nil {
        return "", err
    }

    contentType := http.DetectContentType(buf.Bytes())
    return contentType, nil
}

//    // Check the file content type
//    var fileHeader []byte = make([]byte, 512)
//
//    if _, err := file.Read(fileHeader); err != nil {
//        http.Error(w, err.Error(), http.StatusInternalServerError)
//        return
//    }
//
//    // Reset the file reader to the beginning
//    if _, err := file.Seek(0, 0); err != nil {
//        http.Error(w, err.Error(), http.StatusInternalServerError)
//        return
//    }
//
//    // Detect the file content type
//    contentType := http.DetectContentType(fileHeader)
//
//    // Check if the content type is an image
//    if !strings.HasPrefix(contentType, "image/") {
//        http.Error(w, "Uploaded file is not an image", http.StatusBadRequest)
//        return
//    }
