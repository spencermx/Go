package main

import (
	"log"
	"net/http"
    "goserver/handlers"
//    "github.com/gorilla/handlers"
// "html/template"
)

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

    http.HandleFunc("/getPeople", handlers.GetPeople)

    http.HandleFunc("/uploadImage", handlers.UploadImage)

    http.HandleFunc("/uploadVideo", handlers.UploadVideo)

    http.HandleFunc("/getImages", handlers.GetImages)

    http.HandleFunc("/getVideos", handlers.GetVideos)

//    http.Handle("/getImages", handlers.CORS(
//        handlers.AllowedOrigins([]string{"http://localhost:5173"}),
//        handlers.AllowedMethods([]string{"GET"}),
//    )(http.HandlerFunc(getImages)))     

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

// func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
//     http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
// }
