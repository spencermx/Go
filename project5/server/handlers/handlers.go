package handlers

import (
    "encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"golang.org/x/time/rate"

	"github.com/google/uuid"
	"goserver/awsservices"
	"goserver/common"
	//	"github.com/gorilla/handlers"
	// "html/template"
)

// GLOBAL VARIABLES
var (
	AWS_REGION    string = "us-east-2"
	MAX_FILE_SIZE int64  = 150 << 20 // 150MB
    MAX_AUDIO_DURATION time.Duration = 60 * time.Minute
)

// Create a rate uploadLimiter with a maximum of 10 requests per minute
//var uploadLimiter = rate.NewLimiter(rate.Every(time.Minute), 2)
var uploadLimiter = rate.NewLimiter(rate.Every(time.Hour/10), 7)
var downloadLimiter = rate.NewLimiter(rate.Every(time.Hour/10), 150)

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

	// Check if the request is allowed by the rate uploadLimiter
	if !downloadLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
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

	// Check if the request is allowed by the rate uploadLimiter
	if !uploadLimiter.Allow() {
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

	var multipartFile *MultipartFile = &MultipartFile{File: file}

	if !multipartFile.IsImage() {
		http.Error(w, "The selected file must be an image", http.StatusBadRequest)
		return
	}

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

	var uploader *s3manager.Uploader = s3manager.NewUploader(sess)

	bucketName := os.Getenv("BUCKET_NAME")

	var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, nil, uploader, bucketName, AWS_REGION)

	bucketKey := common.BucketKey{
		Key: uuid.New().String() + "-" + header.Filename,
	}

	err = awsS3.UploadFile(bucketKey, file)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // Redirect to the "/home" route after the upload is complete
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateCaptionsVtt(w http.ResponseWriter, r *http.Request) {
	// Check if the request is allowed by the rate uploadLimiter
	if !downloadLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Parse the query string parameters
	queryParams := r.URL.Query()

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the value of the desired query parameter
	bucketKey := common.BucketKey{Key: queryParams.Get("bucketkey")}
	bucketName := os.Getenv("BUCKET_NAME")

	s3Client := s3.New(sess)
	log.Println("Created S3 client")

	var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, s3Client, nil, bucketName, AWS_REGION)

	containsKey := awsS3.ContainsKey(bucketKey.Key)

	if containsKey {
		log.Printf("Creating Amazon Transcription Client")
		transcribeClient := transcribeservice.New(sess)
		log.Printf("Successfully Created Amazon Transcription Client")

		var awsTranscribe *awsservices.AwsTranscribe = awsservices.NewAwsTranscribe(awsS3, transcribeClient)

        awsTranscribe.CreateVttFile(bucketKey)
	}
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the request is allowed by the rate uploadLimiter
	if !uploadLimiter.Allow() {
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
    // Validate the uploaded file is actually a video
	var multipartFile *MultipartFile = &MultipartFile{File: file}

	if !multipartFile.IsVideo() {
		http.Error(w, "The selected file must be a video", http.StatusBadRequest)
		return
	}
	/**************************************************************************************************/

	/**************************************************************************************************/
    // Validate the length of the audio in the uploaded video is less than the MAX_AUDIO_DURATION 
    duration, err := multipartFile.GetAudioDuration()

    log.Println("selected video audio duration: " + duration.String())
    if duration > MAX_AUDIO_DURATION {
        http.Error(w, "Video is too large for transcription. Maximum size is 1 hour", http.StatusBadRequest)
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

	bucketName := os.Getenv("BUCKET_NAME")
	bucketKey := common.BucketKey{Key: uuid.New().String() + "-" + header.Filename}

	var uploader *s3manager.Uploader = s3manager.NewUploader(sess)
	var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, nil, uploader, bucketName, AWS_REGION)

	err = awsS3.UploadFile(bucketKey, file)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/************************* Enable video transcription with AWS Transcribe *************************/
    // Enable video transcription with AWS Transcribe
	// fmt.Fprintf(w, "File uploaded successfully")

	//log.Printf("Creating Amazon Transcription Client")
	//transcribeClient := transcribeservice.New(sess)
	//log.Printf("Successfully Created Amazon Transcription Client")

	//var awsTranscribe *awsservices.AwsTranscribe = awsservices.NewAwsTranscribe(awsS3, transcribeClient)

	//awsTranscriptionJob, err := awsTranscribe.StartTranscriptionJob(bucketKey)
	//if err != nil {
	//	log.Printf("Error in transcription process: %v", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//err = awsTranscriptionJob.WaitForCompletion()
	//if err != nil {
	//	log.Printf("Error in transcription process: %v", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	// Create VTT file from json
	//awsTranscribe.CreateVttFile(bucketKey)
    /**************************************************************************************************/

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for /getImages")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the request is allowed by the rate uploadLimiter
	if !downloadLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
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

	var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, s3Client, nil, bucketName, AWS_REGION)

	cloudFrontUrls, err := awsS3.GetCloudFrontUrls()

	var clientItems []common.ClientItem

	for _, cloudFrontUrl := range cloudFrontUrls {
		if cloudFrontUrl.BucketKey.IsImage() {
			clientItem := &common.ClientItem{
				CloudFrontUrl: cloudFrontUrl.GetUrl(),
				FileName:      cloudFrontUrl.BucketKey.GetFileNameWithoutExtension(),
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

	// Check if the request is allowed by the rate uploadLimiter
	if !downloadLimiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
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

	var awsS3 *awsservices.AwsS3 = awsservices.NewAwsS3(sess, s3Client, nil, bucketName, AWS_REGION)

	cloudFrontUrls, err := awsS3.GetCloudFrontUrls()

	var clientItems []common.ClientItem

	for _, cloudFrontUrl := range cloudFrontUrls {
		if cloudFrontUrl.BucketKey.IsVideo() {
			clientItem := &common.ClientItem{
				CloudFrontUrl: cloudFrontUrl.GetUrl(),
				FileName:      cloudFrontUrl.BucketKey.GetFileNameWithoutExtensionAndGuid(),
                VideoCaptionsUrl: cloudFrontUrl.GetUrlCaptionsVtt(),
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
