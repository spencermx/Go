package handlers

import (
    "os/exec"
    "encoding/json"
    "strconv"
    "fmt"
    "time"
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

func (m *MultipartFile) GetAudioDuration() (time.Duration, error) {
    // Save the current position of the file
    currentPosition, err := m.File.Seek(0, io.SeekCurrent)
    if err != nil {
        return 0, err
    }

    // Reset the file position after getting the duration
    defer m.File.Seek(currentPosition, io.SeekStart)

    // Use ffprobe to get the audio duration
    cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_entries", "format=duration", "-select_streams", "a", "-")
    cmd.Stdin = m.File
    var out bytes.Buffer
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil {
        return 0, err
    }

    // Parse the JSON output
    var format struct {
        Format struct {
            Duration string `json:"duration"`
        } `json:"format"`
    }
    err = json.Unmarshal(out.Bytes(), &format)
    if err != nil {
        return 0, err
    }

    // Check if an audio duration was found
    if format.Format.Duration == "" {
        return 0, fmt.Errorf("no audio duration found")
    }

    // Convert duration to time.Duration
    durationSeconds, err := strconv.ParseFloat(format.Format.Duration, 64)
    if err != nil {
        return 0, err
    }

    duration := time.Duration(durationSeconds * float64(time.Second))

    return duration, nil
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
