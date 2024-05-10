package main

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
)

// "html/template"
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

func main() {
    fs := http.FileServer(http.Dir("../client/build/"))

    //    http.Handle("/", fs)

    // Create a route handler that serves the static files and fallbacks to index.html
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Check if the requested file exists
        _, err := http.Dir("../client/build").Open(r.URL.Path)
        if err != nil {
            // If the file doesn't exist, serve index.html
            http.ServeFile(w, r, "../client/build/index.html")
            return
        }

        // If the file exists, serve it
        fs.ServeHTTP(w, r)
    })

    http.HandleFunc("/getPeople", getPeople)

    http.HandleFunc("/upload", handleUpload)

    http.HandleFunc("/getImages", getImages)
      
    log.Println("Server running on http://localhost:8080")

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }

//    // Get the environment variable for the server environment
//    env := os.Getenv("SERVER_ENV")
//
//    // Set the default port for local development
//    port := "8080"
//
//    // Check if the server is running in the deployed environment
//    if env == "production" {
//        // Use the specified port for the deployed server
//        port = "8443"
//
//        // Set up HTTPS server with SSL/TLS certificates
//        go func() {
//            httpPort := "8080"
//            log.Printf("HTTP server running on :%s\n", httpPort)
//            http.ListenAndServe(":"+httpPort, http.HandlerFunc(redirectToHTTPS))
//        }()
//
//        log.Printf("HTTPS server running on :%s\n", port)
//        err := http.ListenAndServeTLS(":"+port, "/etc/apache2/ssl/cert.pem", "/etc/apache2/ssl/key.pem", nil)
//        if err != nil {
//            log.Fatalf("HTTPS server ListenAndServeTLS: %v", err)
//        }
//    } else {
//        // Local development environment
//        log.Printf("Server running on http://localhost:%s\n", port)
//        http.ListenAndServe(":"+port, nil)
//    }
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

func getPeople(w http.ResponseWriter, r *http.Request) {
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

func handleUpload(w http.ResponseWriter, r *http.Request) {
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

	defer file.Close()

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
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

func getImages(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for /")

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
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
		imageURL := fmt.Sprintf("https://d271tjczb1hjof.cloudfront.net/%s", *obj.Key)
		image := Image{
			URL: imageURL,
			Alt: *obj.Key, // Use the object key as the alt text
		}
		images = append(images, image)
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
