package awsservices

import (
    "bytes"
    "fmt"
    "strconv"
)

type TranscriptionResult struct {
    JobName   string `json:"jobName"`
    AccountID string `json:"accountId"`
    Results   struct {
        Transcripts []struct {
            Transcript string `json:"transcript"`
        } `json:"transcripts"`
        Items []struct {
            Type         string `json:"type"`
            Alternatives []struct {
                Content string  `json:"content"`
                Confidence string `json:"confidence"`
            } `json:"alternatives"`
            StartTime string `json:"start_time"`
            EndTime   string `json:"end_time"`
        } `json:"items"`
    } `json:"results"`
}

func (r TranscriptionResult) CreateCaptionsVtt() *bytes.Reader {
    // Create a new buffer to store the VTT content
    buffer := new(bytes.Buffer)

    // Write the VTT file header to the buffer
    buffer.WriteString("WEBVTT\n\n")

    // Iterate over the transcription items and write them to the buffer
    for i, item := range r.Results.Items {
        startTime := formatTime(item.StartTime)
        endTime := formatTime(item.EndTime)
        content := item.Alternatives[0].Content

        // Write the caption entry to the buffer
        captionEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, content)
        buffer.WriteString(captionEntry)
    }

    return bytes.NewReader(buffer.Bytes())
}

func formatTime(timeStr string) string {
    seconds, _ := strconv.ParseFloat(timeStr, 64)
    minutes := int(seconds) / 60
    seconds -= float64(minutes * 60)
    return fmt.Sprintf("%02d:%02d:%06.3f", minutes/60, minutes%60, seconds)
}
