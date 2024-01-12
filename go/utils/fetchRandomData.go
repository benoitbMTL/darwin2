package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var fakeNameURL = "https://api.namefake.com"

// Struct to map the JSON data returned from api.namefake.com
type FakeData struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Birthday   string `json:"birth_data"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	EmailU     string `json:"email_u"`
	EmailD     string `json:"email_d"`
	PhoneH     string `json:"phone_h"`
	Useragent  string `json:"useragent"`
	Ipv4       string `json:"ipv4"`
	MiddleName string `json:"maiden_name"`
	CreditCard string `json:"plasticcard"`
}

// Credentials for Juice Shop
type Credentials struct {
	Email     string
	Password  string
	SecAnswer string
}

// Address for Juice Shop
type Address struct {
	Country string
	Name    string
	Mobile  string
	ZipCode string
	Address string
	City    string
}

// Payment for Juice Shop
type Payment struct {
	Name   string
	CardNr string
	Month  string
	Year   string
}

func FetchRandomData() (*FakeData, error) {

	resp, err := http.Get(fakeNameURL)
	if err != nil {
		log.Printf("Error to Get Random Name: %v\n", err) // Log the error
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data FakeData
	// Use bodyBytes for JSON decoding
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Printf("Error to Decode Random Name: %v\n", err) // Log the error
		return nil, err
	}

	return &data, nil
}

func BuildEmail(data FakeData) string {
	email := fmt.Sprintf("%s@%s", data.EmailU, data.EmailD)
	return email
}

func RandomCredentials(fakeData FakeData) Credentials {
	return Credentials{
		Email:     BuildEmail(fakeData),
		Password:  fakeData.Password,
		SecAnswer: fakeData.MiddleName,
	}
}

func RandomAddress(fakeData FakeData) Address {
	countries := []string{"Canada", "United States", "Mexico", "Brazil"}
	cities := []string{"Montreal", "New York", "Mexico City", "Sao Paulo"} // Example cities

	// Generate a random number and convert it to a string
	min := 1000000
	max := 9999999999
	randomNumberStr := fmt.Sprintf("%d", rand.Intn(max-min+1)+min)


	return Address{
		Country: countries[rand.Intn(len(countries))],
		Name:    fakeData.Name,
		Mobile:  randomNumberStr,
		ZipCode: fmt.Sprintf("%05d", rand.Intn(100000)),   // Correct for 5-digit number
		Address: fakeData.Address,
		City:    cities[rand.Intn(len(cities))],
	}
}

func RandomPayment(fakeData FakeData) Payment {
	cardNumber := ensureSixteenDigits(fakeData.CreditCard)

	return Payment{
		Name:   fakeData.Name,
		CardNr: cardNumber,
		Month:  fmt.Sprintf("%02d", rand.Intn(12)+1), // Generates a random month from 01 to 12
		Year:   "2080",
	}
}

func ensureSixteenDigits(number string) string {
	length := len(number)
	switch {
	case length > 16:
		return number[:16] // Trim the number to 16 digits
	case length < 16:
		for i := length; i < 16; i++ {
			number += fmt.Sprintf("%d", rand.Intn(10)) // Add random digits
		}
		return number
	default:
		return number // Already 16 digits
	}
}
