package awsservices

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
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
                Confidence float64 `json:"confidence"`
            } `json:"alternatives"`
            StartTime string `json:"start_time"`
            EndTime   string `json:"end_time"`
        } `json:"items"`
    } `json:"results"`
}

func testinggg() {
    // Read the JSON file
    jsonFile, err := os.Open("transcription.json")
    if err != nil {
        fmt.Println("Error opening JSON file:", err)
        return
    }
    defer jsonFile.Close()

    // Parse the JSON data
    jsonData, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        fmt.Println("Error reading JSON data:", err)
        return
    }

    var result TranscriptionResult
    err = json.Unmarshal(jsonData, &result)
    if err != nil {
        fmt.Println("Error parsing JSON data:", err)
        return
    }

    // Create the VTT file
    vttFile, err := os.Create("captions.vtt")
    if err != nil {
        fmt.Println("Error creating VTT file:", err)
        return
    }
    defer vttFile.Close()

    // Write the VTT header
    _, err = vttFile.WriteString("WEBVTT\n\n")
    if err != nil {
        fmt.Println("Error writing VTT header:", err)
        return
    }

    // Process each item and write to the VTT file
    for i, item := range result.Results.Items {
        if item.Type == "pronunciation" {
            startTime := formatTime(item.StartTime)
            endTime := formatTime(item.EndTime)
            content := item.Alternatives[0].Content

            // Write the caption entry to the VTT file
            captionEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, content)
            _, err := vttFile.WriteString(captionEntry)
            if err != nil {
                fmt.Println("Error writing caption entry:", err)
                return
            }
        }
    }

    fmt.Println("VTT file created successfully.")
}

func formatTime(timeStr string) string {
    seconds, _ := strconv.ParseFloat(timeStr, 64)
    minutes := int(seconds) / 60
    seconds -= float64(minutes * 60)
    return fmt.Sprintf("%02d:%02d:%06.3f", minutes/60, minutes%60, seconds)
}
// package awsservices
// 
// import (
//     "encoding/json"
//     "fmt"
//     "io"
//     "os"
// )
// 
// type TranscribeResult struct {
//     Results struct {
//         Items []struct {
//             Content  string  `json:"content"`
//             StartTime float64 `json:"start_time"`
//             EndTime   float64 `json:"end_time"`
//         } `json:"items"`
//     } `json:"results"`
// }
// 
// func testing() {
//     // Read the JSON file
//     jsonFile, err := os.Open("transcribe_result.json")
//     if err != nil {
//         fmt.Println("Error opening JSON file:", err)
//         return
//     }
// 
//     defer jsonFile.Close()
// 
//     // Parse the JSON data
//     var result TranscribeResult
//     jsonData, _ := io.ReadAll(jsonFile)
//     json.Unmarshal(jsonData, &result)
// 
//     // Create a new WebVTT file
//     vttFile, err := os.Create("captions.vtt")
//     if err != nil {
//         fmt.Println("Error creating VTT file:", err)
//         return
//     }
//     defer vttFile.Close()
// 
//     // Write the WebVTT header
//     vttFile.WriteString("WEBVTT\n\n")
// 
//     // Iterate over the parsed data and write captions to the VTT file
//     for _, item := range result.Results.Items {
//         startTime := formatTime(item.StartTime)
//         endTime := formatTime(item.EndTime)
//         caption := item.Content
// 
//         vttFile.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
//         vttFile.WriteString(caption + "\n\n")
//     }
// 
//     fmt.Println("WebVTT file generated successfully.")
// }
// 
// func formatTime(seconds float64) string {
//     // Format the time in the required format (e.g., "00:00:00.000")
//     // You can use the `time` package or write your own formatting logic
//     // This is just a placeholder example
//     return fmt.Sprintf("%02d:%02d:%06.3f", int(seconds)/3600, int(seconds)%3600/60, seconds)
// }
