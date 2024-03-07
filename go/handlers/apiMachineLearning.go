package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////////////////////////////////////
// MAIN FUNCTION                                                                 //
///////////////////////////////////////////////////////////////////////////////////

func HandleApiMachineLearning(c echo.Context) error {
	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form data")
	}

	requestCountStr := c.Request().FormValue("count")
	if requestCountStr == "" {
		return c.String(http.StatusBadRequest, "Count parameter is required")
	}
	// fmt.Printf("Debug: requestCountStr = %s\n", requestCountStr)

	requestCount, err := strconv.Atoi(requestCountStr)
	if err != nil {
		// fmt.Printf("Debug: Error converting requestCountStr to int: %v\n", err)
		return c.String(http.StatusBadRequest, "Invalid count parameter")
	}

	petNames := []string{"FortiPuma", "FortiFish", "FortiSpider", "FortiTiger", "FortiLion", "FortiShark", "FortiSnake", "FortiMonkey", "FortiFox", "FortiRam", "FortiEagle", "FortiBee", "FortiCat", "FortiDog", "FortiAnt", "FortiWasp", "FortiPanter", "FortiGator", "FortiOwl", "FortiWildcats"}
	petTypes := []string{"Puma", "Fish", "Spider", "Tiger", "Lion", "Shark", "Snake", "Monkey", "Fox", "Ram", "Eagle", "Bee", "Cat", "Dog", "Ant", "Wasp", "Panter", "Gator", "Owl", "Wildcats"}
	petTags := []string{"Friendly", "Playful", "Loyal", "Energetic", "Calm", "Trained", "Rescue", "Fluffy", "Affectionate", "Smart", "Active", "Gentle", "Senior", "Cute", "Good"}
	statuses := []string{"available", "pending", "sold"}

	var lastPetCreated PetstorePet
	var petNew PetstorePet // Declare petNew outside the loop

	for i := 0; i < requestCount; i++ {
		randomName := generateRandomValue(petNames)
		randomPet := generateRandomValue(petTypes)
		randomTag := generateRandomValue(petTags)
		randomStatus := generateRandomValue(statuses)
		randomStatusNew := generateRandomValue(statuses)
		randomIP := utils.GenerateRandomPublicIP()
		randomID := rand.Intn(1001)
		randomPhotoUrl1 := randomStringPhotoUrl(2, 10)
		randomPhotoUrl2 := randomStringPhotoUrl(2, 10)
		userAgent := config.CurrentConfig.USERAGENT
		petStoreURL := config.CurrentConfig.PETSTOREURL

		petNew = PetstorePet{
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
		if err := sendApiPostRequestRandom(petStoreURL, userAgent, petNew, randomIP); err != nil {
			log.Printf("Error sending POST request: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error sending POST request: %v", err))
		}

		// Send PUT request
		if err := sendApiPutRequestRandom(petStoreURL, userAgent, petModified, randomIP); err != nil {
			log.Printf("Error sending PUT request: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error sending PUT request: %v", err))
		}

		// Send GET request
		if err := sendApiGetRequestRandom(petStoreURL, randomStatus, userAgent, randomIP); err != nil {
			log.Printf("Error sending GET request: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error sending GET request: %v", err))
		}

		// Send DELETE request
		if err := sendApiDeleteRequestRandom(petStoreURL, randomID, userAgent, randomIP); err != nil {
			log.Printf("Error sending DELETE request: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error sending DELETE request: %v", err))
		}

    // Last iteration of the loop
    if i == requestCount-1 {
        lastPetCreated = petNew
    }

	}

// Format the last pet created as JSON
lastPetJson, err := json.MarshalIndent(lastPetCreated, "", "  ")
if err != nil {
    // handle error
}

// Return the completion message
var message string
if requestCount == 1 {
    message = fmt.Sprintf("The API traffic generation has successfully completed, with a total of %d request each for POST, PUT, GET, and DELETE types.\n\nLast Pet Created (POST):\n\n%s", requestCount, lastPetJson)
} else {
    message = fmt.Sprintf("The API traffic generation has successfully completed, with a total of %d requests each for POST, PUT, GET, and DELETE types.\n\nLast Pet Created (POST):\n\n%s", requestCount, lastPetJson)
}

return c.String(http.StatusOK, message)

}

///////////////////////////////////////////////////////////////////////////////////
// randomStringPhotoUrl                                                          //
///////////////////////////////////////////////////////////////////////////////////

func randomStringPhotoUrl(minLength, maxLength int) string {
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

///////////////////////////////////////////////////////////////////////////////////
// sendApiPostRequestRandom                                                      //
///////////////////////////////////////////////////////////////////////////////////

func sendApiPostRequestRandom(petStoreURL string, userAgent string, pet PetstorePet, xForwardedFor string) error {
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

	// Handle the response as needed
	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// sendApiPutRequestRandom                                                       //
///////////////////////////////////////////////////////////////////////////////////

func sendApiPutRequestRandom(petStoreURL string, userAgent string, pet PetstorePet, xForwardedFor string) error {
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
// sendApiGetRequestRandom                                                       //
///////////////////////////////////////////////////////////////////////////////////

func sendApiGetRequestRandom(petStoreURL, randomStatus, userAgent, xForwardedFor string) error {
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
// sendApiDeleteRequestRandom                                                    //
///////////////////////////////////////////////////////////////////////////////////

func sendApiDeleteRequestRandom(petStoreURL string, randomID int, userAgent string, xForwardedFor string) error {
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
// generateRandomValue                                                           //
///////////////////////////////////////////////////////////////////////////////////

func generateRandomValue(values []string) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return values[random.Intn(len(values))]
}