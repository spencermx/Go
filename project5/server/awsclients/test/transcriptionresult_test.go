package awsclients_test

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
	"time"
	"goserver/awsclients"
	"goserver/common"
)


// GLOBAL VARIABLES
var (
	AWS_REGION    string = "us-east-2"
	MAX_FILE_SIZE int64  = 150 << 20 // 150MB
    MAX_AUDIO_DURATION time.Duration = 60 * time.Minute
)


func TestSerializeS3TranscriptionJson(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
    //bucketKey := common.BucketKey{Key: "dde20b10-8559-4be5-9d11-7538b24788ab-transcription-output.json" }
    //bucketKey := common.BucketKey{Key: "2143d451-7b44-4bbb-96c9-0ec4c50fcc43-transcription-output.json" }
	bucketKey := common.BucketKey{Key:"0ef5f6b2-2fb8-4d37-9635-ab325d115072-transcription-output.json" } 

    awsS3, err := awsclients.NewAwsClientS3(bucketName, AWS_REGION)

    assert.NotNil(t, err, "Expected an error, but got nil")

    transcriptionResult, err := awsS3.GetTranscriptionJson(bucketKey)

    transcriptionResult.CreateCaptionsVtt()
}

//func TestTranscriptionResultSerialization(t *testing.T) {
//    // Read the JSON file
//    jsonFile, err := os.ReadFile(filepath.Join("../../testdata/", "transcription_result.json"))
//    if err != nil {
//        t.Fatalf("Error reading JSON file: %v", err)
//    }
//
//    // Unmarshal the JSON data into the TranscriptionResult struct
//    var result awsclients.TranscriptionResult
//    err = json.Unmarshal(jsonFile, &result)
//    if err != nil {
//        t.Fatalf("Error unmarshaling JSON: %v", err)
//    }
//
//    // Verify the serialized data
//    expectedJobName := "TranscriptionJobName-dde20b10-8559-4be5-9d11-7538b24788ab"
//    if result.JobName != expectedJobName {
//        t.Errorf("Expected JobName: %s, got: %s", expectedJobName, result.JobName)
//    }
//
//    // Add more assertions for other fields as needed
//
//    // Example assertion for transcript
//    expectedTranscript := "Let's count the 1st 11 patterns."
//    if len(result.Results.Transcripts) == 0 || result.Results.Transcripts[0].Transcript != expectedTranscript {
//        t.Errorf("Expected transcript: %s, got: %s", expectedTranscript, result.Results.Transcripts[0].Transcript)
//    }
//}


//    getObjectInput := &s3.GetObjectInput{
//        Bucket: aws.String(bucketName),
//        Key:    aws.String(bucketKey.Key),
//    }
//
//    getObjectOutput, err := s.GetObject(getObjectInput)
//
//    if err != nil {
//    }
//
//    defer getObjectOutput.Body.Close()
//
//    // Read the JSON content into a byte buffer
//    buf := new(bytes.Buffer)
//
//    _, err = buf.ReadFrom(getObjectOutput.Body)
//
//    if err != nil {
//    }
//
//    // Parse the JSON content into the TranscriptionResult struct
//    var transcriptionResult awsclients.TranscriptionResult
//    fmt.Println(buf.String())
//    err = json.Unmarshal(buf.Bytes(), &transcriptionResult)
//
//    if err != nil {
//    }
