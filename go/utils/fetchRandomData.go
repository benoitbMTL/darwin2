package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
)

// FakeData contains the locally generated identity data used by traffic simulations.
type FakeData struct {
	Name       string `json:"name"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
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

var sampleSequence atomic.Uint64

var firstNames = []string{
	"Alex", "Benoît", "Camille", "Élodie",
	"Jean-Pierre", "Anne-Sophie", "José-Luis", "Marie-Claire",
	"Jean Benoît", "Marie France", "Mary Jane", "Louis Philippe",
}

var lastNames = []string{
	"Martin", "Lévesque", "García", "Smith",
	"Saint-Pierre", "O'Connor", "D'Amours", "Smith-Jones",
	"De la Cruz", "Van Den Berg", "Le Blanc", "Des Rivières",
}

var passwordLetterBlocks = []string{
	"Maple", "River", "Aurora", "Montreal", "Voyage", "Galaxy", "Coffee", "Winter",
}

var passwordSpecialBlocks = []string{"!", "@", "#", "$", "%", "-", "_", ".", "?", "+", "=", "!@"}

// passwordBlockPatterns contains every C/N/S sequence of metric length 3 to 9
// that uses all three categories and never repeats a category in adjacent blocks.
// This produces 1,482 mixed password patterns, so a 3,000-sample run covers
// every pattern at least once regardless of the sequence's starting position.
var passwordBlockPatterns = buildPasswordBlockPatterns()

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
	return generateFakeData(sampleSequence.Add(1) - 1), nil
}

func generateFakeData(sequence uint64) *FakeData {
	index := int(sequence)
	firstName := firstNames[index%len(firstNames)]
	lastName := lastNames[(index/len(firstNames))%len(lastNames)]
	firstAlias := []string{"alex", "benoit", "camille", "elodie", "jean", "anne", "jose", "marie"}[index%8]
	lastAlias := []string{"martin", "levesque", "garcia", "smith", "pierre", "oconnor", "damours", "rivers"}[(index/8)%8]
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.4 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/124.0 Mobile Safari/537.36",
	}

	birthday := time.Date(1950+rand.Intn(53), time.Month(1+rand.Intn(12)), 1+rand.Intn(28), 0, 0, 0, 0, time.UTC)
	emailLocal := emailLocalPart(index, firstAlias, lastAlias)

	return &FakeData{
		Name:       firstName + " " + lastName,
		FirstName:  firstName,
		LastName:   lastName,
		Address:    realisticAddress(index),
		Birthday:   birthday.Format("2006-01-02"),
		Username:   realisticUsername(index, firstAlias, lastAlias),
		Password:   realisticPassword(index),
		EmailU:     emailLocal,
		EmailD:     emailDomain(index / 10),
		PhoneH:     realisticPhone(index),
		Useragent:  userAgents[index%len(userAgents)],
		Ipv4:       randomPublicIPv4(),
		MiddleName: lastNames[(index+5)%len(lastNames)],
		CreditCard: fmt.Sprintf("41111111%08d", rand.Intn(100000000)),
	}
}

func emailLocalPart(index int, first, last string) string {
	digits := fmt.Sprintf("%02d", index%100)
	switch index % 10 {
	case 0:
		return first
	case 1:
		return first + digits
	case 2:
		return first + "." + last
	case 3:
		return first + "." + last + digits
	case 4:
		return first + "-" + digits
	case 5:
		return first + digits + "." + last
	case 6:
		return first + "." + last + "+promo"
	case 7:
		return first + "_" + last + "-" + digits
	case 8:
		return first + digits + "-" + last + digits
	default:
		return first + "-" + digits + "." + last
	}
}

func emailDomain(index int) string {
	words := []string{"atlas", "cedar", "harbor", "maple", "north", "river", "stone", "cloud"}
	first := words[index%len(words)]
	second := words[(index+3)%len(words)]
	digits := fmt.Sprintf("%d", 2+index%97)

	switch index % 10 {
	case 0:
		return first + ".test"
	case 1:
		return first + digits + ".test"
	case 2:
		return first + "-" + second + ".test"
	case 3:
		return digits + ".mail.test"
	case 4:
		return "team-" + digits + ".test"
	case 5:
		return "team-" + digits + ".mail.test"
	case 6:
		return "mail.node" + digits + ".test"
	case 7:
		return "team-" + digits + ".node7.test"
	case 8:
		return first + digits + "-mail.test"
	default:
		return first + digits + "-node7.test"
	}
}

func realisticUsername(index int, first, last string) string {
	digits := fmt.Sprintf("%03d", index%1000)
	variants := []string{
		first,
		first + digits,
		first + "." + last,
		first + "-" + last,
		first + "_" + last,
		first[:1] + last + digits,
		last + "." + first + digits,
		first + "-" + digits + "-" + last,
	}
	return variants[index%len(variants)]
}

func realisticPassword(index int) string {
	pattern := passwordBlockPatterns[index%len(passwordBlockPatterns)]
	var password strings.Builder
	for blockIndex, category := range pattern {
		switch category {
		case 'C':
			word := passwordLetterBlocks[(index+blockIndex)%len(passwordLetterBlocks)]
			if (index+blockIndex)%3 == 0 {
				word = strings.ToUpper(word[:1]) + word[1:]
			}
			password.WriteString(word)
		case 'N':
			width := 2 + (index+blockIndex)%5
			password.WriteString(fmt.Sprintf("%0*d", width, rand.Intn(powerOfTen(width))))
		case 'S':
			password.WriteString(passwordSpecialBlocks[(index+blockIndex)%len(passwordSpecialBlocks)])
		}
	}
	return password.String()
}

func buildPasswordBlockPatterns() []string {
	patterns := make([]string, 0, 1482)
	var visit func(prefix string, targetLength int)
	visit = func(prefix string, targetLength int) {
		if len(prefix) == targetLength {
			if strings.Contains(prefix, "C") && strings.Contains(prefix, "N") && strings.Contains(prefix, "S") {
				patterns = append(patterns, prefix)
			}
			return
		}
		for _, category := range []byte{'C', 'N', 'S'} {
			if len(prefix) == 0 || prefix[len(prefix)-1] != category {
				visit(prefix+string(category), targetLength)
			}
		}
	}
	for length := 3; length <= 9; length++ {
		visit("", length)
	}
	return patterns
}

func powerOfTen(exponent int) int {
	result := 1
	for range exponent {
		result *= 10
	}
	return result
}

func realisticAddress(index int) string {
	number := 1 + index%9999
	unit := 1 + index%400
	addresses := []string{
		fmt.Sprintf("%d Cedar Street", number),
		fmt.Sprintf("%d rue Saint-Paul", number),
		fmt.Sprintf("%dB Maple Road", number),
		fmt.Sprintf("%d avenue du Parc, Apt. %d", number, unit),
		fmt.Sprintf("%d-12 chemin des Érables", number),
		fmt.Sprintf("%d 1/2 boulevard René-Lévesque", number),
		fmt.Sprintf("PO Box %d, Station A", 1000+index%9000),
		fmt.Sprintf("Unit %d, %d North River Drive", unit, number),
	}
	return addresses[index%len(addresses)]
}

func realisticPhone(index int) string {
	area := 200 + index%700
	prefix := 200 + (index*7)%700
	line := (index * 7919) % 10000
	phones := []string{
		fmt.Sprintf("+1-%03d-%03d-%04d", area, prefix, line),
		fmt.Sprintf("+1 %03d %03d %04d", area, prefix, line),
		fmt.Sprintf("(%03d) %03d-%04d", area, prefix, line),
		fmt.Sprintf("%03d.%03d.%04d", area, prefix, line),
		fmt.Sprintf("%03d-%03d-%04d", area, prefix, line),
		fmt.Sprintf("1%03d%03d%04d", area, prefix, line),
	}
	return phones[index%len(phones)]
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
