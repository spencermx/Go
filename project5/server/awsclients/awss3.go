package awsclients

import (
    "fmt"
    "bytes"
    "encoding/json"
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"

	"errors"
	"goserver/common"
	"strings"
)

// GLOBAL VARIABLES
var (
	AWS_REGION    string = "us-east-2"
	MAX_FILE_SIZE int64  = 150 << 20 // 150MB
)

type AwsClientS3 struct {
	AwsSession *session.Session
	S3Client   *s3.S3
    S3Uploader  *s3manager.Uploader
	BucketName string
	Region     string
}

func NewAwsClientS3(awsBucketName string, awsRegion string) (*AwsClientS3, error) {
	// Create a new AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
        return nil, err
	}

	log.Println("Created AWS session")

	var uploader *s3manager.Uploader = s3manager.NewUploader(awsSession)

	// Create an S3 client
	s3Client := s3.New(awsSession)

    awsClientS3 := &AwsClientS3{
		AwsSession: awsSession,
        S3Client: s3Client,
        S3Uploader: uploader,
		BucketName: awsBucketName,
		Region:     awsRegion,
	}

    return awsClientS3, nil
}

func (s *AwsClientS3) ContainsKey(key string) bool {
    keys, _ := s.GetKeys()

    for _, item := range keys {
        if item.Key == key {
            return true
        }
    }

    return false
}

// return strings like "account/spencer/temp.png" "first.png"
func (s *AwsClientS3) GetKeys() ([]*common.BucketKey, error) {
	// List objects in the S3 bucket
	result, err := s.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		return nil, err
	}

	// Extract the object keys from the result
	keys := make([]*common.BucketKey, len(result.Contents))
	for i, obj := range result.Contents {
        keys[i] = &common.BucketKey{
            Key: *obj.Key,
        }
	}

	return keys, nil
}

func (s *AwsClientS3) GetFileByKey(bucketKey *common.BucketKey) ([]byte, error) {
	// Get the object from the S3 bucket
	result, err := s.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(bucketKey.Key),
	})

	if err != nil {
		log.Printf("Failed to get object from S3 bucket: %v", err)
		return nil, err
	}

	defer result.Body.Close()

	// Read the object contents
	fileBytes, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Failed to read object contents: %v", err)
		return nil, err
	}

	return fileBytes, nil
}

func (s *AwsClientS3) GetCloudFrontUrls() ([]*common.FileUrl, error) {
    bucketKeys, err := s.GetKeys()

    if err != nil {
        return nil, err
    }

    var cloudFrontUrls []*common.FileUrl

    for _, bucketKey := range bucketKeys {
        fileUrl := common.NewFileUrlWithDefaultRootUrl(*bucketKey)

        cloudFrontUrls = append(cloudFrontUrls, fileUrl)
    }

    return cloudFrontUrls, nil
}

func (s *AwsClientS3) GetKeysByGuid(uuid uuid.UUID) ([]*common.BucketKey, error) {
    bucketKeys, err := s.GetKeys()

    if err != nil {
        return nil, err
    }

    var bucketsKeysForGuid []*common.BucketKey

    for _, bucketKey := range bucketKeys {
        if strings.Contains(bucketKey.Key, uuid.String()) {
            bucketsKeysForGuid = append(bucketsKeysForGuid, bucketKey)
        } 
    }

    if len(bucketsKeysForGuid) == 0 {
        return nil, errors.New("no bucket keys found containing that guid")
    }

    return bucketsKeysForGuid, nil
}

func (s *AwsClientS3) UploadFileCaptions(bucketKey common.BucketKey, buffer *bytes.Reader) error {
    uploader := s3manager.NewUploader(s.AwsSession)
    uploadInput := &s3manager.UploadInput{
        Bucket: aws.String(s.BucketName),
        Key:    aws.String(bucketKey.GetKeyForCaptions()),
        Body:   buffer,
    }
    _, err := uploader.Upload(uploadInput)

    if err != nil {
        return err
    }

    return nil
}

func (s *AwsClientS3) UploadFile(bucketKey common.BucketKey, file multipart.File) error {
    _, err := s.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(bucketKey.Key),
		Body:   file,
	})

	if err != nil {
        return err
	}

    return nil
}

func (s *AwsClientS3) GetTranscriptionJson(bucketKey common.BucketKey) (*TranscriptionResult, error) {
    getObjectInput := &s3.GetObjectInput{
        Bucket: aws.String(s.BucketName),
        Key:    aws.String(bucketKey.GetKeyForTranscription()),
    }
    getObjectOutput, err := s.S3Client.GetObject(getObjectInput)

    if err != nil {
        return nil, err
    }

    defer getObjectOutput.Body.Close()

    // Read the JSON content into a byte buffer
    buf := new(bytes.Buffer)

    _, err = buf.ReadFrom(getObjectOutput.Body)

    if err != nil {
       return nil, err
    }

    // Parse the JSON content into the TranscriptionResult struct
    var transcriptionResult TranscriptionResult
    fmt.Println(buf.String())
    err = json.Unmarshal(buf.Bytes(), &transcriptionResult)

    if err != nil {
       return nil, err
    }

    return &transcriptionResult, nil

}




















