package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// FakeData contains the locally generated identity data used by traffic simulations.
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
	firstNames := []string{
		"Alex", "Camille", "Charlie", "Jordan", "Morgan", "Robin", "Sam", "Taylor",
		"Anne-Sophie", "Jean-Pierre", "Louis-Philippe", "Marie-Claire",
	}
	lastNames := []string{"Bernard", "Brown", "Garcia", "Johnson", "Martin", "Miller", "Smith", "Wilson"}
	streetNames := []string{"Cedar Street", "Lake Avenue", "Maple Road", "Oak Street", "Park Avenue", "Pine Road"}
	maidenNames := []string{"Anderson", "Clark", "Davis", "Moore", "Thomas", "White"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.4 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
	}

	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	username := fmt.Sprintf("%s.%s%d", strings.ToLower(firstName), strings.ToLower(lastName), rand.Intn(10000))
	birthday := time.Date(1950+rand.Intn(53), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), 0, 0, 0, 0, time.UTC)

	return &FakeData{
		Name:       firstName + " " + lastName,
		Address:    fmt.Sprintf("%d %s", 1+rand.Intn(9999), streetNames[rand.Intn(len(streetNames))]),
		Birthday:   birthday.Format("2006-01-02"),
		Username:   username,
		Password:   fmt.Sprintf("Demo!%06d", rand.Intn(1000000)),
		EmailU:     strings.ReplaceAll(username, ".", "-"),
		EmailD:     randomEmailDomain(),
		PhoneH:     fmt.Sprintf("+1-555-%03d-%04d", rand.Intn(1000), rand.Intn(10000)),
		Useragent:  userAgents[rand.Intn(len(userAgents))],
		Ipv4:       randomPublicIPv4(),
		MiddleName: maidenNames[rand.Intn(len(maidenNames))],
		CreditCard: fmt.Sprintf("41111111%08d", rand.Intn(100000000)),
	}, nil
}

func randomEmailDomain() string {
	words := []string{"alpine", "atlas", "blue", "cedar", "cloud", "green", "harbor", "maple", "north", "river", "star", "stone"}
	first := words[rand.Intn(len(words))]

	// Generate both single-label domains (atlas.test) and hyphenated domains
	// (blue-harbor.test). The reserved .test suffix prevents real delivery.
	if rand.Intn(2) == 0 {
		return first + ".test"
	}

	second := words[rand.Intn(len(words))]
	for second == first {
		second = words[rand.Intn(len(words))]
	}
	return first + "-" + second + ".test"
}

func randomPublicIPv4() string {
	prefixes := [][3]int{{192, 0, 2}, {198, 51, 100}, {203, 0, 113}}
	prefix := prefixes[rand.Intn(len(prefixes))]
	return fmt.Sprintf("%d.%d.%d.%d", prefix[0], prefix[1], prefix[2], 1+rand.Intn(254))
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
		ZipCode: fmt.Sprintf("%05d", rand.Intn(100000)), // Correct for 5-digit number
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
