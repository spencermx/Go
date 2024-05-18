package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func (m *MultipartFile) GenerateThumbnail() ([]byte, error) {
    videoDuration, err := m.GetVideoDuration()

	middleDuration := videoDuration / 2

    timeStamp := formatDuration(middleDuration)

    fileBytes, err := m.getBytesFromFile()

	// Create an FFmpeg command to generate the thumbnail
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",              // Read input from stdin
		"-ss", timeStamp,            // "00:00:01" Set the thumbnail timestamp (adjust as needed)
		"-vframes", "1",             // Extract a single frame
		// "-vf", "scale=320:240",      // Set the thumbnail dimensions (adjust as needed)
		"-f", "image2pipe",          // Output to stdout as an image
		"-c:v", "mjpeg",             // Use MJPEG codec for the thumbnail
		"pipe:1",                    // Write output to stdout
	)

	// Set the command's stdin and stdout
	cmd.Stdin = &fileBytes
	var thumbnailBytes bytes.Buffer
	cmd.Stdout = &thumbnailBytes

	// Run the FFmpeg command
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to generate thumbnail: %v", err)
	}

	return thumbnailBytes.Bytes(), nil
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

func (m *MultipartFile) GetVideoDuration() (time.Duration, error) {
    // Save the current position of the file
    currentPosition, err := m.File.Seek(0, io.SeekCurrent)
    if err != nil {
        return 0, err
    }

    // Reset the file position after getting the duration
    defer m.File.Seek(currentPosition, io.SeekStart)

    /********************************************************* GET VIDEO **************************************************************/
    cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-")
    
    cmd.Stdin = m.File

    var out bytes.Buffer

    cmd.Stdout = &out
    
    err = cmd.Run()
    if err != nil {
        return 0, err
    }

    // Parse the duration from the output buffer
    durationStr := strings.TrimSpace(out.String())
    duration, err := strconv.ParseFloat(durationStr, 64)
    if err != nil {
        return 0, err
    }

    // Convert the duration to time.Duration
    videoDuration := time.Duration(duration * float64(time.Second))

    return videoDuration, nil
}

func  formatDuration(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d.%03d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
		int(duration.Milliseconds())%1000,
	)
}

func readFileToBytes(f multipart.File) ([]byte, error) {
    var buf bytes.Buffer
    _, err := io.Copy(&buf, f)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
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

func (m *MultipartFile) getBytesFromFile() (bytes.Buffer, error) {
	// Read the file into memory
    var fileBytes bytes.Buffer

	_, err := io.Copy(&fileBytes, m.File)
	if err != nil {
		return fileBytes, err 
	}

    return fileBytes, nil
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
