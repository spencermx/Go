package handlers

import (
    "log"
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
    var thumbnailBytes bytes.Buffer

    log.Println("Generating thumbnail...")

    videoDuration, err := m.GetVideoDuration()
    if err != nil {
        log.Printf("Failed to get video duration: %v", err)
        return thumbnailBytes.Bytes(), err
    }

    log.Printf("Video duration: %v", videoDuration)

    middleDuration := videoDuration / 2
    timeStamp := formatDuration(middleDuration)

    log.Printf("Thumbnail timestamp: %s", timeStamp)

    fileBytes, err := m.getBytesFromFile()
    if err != nil {
        log.Printf("Failed to get bytes from file: %v", err)
        return nil, err
    }

    log.Println("Creating FFmpeg command...")

    // Create an FFmpeg command to generate the thumbnail
    cmd := exec.Command("ffmpeg",
        "-i", "pipe:0", // Read input from stdin
        "-ss", timeStamp, // "00:00:01" Set the thumbnail timestamp (adjust as needed)
        "-vframes", "1", // Extract a single frame
        "-vf", "scale=600:400", // Set the thumbnail dimensions (adjust as needed)
        "-f", "image2pipe", // Output to stdout as an image
        "-c:v", "mjpeg", // Use MJPEG codec for the thumbnail
        "pipe:1", // Write output to stdout
    )

    log.Println("Setting command's stdin and stdout...")

    // Set the command's stdin and stdout
    cmd.Stdin = &fileBytes
    cmd.Stdout = &thumbnailBytes

    log.Println("Running FFmpeg command...")

    // Run the FFmpeg command
    err = cmd.Run()
    if err != nil {
        log.Printf("Failed to generate thumbnail: %v", err)
        return nil, fmt.Errorf("failed to generate thumbnail: %v", err)
    }

    log.Println("Thumbnail generated successfully")

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
    // Save the current position of the file
    currentPosition, err := m.File.Seek(0, io.SeekCurrent)
    if err != nil {
        log.Printf("Failed to get current file position: %v", err)
        return bytes.Buffer{}, err
    }

    log.Printf("Current file position: %d", currentPosition)

    // Read the file into memory
    var fileBytes bytes.Buffer
    _, err = io.Copy(&fileBytes, m.File)
    if err != nil {
        log.Printf("Failed to read file into memory: %v", err)
        return bytes.Buffer{}, err
    }

    log.Printf("File read into memory. Size: %d bytes", fileBytes.Len())

    // Reset the file cursor position to the beginning
    _, err = m.File.Seek(currentPosition, io.SeekStart)
    if err != nil {
        log.Printf("Failed to reset file position: %v", err)
        return bytes.Buffer{}, err
    }

    log.Printf("File position reset to: %d", currentPosition)

    return fileBytes, nil
}

//func (m *MultipartFile) getBytesFromFile() (bytes.Buffer, error) {
//	// Read the file into memory
//    var fileBytes bytes.Buffer
//
//	_, err := io.Copy(&fileBytes, m.File)
//	if err != nil {
//		return fileBytes, err 
//	}
//
//    return fileBytes, nil
//}

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
