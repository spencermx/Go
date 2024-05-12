package awsservices

import (
    "strings"
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

const maxTimeDiff = 3.0 // Maximum time difference in seconds between consecutive words to group them together

func (r TranscriptionResult) CreateCaptionsVtt() *bytes.Reader {
    buffer := new(bytes.Buffer)
    buffer.WriteString("WEBVTT\n\n")

    var currentCaptionWords []string

    var originalStartTime string
    var originalStartTimeValue float64
    var endTimeOfLastCaptionAdded string
    var captionIndex int = 1
    
    for i, item := range r.Results.Items {
        currentStartTime := formatTime(item.StartTime)
        currentStartTimeValue, _ := strconv.ParseFloat(item.StartTime, 64)
        currentEndTime := formatTime(item.EndTime)
        currentEndTimeValue, _ := strconv.ParseFloat(item.EndTime, 64)
        currentWord := item.Alternatives[0].Content

        if (i == 0) {
            // first iteration set original
            originalStartTime = currentStartTime
            originalStartTimeValue = currentStartTimeValue
        }

        if (item.Type == "punctuation") {
            currentCaptionWords = append(currentCaptionWords, currentWord)
            continue
        }
        
        if (i == len(r.Results.Items)-1) {
            // last iteration
            currentCaptionWords = append(currentCaptionWords, currentWord)
            captionEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", captionIndex, originalStartTime, currentEndTime, joinWordsWithPunctuation(currentCaptionWords))
            buffer.WriteString(captionEntry)
        }

        if (currentEndTimeValue - originalStartTimeValue > maxTimeDiff) {
            // write the current caption then start a new originalStartTimeValue
            captionEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", captionIndex, originalStartTime, endTimeOfLastCaptionAdded, joinWordsWithPunctuation(currentCaptionWords))
            buffer.WriteString(captionEntry)
            captionIndex += 1

            // start a new caption
            currentCaptionWords = []string{currentWord}
            originalStartTime = currentStartTime
            originalStartTimeValue = currentStartTimeValue
            endTimeOfLastCaptionAdded = currentEndTime
        } else {
            // continue adding to the current caption
            currentCaptionWords = append(currentCaptionWords, currentWord)
            endTimeOfLastCaptionAdded = currentEndTime 
        }
    }

    return bytes.NewReader(buffer.Bytes())
}

func formatTime(timeStr string) string {
    seconds, _ := strconv.ParseFloat(timeStr, 64)
    minutes := int(seconds) / 60
    seconds -= float64(minutes * 60)
    return fmt.Sprintf("%02d:%02d:%06.3f", minutes/60, minutes%60, seconds)
}

func joinWordsWithPunctuation(words []string) string {
    var result strings.Builder
    for i := 0; i < len(words); i++ {
        result.WriteString(words[i])
        if i < len(words)-1 {
            if isPunctuation(words[i+1]) {
                // If the next word is punctuation, join directly
                continue
            } else {
                // If the next word is not punctuation, join with a space
                result.WriteString(" ")
            }
        }
    }
    return result.String()
}

func isPunctuation(word string) bool {
    // Define the set of punctuation characters
    punctuationChars := ",.!?:;-"
    return len(word) == 1 && strings.ContainsAny(word, punctuationChars)
}

//func (r TranscriptionResult) CreateCaptionsVtt() *bytes.Reader {
//    // Create a new buffer to store the VTT content
//    buffer := new(bytes.Buffer)
//
//    // Write the VTT file header to the buffer
//    buffer.WriteString("WEBVTT\n\n")
//
//    // Iterate over the transcription items and write them to the buffer
//    for i, item := range r.Results.Items {
//        startTime := formatTime(item.StartTime)
//        endTime := formatTime(item.EndTime)
//        content := item.Alternatives[0].Content
//
//        // Write the caption entry to the buffer
//        captionEntry := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, content)
//        buffer.WriteString(captionEntry)
//    }
//
//    return bytes.NewReader(buffer.Bytes())
//}

