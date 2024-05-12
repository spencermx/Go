package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	//"github.com/aws/aws-sdk-go/service/transcribeservice"
	"golang.org/x/time/rate"

	"goserver/awsservices"
	"goserver/common"
	"github.com/google/uuid"
	//	"github.com/gorilla/handlers"
	//
	// "html/template"
)

// GLOBAL VARIABLES
var AWS_REGION string = "us-east-2"
var MAX_FILE_SIZE int64 = 150 << 20 // 150MB

var (
	// Create a rate limiter with a maximum of 10 requests per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute), 2)
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}


func GetPeople(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
	// Initialize the in-memory data structure with some sample data
	people := []Person{
		{ID: 1, Name: "John Doe", Age: 30},
		{ID: 2, Name: "Jane Smith", Age: 25},
		{ID: 3, Name: "Bob Johnson", Age: 35},
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the `people` slice as JSON and write it to the response
	json.NewEncoder(w).Encode(people)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
  	log.Printf("/uploadImage")
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// Check if the request is allowed by the rate limiter
	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(MAX_FILE_SIZE) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

    /**************************************************************************************************/
    // First Image Validation Check, may become obsolete after fileheader checks
    // var fileName *File = &File { Key: header.Filename } 

	// if (err != nil || !fileName.IsImage()) {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
    /**************************************************************************************************/

    /**************************************************************************************************/
    // NEW IMAGE VALIDATION
    var multipartFile *MultipartFile = &MultipartFile { File: file }

    if !multipartFile.IsImage() {
        http.Error(w, "The selected file must be an image", http.StatusBadRequest)
		return
    }
    /**************************************************************************************************/

    /**************************************************************************************************/

	defer file.Close()

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})

	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Created AWS session")

	// Create an S3 uploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(header.Filename),
		Body:   file,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully")
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
	// Check if the request is allowed by the rate limiter
	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(MAX_FILE_SIZE)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, header, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    /**************************************************************************************************/
    // NEW VIDEO VALIDATION
    var multipartFile *MultipartFile = &MultipartFile { File: file }

    if !multipartFile.IsVideo() {
        http.Error(w, "The selected file must be a video", http.StatusBadRequest)
		return
    }
    /**************************************************************************************************/

	defer file.Close()

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})

	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Created AWS session")

	// Create an S3 uploader
	uploader := s3manager.NewUploader(sess)

    var bucketName string = os.Getenv("BUCKET_NAME")
    
    uuid, _ := uuid.NewRandom()

    //var uuid string = "9e1e2dd4-c836-43af-ba21-090b9a1032d3"
    key := fmt.Sprintf("%s-%s", uuid, header.Filename)

	// Upload the file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully")

    // Create a new Amazon Transcribe client
    // log.Printf("Creating Amazon Transcription Client")

    // transcribeClient := transcribeservice.New(sess)

    // log.Printf("Successfully Created Amazon Transcription Client")
    // 
    // var videoS3Uri string = fmt.Sprintf("s3://%s/%s", bucketName, key)
    // var videoTranscriptOutput string = fmt.Sprintf("%s-transcription-output.json", uuid)
    // 
    // transcriptionJobInput := &transcribeservice.StartTranscriptionJobInput{
    //     Media: &transcribeservice.Media{
    //         MediaFileUri: aws.String(videoS3Uri),
    //     },
    //     OutputBucketName: aws.String(bucketName),
    //     OutputKey:        aws.String(videoTranscriptOutput),
    //     TranscriptionJobName:  aws.String("TranscriptionJobName-" + uuid), 
    //     LanguageCode:     aws.String("en-US"), // Set the language code
    //     // Set any other necessary options
    // }

    // _, err = transcribeClient.StartTranscriptionJob(transcriptionJobInput)

    // if err != nil {
	// 	log.Printf("Transcription failure: %v", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
    // }


    // Save the transcription job name for later use

//    transcriptionjobname := transcriptionjoboutput.transcriptionjob.transcriptionjobname
//   // wait for the transcription job to complete
//    for {
//        jobstatusoutput, err := transcribeclient.gettranscriptionjob(&transcribeservice.gettranscriptionjobinput{
//            transcriptionjobname: transcriptionjobname,
//        })
//        if err != nil {
//            log.printf("failed to get transcription job status: %v", err)
//            http.error(w, err.error(), http.statusinternalservererror)
//            return
//        }
//
//        jobstatus := *jobstatusoutput.transcriptionjob.transcriptionjobstatus
//
//        if jobstatus == transcribeservice.transcriptionjobstatuscompleted {
//            log.printf("transcription job %s completed successfully", *transcriptionjobname)
//            break
//        } else if jobstatus == transcribeservice.transcriptionjobstatusfailed {
//            log.printf("transcription job %s failed with status %s", *transcriptionjobname, jobstatus)
//            http.error(w, fmt.sprintf("transcription job failed with status %s", jobstatus), http.statusinternalservererror)
//            return
//        } else {
//            log.printf("transcription job %s is in progress with status %s", *transcriptionjobname, jobstatus)
//            time.sleep(5 * time.second) // wait for a few seconds before checking again
//        }
//    }
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for /getImages")

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Created AWS session")

	// Create an S3 client
	s3Client := s3.New(sess)
	log.Println("Created S3 client")

    var bucketName string = os.Getenv("BUCKET_NAME")
    
    var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, s3Client, bucketName, AWS_REGION)  

    cloudFrontUrls, err := awsS3.GetCloudFrontUrls()

    var clientItems []common.ClientItem

    for _, cloudFrontUrl := range cloudFrontUrls {
        if cloudFrontUrl.BucketKey.IsImage() {
            clientItem := &common.ClientItem { 
                CloudFrontUrl: cloudFrontUrl.GetUrl(),
                FileName: cloudFrontUrl.BucketKey.GetFileNameWithoutExtension(),
            }    

            clientItems = append(clientItems, *clientItem)
        }
    }

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Encode the image objects as JSON and write the response
	err = json.NewEncoder(w).Encode(clientItems)
	if err != nil {
		log.Printf("Failed to encode image objects as JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for /getVideos")

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Created AWS session")

	// Create an S3 client
	s3Client := s3.New(sess)
	log.Println("Created S3 client")

    var bucketName string = os.Getenv("BUCKET_NAME")
    
    var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, s3Client, bucketName, AWS_REGION)  

    cloudFrontUrls, err := awsS3.GetCloudFrontUrls()

    var clientItems []common.ClientItem

    for _, cloudFrontUrl := range cloudFrontUrls {
        if cloudFrontUrl.BucketKey.IsVideo() {
            clientItem := &common.ClientItem { 
                CloudFrontUrl: cloudFrontUrl.GetUrl(),
                FileName: cloudFrontUrl.BucketKey.GetFileNameWithoutExtension(),
            }    

            clientItems = append(clientItems, *clientItem)
        }
    }

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Encode the image objects as JSON and write the response
	err = json.NewEncoder(w).Encode(clientItems)
	if err != nil {
		log.Printf("Failed to encode image objects as JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
