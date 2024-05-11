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
	"golang.org/x/time/rate"
	//	"github.com/gorilla/handlers"
	//
	// "html/template"
)

// GLOBAL VARIABLES
var AWS_REGION string = "us-east-2"

var (
	// Create a rate limiter with a maximum of 10 requests per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute), 2)
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
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
	err := r.ParseMultipartForm(10 << 20) // Max size of 10MB

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
	err := r.ParseMultipartForm(10 << 20) // Max size of 10MB

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
func GetImages(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	log.Println("Handling request for /getImages")

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
	svc := s3.New(sess)
	log.Println("Created S3 client")

	// List objects in the S3 bucket
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
	})
	if err != nil {
		log.Printf("Failed to list objects in S3 bucket: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Listed %d objects in S3 bucket", len(result.Contents))

	// Create a slice to store the image objects
	var images []Image

	//  imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
	// Iterate over the objects and create the image objects
	for _, obj := range result.Contents {
        file := &File { Key: *obj.Key } 

        if file.IsImage() {
            imageURL := fmt.Sprintf("https://d271tjczb1hjof.cloudfront.net/%s", file.Key)
            image := Image{
                URL: imageURL,
                Alt: file.Key, // Use the object key as the alt text
            }
            images = append(images, image)
        }
	}
	log.Printf("Created %d image objects", len(images))

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Encode the image objects as JSON and write the response
	err = json.NewEncoder(w).Encode(images)
	if err != nil {
		log.Printf("Failed to encode image objects as JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    } 
	log.Println("Handling request for /getVideos")

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
	svc := s3.New(sess)
	log.Println("Created S3 client")

	// List objects in the S3 bucket
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
	})

	if err != nil {
		log.Printf("Failed to list objects in S3 bucket: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Listed %d objects in S3 bucket", len(result.Contents))

	// Create a slice to store the image objects
	var images []Image

	//  imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
	// Iterate over the objects and create the image objects
	for _, object := range result.Contents {
        file := &File { Key: *object.Key }

        if file.IsVideo() {
            imageURL := fmt.Sprintf("https://d271tjczb1hjof.cloudfront.net/%s", *object.Key)
            
            image := Image{
                URL: imageURL,
                Alt: *object.Key, // Use the object key as the alt text
            }
            images = append(images, image)
        }
	}
	log.Printf("Created %d video objects", len(images))

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Encode the image objects as JSON and write the response
	err = json.NewEncoder(w).Encode(images)
	if err != nil {
		log.Printf("Failed to encode video objects as JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
