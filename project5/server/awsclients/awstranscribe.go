package awsclients

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"goserver/common"
)

type AwsClientTranscribe struct {
	AwsS3            *AwsClientS3
	TranscribeClient *transcribeservice.TranscribeService
}

func NewAwsClientTranscribe(awsS3 *AwsClientS3, transcribeClient *transcribeservice.TranscribeService) *AwsClientTranscribe {
	awsClientTranscribe := &AwsClientTranscribe{AwsS3: awsS3, TranscribeClient: transcribeClient}

	return awsClientTranscribe
}

// *transcribeservice.StartTranscriptionJobOutput
func (t *AwsClientTranscribe) StartTranscriptionJob(bucketKey common.BucketKey) (*AwsTranscriptionJob, error) {
	var videoS3Uri string = fmt.Sprintf("s3://%s/%s", t.AwsS3.BucketName, bucketKey.Key)
	var videoTranscriptOutput string = bucketKey.GetKeyForTranscription()
	var transcriptionJobName string = bucketKey.GetTranscriptionJobName()

	transcriptionJobInput := &transcribeservice.StartTranscriptionJobInput{
		Media: &transcribeservice.Media{
			MediaFileUri: aws.String(videoS3Uri),
		},
		OutputBucketName:     aws.String(t.AwsS3.BucketName),
		OutputKey:            aws.String(videoTranscriptOutput),
		TranscriptionJobName: aws.String(transcriptionJobName),
		LanguageCode:         aws.String("en-US"), // Set the language code
		// Set any other necessary options
	}

	startTranscriptionJobOutput, err := t.TranscribeClient.StartTranscriptionJob(transcriptionJobInput)
	if err != nil {
		return nil, errors.New("Failed to start transcription job")
	}

	awsTranscriptionJob := &AwsTranscriptionJob{
		BucketKey:                   bucketKey,
		StartTranscriptionJobInput:  transcriptionJobInput,
		StartTranscriptionJobOutput: startTranscriptionJobOutput,
		TranscribeClient:            t.TranscribeClient,
	}

	return awsTranscriptionJob, nil
}

func (t *AwsClientTranscribe) CreateVttFile(bucketKey common.BucketKey) error {
    transcriptionResult, err := t.AwsS3.GetTranscriptionJson(bucketKey)

	if err != nil {
        return err
	}

    vttBuffer := transcriptionResult.CreateCaptionsVtt()

    t.AwsS3.UploadFileCaptions(bucketKey, vttBuffer)

    return nil
}

