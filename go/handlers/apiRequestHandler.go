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
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	Name      string   `json:"name"`
	Category  Category `json:"category"`
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
	PetID int `json:"option"`
}

func HandleApiGet(c echo.Context) error {
	log.Printf("Start HandleApiGet")

	var status PetstoreGetRequest
	if err := c.Bind(&status); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Log the Status
	log.Printf("Status: %s", status.Option)

	apiURL := fmt.Sprintf("%s/%s", config.CurrentConfig.PETSTOREURL, status.Option)
	log.Printf("API URL: %s", apiURL)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Generate curl command string
	curlCommand := utils.GenerateCurlCommand(req, nil)

	// Log the curl command
	log.Printf("Curl Command: %s", curlCommand)

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
	log.Printf("Response Body: %s", string(body))

	contentType := resp.Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	// Construct the response object with the curl command and the URL
	response := map[string]interface{}{
		"data":        nil,
		"url":         req.URL.String(),
		"curlCommand": curlCommand,
	}

	if strings.Contains(contentType, "application/json") {
		var jsonData interface{}
		err := json.Unmarshal(body, &jsonData)
		if err != nil {
			log.Printf("Received non-JSON content despite 'application/json' content-type: %v", err)
			response["data"] = string(body)
		} else {
			response["data"] = jsonData
		}
	} else if strings.Contains(contentType, "text/plain") || strings.Contains(contentType, "text/html") {
		response["data"] = string(body)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported content type"})
	}

	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// POST                                                                          //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiPost(c echo.Context) error {
	apiURL := config.CurrentConfig.PETSTOREURL
	log.Printf("API URL: %s", apiURL)

	// Read the request body
	var request PetstorePostRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	log.Printf("Decoded request: %+v", request)

	// Prepare the payload for the POST request
	payload, err := json.Marshal(request.Option)
	if err != nil {
		log.Printf("Error marshalling the payload: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	log.Printf("Marshalled payload: %s", string(payload))

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
	log.Printf("Curl command: %s", curlCommand)

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
	log.Printf("Response status: %s", resp.Status)

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
	log.Printf("Content-Type: %s", contentType)

	if strings.Contains(contentType, "application/json") {
		var pets PetstorePet
		if err := json.Unmarshal(responseBody, &pets); err != nil {
			log.Printf("Error unmarshalling response body: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		response["data"] = pets
	} else if strings.Contains(contentType, "text/plain") {
		response["data"] = string(responseBody)
	} else if strings.Contains(contentType, "text/html") {
		response["data"] = string(responseBody)
	} else {
		log.Printf("Unsupported content type")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported content type"})
	}

	return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// PUT                                                                           //
///////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////
// PUT                                                                           //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiPut(c echo.Context) error {
    apiURL := config.CurrentConfig.PETSTOREURL
    log.Printf("API URL: %s", apiURL)

    // Read the request body
    var request PetstorePutRequest
    if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
        log.Printf("Error decoding request: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    log.Printf("Decoded request: %+v", request)

    // Prepare the payload for the PUT request
    payload, err := json.Marshal(request.Option)
    if err != nil {
        log.Printf("Error marshalling the payload: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    log.Printf("Marshalled payload: %s", string(payload))

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
    log.Printf("Curl command: %s", curlCommand)

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
    log.Printf("Response status: %s", resp.Status)

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
    log.Printf("Content-Type: %s", contentType)

    if strings.Contains(contentType, "application/json") {
        var pets PetstorePet
        if err := json.Unmarshal(responseBody, &pets); err != nil {
            log.Printf("Error unmarshalling response body: %v", err)
            return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
        }
        response["data"] = pets
    } else if strings.Contains(contentType, "text/plain") {
        response["data"] = string(responseBody)
    } else if strings.Contains(contentType, "text/html") {
        response["data"] = string(responseBody)
    } else {
        log.Printf("Unsupported content type")
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported content type"})
    }

    return c.JSON(http.StatusOK, response)
}

///////////////////////////////////////////////////////////////////////////////////
// DELETE                                                                        //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiDelete(c echo.Context) error {
    // Read the request body into PetstoreDeleteRequest
    var request PetstoreDeleteRequest
    if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
        log.Printf("Error decoding request: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    log.Printf("Decoded request: %+v", request)

    // Construct the URL with the pet ID from the request
    apiURL := fmt.Sprintf("%s/%d", config.CurrentConfig.PETSTOREURL, request.PetID)
    log.Printf("API URL: %s", apiURL)

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
    log.Printf("Curl command: %s", curlCommand)

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
    log.Printf("Response status: %s", resp.Status)

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
    log.Printf("Content-Type: %s", contentType)

    if strings.Contains(contentType, "application/json") {
        var jsonResponse map[string]interface{}
        if err := json.Unmarshal(responseBody, &jsonResponse); err != nil {
            response["data"] = string(responseBody)
        } else {
            response["data"] = jsonResponse
        }
    } else if strings.Contains(contentType, "text/plain") || strings.Contains(contentType, "text/html") {
        response["data"] = string(responseBody)
    } else {
        log.Printf("Unsupported content type")
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported content type"})
    }

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

	// Print equivalent curl command
	// fmt.Printf("curl -X PUT %s -H \"Accept: application/json\" -H \"Content-Type: application/json\" -H \"User-Agent: %s\" -H \"X-Forwarded-For: %s\" -d '%s'\n", petStoreURL, userAgent, xForwardedFor, jsonData)

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

///////////////////////////////////////////////////////////////////////////////////
// handleAPITrafficGenerator                                                     //
///////////////////////////////////////////////////////////////////////////////////

func handleAPITrafficGenerator(c echo.Context) error {
	requestCount := 1800
	petNames := []string{"FortiPuma", "FortiFish", "FortiSpider", "FortiTiger", "FortiLion", "FortiShark", "FortiSnake", "FortiMonkey", "FortiFox", "FortiRam", "FortiEagle", "FortiBee", "FortiCat", "FortiDog", "FortiAnt", "FortiWasp", "FortiPanter", "FortiGator", "FortiOwl", "FortiWildcats"}
	petTypes := []string{"Puma", "Fish", "Spider", "Tiger", "Lion", "Shark", "Snake", "Monkey", "Fox", "Ram", "Eagle", "Bee", "Cat", "Dog", "Ant", "Wasp", "Panter", "Gator", "Owl", "Wildcats"}
	petTags := []string{"Friendly", "Playful", "Loyal", "Energetic", "Calm", "Trained", "Rescue", "Fluffy", "Affectionate", "Smart", "Active", "Gentle", "Senior", "Cute", "Good"}
	statuses := []string{"available", "pending", "sold"}

	for i := 0; i < requestCount; i++ {
		randomName := generateRandomValue(petNames)
		randomPet := generateRandomValue(petTypes)
		randomTag := generateRandomValue(petTags)
		randomStatus := generateRandomValue(statuses)
		randomStatusNew := generateRandomValue(statuses)
		randomIP := utils.GenerateRandomPublicIP()
		randomID := rand.Intn(1001)
		randomPhotoUrl1 := randomString(2, 10)
		randomPhotoUrl2 := randomString(2, 10)
		userAgent := config.CurrentConfig.USERAGENT
		petStoreURL := config.CurrentConfig.PETSTOREURL

		petNew := PetstorePet{
			ID: randomID,
			Category: Category{
				ID:   randomID,
				Name: randomPet,
			},
			Name:      randomName,
			PhotoUrls: []string{randomPhotoUrl1 + ".png", randomPhotoUrl2 + ".png"},
			Tags: []Tag{
				{
					ID:   randomID,
					Name: randomTag,
				},
			},
			Status: randomStatus,
		}

		petModified := PetstorePet{
			ID: randomID,
			Category: Category{
				ID:   randomID,
				Name: randomPet,
			},
			Name:      randomName,
			PhotoUrls: []string{randomPhotoUrl1 + ".png", randomPhotoUrl2 + ".png"},
			Tags: []Tag{
				{
					ID:   randomID,
					Name: randomTag,
				},
			},
			Status: randomStatusNew,
		}

		// Send POST request
		err := sendApiPostRequest(petStoreURL, userAgent, petNew, randomIP)
		if err != nil {
			log.Fatalf("Error sending POST request: %v", err)
		}

		// Send PUT request
		err = sendApiPutRequest(petStoreURL, userAgent, petModified, randomIP)
		if err != nil {
			log.Fatalf("Error sending PUT request: %v", err)
		}

		// Send GET request
		err = sendApiGetRequest(petStoreURL, randomStatus, userAgent, randomIP)
		if err != nil {
			log.Fatalf("Error sending PUT request: %v", err)
		}

		// Send DELETE request
		err = sendApiDeleteRequest(petStoreURL, randomID, userAgent, randomIP)
		if err != nil {
			log.Fatalf("Error sending PUT request: %v", err)
		}
	}

	// Return the completion message
	message := fmt.Sprintf("API traffic generation is complete. We have sent %d random requests of types POST, PUT, GET, and DELETE.", requestCount)
	return c.String(http.StatusOK, message)
}

///////////////////////////////////////////////////////////////////////////////////
// GENERATE RANDOM API TRAFFIC                                                   //
///////////////////////////////////////////////////////////////////////////////////

func generateRandomValue(values []string) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return values[random.Intn(len(values))]
}

func randomString(minLength, maxLength int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@_-[]{}$!")
	var startingLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	length := random.Intn(maxLength-minLength+1) + minLength

	var sb strings.Builder
	sb.WriteRune(startingLetters[random.Intn(len(startingLetters))]) // Start with A-Z

	for i := 1; i < length; i++ {
		sb.WriteRune(letters[random.Intn(len(letters))])
	}

	return sb.String()
}
