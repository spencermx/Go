package awsservices

import (
	"log"
    "io"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

    "goserver/common"
	//	"github.com/gorilla/handlers"
	// "html/template"
)

// GLOBAL VARIABLES
var (
	AWS_REGION    string = "us-east-2"
	MAX_FILE_SIZE int64  = 150 << 20 // 150MB
)

type AwsS3 struct {
	AwsSession *session.Session
	S3Client   *s3.S3
	BucketName string
	Region     string
}

func NewAwsS3(awsSession *session.Session, s3Client *s3.S3, bucketName string, region string) *AwsS3 {
	return &AwsS3{
		AwsSession: awsSession,
        S3Client: s3Client,
		BucketName: bucketName,
		Region:     region,
	}
}

// return strings like "account/spencer/temp.png" "first.png"
func (s *AwsS3) GetKeys() ([]*common.BucketKey, error) {
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

func (s *AwsS3) GetFileByKey(bucketKey *common.BucketKey) ([]byte, error) {
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

func (s *AwsS3) GetCloudFrontUrls() ([]*common.FileUrl, error) {
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
