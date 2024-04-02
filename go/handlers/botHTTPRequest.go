package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"golang.org/x/net/http2"

	"github.com/labstack/echo/v4"
)

// HTTPRequestData represents all the parameters for the HTTP request from the frontend.
type HTTPRequestData struct {
	Method          string `json:"method"`
	URL             string `json:"url"`
	LoopCount       int    `json:"loopCount"`
	FollowRedirects bool   `json:"followRedirects"`
	DataContent     string `json:"dataContent"`
	UserAgent       string `json:"userAgent"`
	ContentType     string `json:"contentType"`
	Cookie          string `json:"cookie"`
	XForwardedFor   string `json:"xForwardedFor"`
	Referer         string `json:"referer"`
}

type HTTPResponseData struct {
	RequestHeaders  string `json:"requestHeaders"`
	ResponseHeaders string `json:"responseHeaders"`
	ResponseBody    string `json:"responseBody"`
}

func HandleHTTPRequest(c echo.Context) error {
	var requestData HTTPRequestData

	// Parse the JSON body from the incoming request.
	if err := c.Bind(&requestData); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	// Convert requestData to JSON format with indentation
	requestDataJson, err := json.MarshalIndent(requestData, "", "  ")
	if err != nil {
		// Handle error
		fmt.Printf("Error marshalling requestData to JSON: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error processing request")
	}

	// Debug: Print the received request data.
	fmt.Printf("Received request:\n%s\n", requestDataJson)

	// Initialize the transport to always enable HTTP/2 and skip TLS verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Explicitly enable HTTP/2
	http2.ConfigureTransport(transport)

	// Create the HTTP client with configured transport and custom CheckRedirect
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Control redirect behavior based on user's choice
			if !requestData.FollowRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Prepare the HTTP request based on the received data.
	req, err := http.NewRequest(requestData.Method, requestData.URL, bytes.NewBuffer([]byte(requestData.DataContent)))
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Failed to create request: %v", err))
	}

	// Set headers based on requestData
	if requestData.UserAgent != "" {
		req.Header.Set("User-Agent", requestData.UserAgent)
	}
	if requestData.ContentType != "" {
		req.Header.Set("Content-Type", requestData.ContentType)
	}
	if requestData.Cookie != "" {
		req.Header.Set("Cookie", requestData.Cookie)
	}
	if requestData.XForwardedFor != "" {
		req.Header.Set("X-Forwarded-For", requestData.XForwardedFor)
	}
	if requestData.Referer != "" {
		req.Header.Set("Referer", requestData.Referer)
	}

	// Dump request headers (before the request is sent)
	requestDump, err := httputil.DumpRequestOut(req, true) // Use DumpRequestOut to include the body if present
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Failed to dump request: %v", err))
	}

	// Before starting the request loop, check for loop count validity
	if requestData.LoopCount < 1 {
		requestData.LoopCount = 1 // Ensure there's at least one iteration
	}

	var lastResponseData HTTPResponseData // To store data from the last iteration

	for i := 0; i < requestData.LoopCount; i++ {

		fmt.Printf("Loop Request: %d\n", i)

		// Perform the HTTP request
		resp, err := client.Do(req)
		if err != nil {
			return c.String(http.StatusOK, fmt.Sprintf("Request failed: %v", err))
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.String(http.StatusOK, fmt.Sprintf("Failed to read response body: %v", err))
		}

		// Dump response headers
		responseDump, err := httputil.DumpResponse(resp, false) // false because we've already read the body
		if err != nil {
			return c.String(http.StatusOK, fmt.Sprintf("Failed to dump response: %v", err))
		}

		// Respond back to the client with request and response headers, and body
		lastResponseData = HTTPResponseData{
			RequestHeaders:  string(requestDump),
			ResponseHeaders: string(responseDump),
			ResponseBody:    string(body),
		}
	}

	//	fmt.Printf("Request Headers:\n%s\n", lastResponseData.RequestHeaders)
	//	fmt.Printf("Response Headers:\n%s\n", lastResponseData.ResponseHeaders)
	//	fmt.Printf("Response Body:\n%s\n", lastResponseData.ResponseBody)

	return c.JSON(http.StatusOK, lastResponseData)
}
