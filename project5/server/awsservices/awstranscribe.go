package awsservices

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type TranscribeResult struct {
    Results struct {
        Items []struct {
            Content  string  `json:"content"`
            StartTime float64 `json:"start_time"`
            EndTime   float64 `json:"end_time"`
        } `json:"items"`
    } `json:"results"`
}

func testing() {
    // Read the JSON file
    jsonFile, err := os.Open("transcribe_result.json")
    if err != nil {
        fmt.Println("Error opening JSON file:", err)
        return
    }

    defer jsonFile.Close()

    // Parse the JSON data
    var result TranscribeResult
    jsonData, _ := io.ReadAll(jsonFile)
    json.Unmarshal(jsonData, &result)

    // Create a new WebVTT file
    vttFile, err := os.Create("captions.vtt")
    if err != nil {
        fmt.Println("Error creating VTT file:", err)
        return
    }
    defer vttFile.Close()

    // Write the WebVTT header
    vttFile.WriteString("WEBVTT\n\n")

    // Iterate over the parsed data and write captions to the VTT file
    for _, item := range result.Results.Items {
        startTime := formatTime(item.StartTime)
        endTime := formatTime(item.EndTime)
        caption := item.Content

        vttFile.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
        vttFile.WriteString(caption + "\n\n")
    }

    fmt.Println("WebVTT file generated successfully.")
}

func formatTime(seconds float64) string {
    // Format the time in the required format (e.g., "00:00:00.000")
    // You can use the `time` package or write your own formatting logic
    // This is just a placeholder example
    return fmt.Sprintf("%02d:%02d:%06.3f", int(seconds)/3600, int(seconds)%3600/60, seconds)
}
