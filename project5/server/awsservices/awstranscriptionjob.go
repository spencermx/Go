package awsservices

import (
	"errors"
	"log"
	"time"

	"goserver/common"

	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

type AwsTranscriptionJob struct {
	BucketKey                   common.BucketKey
	StartTranscriptionJobInput  *transcribeservice.StartTranscriptionJobInput
	StartTranscriptionJobOutput *transcribeservice.StartTranscriptionJobOutput
	TranscribeClient            *transcribeservice.TranscribeService
}

func NewAwsTranscriptionJob(awsS3 *AwsS3, transcribeClient *transcribeservice.TranscribeService) *AwsTranscribe {
	awsTranscribe := &AwsTranscribe{AwsS3: awsS3, TranscribeClient: transcribeClient}

	return awsTranscribe
}

func (t *AwsTranscriptionJob) WaitForCompletion() error {
	// Save the transcription job name for later use
	transcriptionjobname := t.StartTranscriptionJobOutput.TranscriptionJob.TranscriptionJobName
	// wait for the transcription job to complete

	for {
		jobStatusOutput, err := t.TranscribeClient.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
			TranscriptionJobName: transcriptionjobname,
		})
		if err != nil {
			log.Printf("failed to get transcription job status: %v", err)
			return errors.New("failed to query transcription job status")
		}

		jobstatus := *jobStatusOutput.TranscriptionJob.TranscriptionJobStatus

		if jobstatus == transcribeservice.TranscriptionJobStatusCompleted {
			log.Printf("transcription job %s completed successfully", *transcriptionjobname)
			break
		} else if jobstatus == transcribeservice.TranscriptionJobStatusFailed {
			log.Printf("transcription job %s failed with status %s", *transcriptionjobname, jobstatus)
			return errors.New("transcription job failed")
		} else {
			log.Printf("transcription job %s is in progress with status %s", *transcriptionjobname, jobstatus)
			time.Sleep(5 * time.Second) // wait for a few seconds before checking again
		}
	}
	return nil
}
