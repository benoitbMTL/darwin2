package main

import (
    "fmt"
    "net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
    // Echo back the request method and URL
    fmt.Fprintf(w, "Request Method: %s\nRequested URL: %s", r.Method, r.URL.Path)
}

func main() {
    // Route to handle requests
    http.HandleFunc("/", echoHandler)

    // Start the server on port 8080
    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
