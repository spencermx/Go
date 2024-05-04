package main

import (
//    "fmt"
//    "html/template"
    "log"
    "net/http"
//    "os"
//
//    "github.com/aws/aws-sdk-go/aws"
//    "github.com/aws/aws-sdk-go/aws/session"
//    "github.com/aws/aws-sdk-go/service/s3"
//    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)
func main() {
    fs := http.FileServer(http.Dir("../client/build/"))
    http.Handle("/", fs)

    log.Println("Server running on http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
// func main() {
//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         http.ServeFile(w, r, "static/index.html")
//     })
// 
//     http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
// 
//     err := http.ListenAndServe(":8080", nil)
//     if err != nil {
//         log.Fatal(err)
//     }
// }

// func main() {
// 
//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         http.ServeFile(w, r, "hello.html")
//     })
// 
//     // Start the server
//     log.Println("Server is running on http://localhost:8080")
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

    // Create a log file
//    logFile, err := os.Create("server.log")
//    if err != nil {
//        log.Fatalf("Failed to create log file: %v", err)
//    }
//    defer logFile.Close()
//
//    // Set the log output to the file
//    log.SetOutput(logFile)
//
//    log.Println("Server is starting...")
//    
//    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//        http.ServeFile(w, r, "index.html")
//    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//        http.ServeFile(w, r, "index.html")
//    })
//    // Serve the HTML file
//   // http.HandleFunc("/", handleHome1)
//   // log.Println("Registered route: /")
//
//    // Handle the file upload
//    http.HandleFunc("/upload", handleUpload)
//    log.Println("Registered route: /upload")
//
//    // Start the server
//    log.Println("Server is running on http://localhost:8080")
//    log.Fatal(http.ListenAndServe(":8080", nil))


// func handleUpload(w http.ResponseWriter, r *http.Request) {
//     // Parse the multipart form
//     err := r.ParseMultipartForm(10 << 20) // Max size of 10MB
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
// 
//     // Get the uploaded file
//     file, header, err := r.FormFile("file")
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
//     
//     defer file.Close()
// 
//     // Create a new AWS session
//     sess, err := session.NewSession(&aws.Config{
//     	Region: aws.String("us-east-2"),
//     })
// 
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// 
//     // Create an S3 uploader
//     uploader := s3manager.NewUploader(sess)
// 
//     // Upload the file to S3
//     _, err = uploader.Upload(&s3manager.UploadInput{
//         Bucket: aws.String("goserverbucket"),
//         Key:    aws.String(header.Filename),
//         Body:   file,
//     })
// 
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     // Redirect to the home page
//     http.Redirect(w, r, "/", http.StatusSeeOther)
//     log.Println("Redirected to home page")
// }
// 
// func handleHome(w http.ResponseWriter, r *http.Request) {
//     log.Println("Handling request for /")
// 
//     // Create a new AWS session
//     sess, err := session.NewSession(&aws.Config{
//         Region: aws.String("us-east-2"),
//     })
//     if err != nil {
//         log.Printf("Failed to create AWS session: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     log.Println("Created AWS session")
// 
//     // Create an S3 client
//     svc := s3.New(sess)
//     log.Println("Created S3 client")
// 
//     // List objects in the S3 bucket
//     result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
//         Bucket: aws.String("goserverbucket"),
//     })
//     if err != nil {
//         log.Printf("Failed to list objects in S3 bucket: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     log.Printf("Listed %d objects in S3 bucket", len(result.Contents))
// 
//     // Create a slice to store the image URLs
//     var imageURLs []string
// 
//     // Iterate over the objects and create the image URLs
//     for _, obj := range result.Contents {
//         imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
//         imageURLs = append(imageURLs, imageURL)
//     }
//     log.Printf("Created %d image URLs", len(imageURLs))
// 
//     // Create a template data structure
//     data := struct {
//         ImageURLs []string
//     }{
//         ImageURLs: imageURLs,
//     }
// 
//     // Render the HTML template with the template data
//     tmpl := template.Must(template.ParseFiles("index.html"))
//     err = tmpl.Execute(w, data)
//     if err != nil {
//         log.Printf("Failed to render template: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     log.Println("Rendered template")
// }

// package main
// 
// 
// import (
//     "fmt"
//     "html/template"
//     "net/http"
//     "log"
//     "github.com/aws/aws-sdk-go/aws"
//     "github.com/aws/aws-sdk-go/aws/session"
//     "github.com/aws/aws-sdk-go/service/s3"
//     "github.com/aws/aws-sdk-go/service/s3/s3manager"
// ) 
// 
// func main() {
//     // Serve the HTML file
//     http.HandleFunc("/", handleHome)
// 
//     // Handle the file upload
//  //   http.HandleFunc("/upload", handleUpload)
// 
//     // Serve the HTML file
//     // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//     //     http.ServeFile(w, r, "upload.html")
//     // })
// 
// //    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// //        fmt.Fprintf(w, "Hello from Fort Collins")
// //    })
// 
//     // Start the server
//     fmt.Println("Server is running on http://localhost:8080")
//     http.ListenAndServe(":8080", nil)
// 
//     fmt.Println("xxServer is running on http://localhost:8080")
// }
// func handleHome(w http.ResponseWriter, r *http.Request) {
//     // Create a new AWS session
//     sess, err := session.NewSession(&aws.Config{
//         Region: aws.String("us-east-2"),
//     })
//     if err != nil {
//         log.Printf("Failed to create AWS session: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// 
//     // Create an S3 client
//     svc := s3.New(sess)
// 
//     // List objects in the S3 bucket
//     result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
//         Bucket: aws.String("goserverbucket"),
//     })
//     if err != nil {
//         log.Printf("Failed to list objects in S3 bucket: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// 
//     // Create a slice to store the image URLs
//     var imageURLs []string
// 
//     // Iterate over the objects and create the image URLs
//     for _, obj := range result.Contents {
//         imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
//         imageURLs = append(imageURLs, imageURL)
//     }
// 
//     // Render the HTML template with the image URLs
//     tmpl := template.Must(template.ParseFiles("index.html"))
//     err = tmpl.Execute(w, imageURLs)
//     if err != nil {
//         log.Printf("Failed to render template: %v", err)
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// }
// 
// func handleUpload(w http.ResponseWriter, r *http.Request) {
//     // Parse the multipart form
//     err := r.ParseMultipartForm(10 << 20) // Max size of 10MB
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
// 
//     // Get the uploaded file
//     file, header, err := r.FormFile("file")
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
//     
//     defer file.Close()
// 
//     // Create a new AWS session
//     sess, err := session.NewSession(&aws.Config{
//     	Region: aws.String("us-east-2"),
//     })
// 
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// 
//     // Create an S3 uploader
//     uploader := s3manager.NewUploader(sess)
// 
//     // Upload the file to S3
//     _, err = uploader.Upload(&s3manager.UploadInput{
//         Bucket: aws.String("goserverbucket"),
//         Key:    aws.String(header.Filename),
//         Body:   file,
//     })
// 
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// 
//     fmt.Fprintf(w, "File uploaded successfully")
// }
// 
// // func main() {
// //     // Create a new session using the default AWS configuration
// //     sess, err := session.NewSession(&aws.Config{
// //         Region: aws.String("us-west-2"), // Replace with your desired AWS region
// //     })
// //     if err != nil {
// //         fmt.Println("Error creating session:", err)
// //         return
// //     }
// // 
// //     // Create an S3 client
// //     s3Client := s3.New(sess)
// // 
// //     // Specify the S3 bucket and object key
// //     bucket := "your-bucket-name"
// //     key := "path/to/your/object"
// // 
// //     // Get the object from the S3 bucket
// //     result, err := s3Client.GetObject(&s3.GetObjectInput{
// //         Bucket: aws.String(bucket),
// //         Key:    aws.String(key),
// //     })
// //     if err != nil {
// //         fmt.Println("Error getting object:", err)
// //         return
// //     }
// //     defer result.Body.Close()
// // 
// //     // Process the object data
// //     // ...
// // 
// //     fmt.Println("Successfully accessed S3 object")
// // }