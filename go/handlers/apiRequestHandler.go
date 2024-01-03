package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////////////////////////////////////
// STRUCTURE                                                                     //
///////////////////////////////////////////////////////////////////////////////////

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PetstorePet struct {
	ID        int      `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}

type PetstorePetArray []PetstorePet

type PetstoreGetRequest struct {
	Option string `json:"option"`
}

type PetstorePostRequest struct {
	Option PetstorePet `json:"option"`
}

type PetstorePutRequest struct {
	Option PetstorePet `json:"option"`
}

type PetstoreDeleteRequest struct {
	PetID string `json:"option"`
}

///////////////////////////////////////////////////////////////////////////////////
// MAIN                                                                          //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiGet(c echo.Context) error {
	// fmt.Printf("Start HandleApiGet\n")

	var status PetstoreGetRequest
	if err := c.Bind(&status); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Log the Status
	// fmt.Printf("Status: %s\n", status.Option)

	apiURL := fmt.Sprintf("%s/%s", config.CurrentConfig.PETSTOREURL, status.Option)
	// fmt.Printf("API URL: %s\n", apiURL)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Generate curl command string
	curlCommand := utils.GenerateCurlCommand(req, nil)

	// Log the curl command
	// fmt.Printf("Curl Command: %s\n", curlCommand)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Log the response body for debugging
	// fmt.Printf("Response Body: %s", string(body))

	contentType := resp.Header.Get("Content-Type")
	// fmt.Printf("Content-Type: %s\n", contentType)

	// Construct the response object with the curl command and the URL
	response := map[string]interface{}{
		"data":        nil,
		"url":         req.URL.String(),
		"curlCommand": curlCommand,
	}

	if strings.HasPrefix(contentType, "application/json") {
		var jsonData PetstorePetArray
		if err := json.Unmarshal(body, &jsonData); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			response["data"] = string(body) // return raw body as fallback
		} else {
			response["data"] = jsonData
		}
	} else {
		// If not JSON, return the raw body
		response["data"] = string(body)
	}

	// fmt.Printf("Response data: %v\n", response["data"])
	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// POST                                                                          //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiPost(c echo.Context) error {
	apiURL := config.CurrentConfig.PETSTOREURL
	// fmt.Printf("API URL: %s", apiURL)

	// Read the request body
	var request PetstorePostRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// fmt.Printf("Decoded request: %+v", request)

	// Prepare the payload for the POST request
	payload, err := json.Marshal(request.Option)
	if err != nil {
		log.Printf("Error marshalling the payload: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// fmt.Printf("Marshalled payload: %s", string(payload))

	// Create a new POST request using the received body
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set headers for the request
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Generate curl command string
	curlCommand := utils.GenerateCurlCommand(req, payload)
	// fmt.Printf("Curl command: %s", curlCommand)

	// Create a custom HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()
	// fmt.Printf("Response status: %s", resp.Status)

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// log.Printf("Response body: %s", string(responseBody))

	// Construct the response object with the curl command
	response := map[string]interface{}{
		"data":        nil,
		"url":         req.URL.String(),
		"curlCommand": curlCommand,
	}

	contentType := resp.Header.Get("Content-Type")
	// fmt.Printf("Content-Type: %s", contentType)

	if strings.HasPrefix(contentType, "application/json") {
		var jsonData PetstorePet
		if err := json.Unmarshal(responseBody, &jsonData); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			response["data"] = string(responseBody) // return raw body as fallback
		} else {
			response["data"] = jsonData
		}
	} else {
		// If not JSON, return the raw body
		response["data"] = string(responseBody)
	}

	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// PUT                                                                           //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiPut(c echo.Context) error {
	apiURL := config.CurrentConfig.PETSTOREURL
	// fmt.Printf("API URL: %s", apiURL)

	// Read the request body
	var request PetstorePutRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// fmt.Printf("Decoded request: %+v", request)

	// Prepare the payload for the PUT request
	payload, err := json.Marshal(request.Option)
	if err != nil {
		log.Printf("Error marshalling the payload: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// fmt.Printf("Marshalled payload: %s", string(payload))

	// Create a new PUT request using the received body
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set headers for the request
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Generate curl command string
	curlCommand := utils.GenerateCurlCommand(req, payload)
	// fmt.Printf("Curl command: %s", curlCommand)

	// Create a custom HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()
	// fmt.Printf("Response status: %s", resp.Status)

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Construct the response object with the curl command
	response := map[string]interface{}{
		"data":        nil,
		"url":         req.URL.String(),
		"curlCommand": curlCommand,
	}

	contentType := resp.Header.Get("Content-Type")
	// fmt.Printf("Content-Type: %s", contentType)

	if strings.HasPrefix(contentType, "application/json") {
		var jsonData PetstorePet
		if err := json.Unmarshal(responseBody, &jsonData); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			response["data"] = string(responseBody) // return raw body as fallback
		} else {
			response["data"] = jsonData
		}
	} else {
		// If not JSON, return the raw body
		response["data"] = string(responseBody)
	}

	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// DELETE                                                                        //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiDelete(c echo.Context) error {
	// Read the request body into PetstoreDeleteRequest
	var request PetstoreDeleteRequest
	// fmt.Printf("Raw request: %+v\n", request)

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// fmt.Printf("Decoded request: %+v\n", request)

	// Construct the URL with the pet ID from the request
	apiURL := fmt.Sprintf("%s/%s", config.CurrentConfig.PETSTOREURL, request.PetID)
	// fmt.Printf("API URL: %s\n", apiURL)

	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set headers for the request
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Generate curl command string
	curlCommand := utils.GenerateCurlCommand(req, nil)
	// fmt.Printf("Curl command: %s\n", curlCommand)

	// Create a custom HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()
	// fmt.Printf("Response status: %s\n", resp.Status)

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// fmt.Printf("Response Body: %s\n", responseBody)

	// Construct the response object with the curl command
	response := map[string]interface{}{
		"data":        nil,
		"url":         req.URL.String(),
		"curlCommand": curlCommand,
	}

	response["data"] = string(responseBody)

	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// sendApiGetRequest                                                                //
///////////////////////////////////////////////////////////////////////////////////

func sendApiGetRequest(petStoreURL, randomStatus, userAgent, xForwardedFor string) error {
	// Construct the URL with query parameters
	fullURL := fmt.Sprintf("%s/findByStatus?status=%s", petStoreURL, url.QueryEscape(randomStatus))

	// Create the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Forwarded-For", xForwardedFor)

	// Send the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and log the response body (optional)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Response Body: %s\n", body)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// sendApiPostRequest                                                               //
///////////////////////////////////////////////////////////////////////////////////

func sendApiPostRequest(petStoreURL string, userAgent string, pet PetstorePet, xForwardedFor string) error {
	jsonData, err := json.Marshal(pet)
	if err != nil {
		log.Printf("Error marshalling pet data: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", petStoreURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v\n", err)
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Forwarded-For", xForwardedFor)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Read and log the response body (optional)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Response Body: %s\n", body)

	// Print various response fields
	// fmt.Printf("Status: %s\n", resp.Status)
	// fmt.Printf("StatusCode: %d\n", resp.StatusCode)
	// fmt.Printf("Header: %v\n", resp.Header)
	// fmt.Printf("Body: %s\n", body)
	// fmt.Printf("ContentLength: %d\n", resp.ContentLength)
	// fmt.Printf("TransferEncoding: %v\n", resp.TransferEncoding)
	// fmt.Printf("Close: %v\n", resp.Close)
	// fmt.Printf("Uncompressed: %v\n", resp.Uncompressed)

	// Handle the response as needed
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// sendApiPutRequest                                                                //
///////////////////////////////////////////////////////////////////////////////////

func sendApiPutRequest(petStoreURL string, userAgent string, pet PetstorePet, xForwardedFor string) error {
	jsonData, err := json.Marshal(pet)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", petStoreURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Forwarded-For", xForwardedFor)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and log the response body (optional)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Response Body: %s\n", body)

	// Handle the response as needed
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// sendApiDeleteRequest                                                             //
///////////////////////////////////////////////////////////////////////////////////

func sendApiDeleteRequest(petStoreURL string, randomID int, userAgent string, xForwardedFor string) error {
	// Construct the URL
	fullURL := fmt.Sprintf("%s/%d", petStoreURL, randomID)

	// Create the request
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Forwarded-For", xForwardedFor)

	// Send the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and log the response body (optional)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Response Body: %s\n", body)

	return nil
}
