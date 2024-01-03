package utils

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateAPIToken() string {
	tokenData := fmt.Sprintf(`{"username":"%s","password":"%s","vdom":"%s"}`, config.CurrentConfig.USERNAMEAPI, config.CurrentConfig.PASSWORDAPI, config.CurrentConfig.VDOMAPI)
	return base64.StdEncoding.EncodeToString([]byte(tokenData))
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SendRequest(method, url, token string, data interface{}) ([]byte, error) {
	// fmt.Println("Sending HTTP Request...")
	var req *http.Request
	var err error

	reqData := config.Request{
		Data: data,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("Error marshaling JSON data: %v\n", err)
		return nil, err
	}

	// Convert jsonData to string for comparison
	jsonDataStr := string(jsonData)

	if jsonDataStr != "" && jsonDataStr != `{"data":null}` {
		// Create a new request with JSON data
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating HTTP request: %v\n", err)
			return nil, err
		}
	} else {
		// Create a new request without data
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Printf("Error creating HTTP request: %v\n", err)
			return nil, err
		}
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	// Create a custom HTTP client with SSL/TLS certificate verification disabled
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// fmt.Println("Sending HTTP Request...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	// fmt.Println("Waiting for response...")
	time.Sleep(time.Duration(1000) * time.Millisecond)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	// fmt.Println("Received HTTP Response")
	return body, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckOperationStatus(result []byte) bool {
	var res map[string]interface{}
	json.Unmarshal(result, &res)

	// Check if the result map is empty
	if len(res) == 0 {
		return false
	}

	// Check if the "data" field is null
	if data, ok := res["data"]; ok && data == nil {
		return false
	}

	if _, ok := res["results"].(map[string]interface{})["errcode"]; ok {
		// The result contains an error code, so the operation failed
		return false
	}
	// The operation succeeded
	return true
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
