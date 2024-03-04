package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const Loop = 300

// Define chromeDriverPath
// https://googlechromelabs.github.io/chrome-for-testing/
// https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.94/linux64/chromedriver-linux64.zip
const (
	chromeDriverPath = "./selenium/chromedriver" // Path to ChromeDriver
	port             = 4444                      // Port on which ChromeDriver will listen
)

// Define your variable for the number of seconds
var waitDurationInSeconds time.Duration = 1 // seconds

// MAIN START ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleSelenium(c echo.Context) error {

	// Debug: Check if chromedriver exists
	fmt.Println("Checking if ChromeDriver exists at:", chromeDriverPath)
	if _, err := os.Stat(chromeDriverPath); os.IsNotExist(err) {
		fmt.Println("ChromeDriver not found at:", chromeDriverPath)
		return echo.NewHTTPError(http.StatusNotFound, "ChromeDriver not found")
	} else if err != nil {
		fmt.Println("Error checking ChromeDriver:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking ChromeDriver")
	} else {
		fmt.Println("ChromeDriver found at:", chromeDriverPath)
	}



	// Start ChromeDriver service
	fmt.Printf("\033[1m\033[34m\nStart ChromeDriver service\033[0m\n")
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v\n", err)
		return err
	}
	
	defer service.Stop()

	// Set general capabilities for Selenium WebDriver
	fmt.Printf("Set general capabilities for Selenium WebDriver\n")
	caps := selenium.Capabilities{"browserName": "chrome"}

	// Define Chrome-specific capabilities
	chromeCaps := chrome.Capabilities{
		Prefs: map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
			// 1 allows notifications (default)
			// 2 blocks all notifications
			// 0 asks user every time (prompt)
		},
		Args: []string{
			"--no-sandbox",
			"--ignore-certificate-errors",
		},
	}

	// Merge Chrome-specific capabilities into the general capabilities
	caps.AddChrome(chromeCaps)

	// Create a new instance of the WebDriver, connecting to the ChromeDriver server
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Printf("Error connecting to ChromeDriver: %v\n", err)
		return err
	}
	defer webDriver.Quit() // Ensure that WebDriver is closed when the function returns

	// Execute the Selenium actions defined in SeleniumActions
	for i := 0; i < Loop; i++ {
		// Generate random data for each iteration
		fakeData, err := utils.FetchRandomData()
		if err != nil {
			return err
		}

		credentials := utils.RandomCredentials(*fakeData)
		address := utils.RandomAddress(*fakeData)
		payment := utils.RandomPayment(*fakeData)

		if err := SeleniumActions(webDriver, credentials, address, payment); err != nil {
			return err // Return any errors encountered during Selenium actions
		}
	}

	return nil // No error occurred
}

// ACTIONS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SeleniumActions(webDriver selenium.WebDriver, credentials utils.Credentials, address utils.Address, payment utils.Payment) error {

	// Define JuiceShop URL
	juiceshopUrl := config.CurrentConfig.JUICESHOPURL

	// Navigate to the URL
	fmt.Printf("\n\033[31m----------------------SCRIPT STARTING----------------------\033[0m\n")
	fmt.Printf("Navigate to the URL: %s\n", juiceshopUrl)
	if err := webDriver.Get(juiceshopUrl); err != nil {
		fmt.Println("Failed to load page:", err)
		return err
	}

	// Allow time for the page to load
	sleepForShortDuration()

	// Maximize the Browser Window
	err := webDriver.MaximizeWindow("")
	if err != nil {
		fmt.Println("Failed to maximize window:", err)
	}

	// Allow time for the next action
	sleepForShortDuration()

	fmt.Printf("\033[34m\nDismiss Welcome Page & Cookie Message\033[0m\n")
	dismissWelcomePage(webDriver)

	if rand.Intn(2) == 0 {
		fmt.Printf("\033[34m\nClicking all products on the landing page\033[0m\n")
		clickAllProducts(webDriver)
	}

	fmt.Printf("\033[34m\nCreating a new account\033[0m\n")
	createAccount(webDriver, credentials)

	fmt.Printf("\033[34m\nLogging in\033[0m\n")
	login(webDriver, juiceshopUrl, credentials)

	fmt.Printf("\033[34m\nAdding a new address\033[0m\n")
	addNewAddress(webDriver, juiceshopUrl, address)

	fmt.Printf("\033[34m\nAdding a new payment method\033[0m\n")
	addNewPayment(webDriver, juiceshopUrl, payment)

	fmt.Printf("\033[34m\nAdding items to the shopping cart\033[0m\n")
	addItemsToShoppingCart(webDriver, juiceshopUrl)

	fmt.Printf("\033[34m\nChecking out the shopping cart\033[0m\n")
	checkoutShoppingCart(webDriver, juiceshopUrl)

	fmt.Printf("\033[34m\nLogout\033[0m\n")
	logout(webDriver)

	fmt.Printf("\033[31m----------------------END OF SCRIPT----------------------\033[0m\n")

	sleepForLongDuration()

	return nil
}

// DISMISS WELCOME PAGE & COOKIE WARNING ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func dismissWelcomePage(webDriver selenium.WebDriver) {
	sleepForShortDuration()

	// Click the 'Dismiss' button to close the welcome banner
	fmt.Println("Clicking the 'Dismiss' button to close the welcome banner")
	clickButton(webDriver, `//button[@aria-label="Close Welcome Banner"]`)

	// Click the 'Me want it!' link to dismiss any additional messages
	fmt.Println("Clicking the 'Me want it!' link")
	clickButton(webDriver, `//a[contains(text(), 'Me want it!')]`)

	fmt.Println("[+] Done Dismiss Welcome Page & Cookie Message")

	sleepForShortDuration()
}

// CLICK ALL PRODUCTS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func clickAllProducts(webDriver selenium.WebDriver) {
	sleepForShortDuration()

	// Find all product elements on the page
	products, err := webDriver.FindElements(selenium.ByXPATH, `//div[@aria-label="Click for more information about the product"]`)
	if err != nil {
		fmt.Println("Error finding product elements:", err)
		return
	}

	for index, product := range products {

		// Scroll the product into view
		scrollIntoView(webDriver, product)

		fmt.Printf("Processing product %d\n", index+1)

		// Click the product
		err = product.Click()
		if err != nil {
			fmt.Println("Error clicking on product:", err)
			continue
		}

		// Click the close dialog button
		clickButton(webDriver, `//button[@aria-label="Close Dialog"]`)
	}

	fmt.Println("[+] Done clicking around on the Startpage")
	sleepForShortDuration()
}

// CREATE ACCOUNT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createAccount(webDriver selenium.WebDriver, credentials utils.Credentials) {
	// Click Account menu button
	fmt.Println("Clicking Account menu button")
	clickButton(webDriver, `//button[@aria-label="Show/hide account menu"]`)

	// Click the navbar login button
	fmt.Println("Clicking the navbar login button")
	clickButton(webDriver, `//button[@id="navbarLoginButton"]`)

	// Click the 'Not yet a customer?' link
	fmt.Println("Clicking the 'Not yet a customer?' link")
	clickButton(webDriver, `//a[contains(text(), 'Not yet a customer?')]`)

	// Fill in the email, password, repeat password, and security answer fields
	fmt.Printf("Filling in the email: %s\n", credentials.Email)
	fillInputField(webDriver, `//input[@id="emailControl"]`, credentials.Email)

	fmt.Printf("Filling in the password: %s\n", credentials.Password)
	fillInputField(webDriver, `//input[@id="passwordControl"]`, credentials.Password)

	fmt.Printf("Filling in the repeat password: %s\n", credentials.Password)
	fillInputField(webDriver, `//input[@id="repeatPasswordControl"]`, credentials.Password)

	// Click the security question dropdown and select an option

	fmt.Println("Clicking the security question dropdown and selecting an option")
	selectDropdownOption(webDriver, selenium.ByName, "securityQuestion", selenium.ByID, "mat-option-3")

	fmt.Printf("Filling in the security answer: %s\n", credentials.SecAnswer)
	fillInputField(webDriver, `//input[@id="securityAnswerControl"]`, credentials.SecAnswer)

	// Click the 'Register' button
	fmt.Println("Clicking on 'Register' button")
	clickButton(webDriver, `//button[@id="registerButton"]`)

	fmt.Println("[+] Done Registering a new account")
}

// LOGIN ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func login(webDriver selenium.WebDriver, url string, credentials utils.Credentials) {

	// Navigate to the login page
	fmt.Printf("Navigate to the login page.\n")
	loginURL := url + "/#/login"
	webDriver.Get(loginURL)

	// Fill in the email and password
	fmt.Printf("Fill in the email.\n")
	emailInput, _ := webDriver.FindElement(selenium.ByXPATH, `//input[@id="email"]`)
	emailInput.SendKeys(credentials.Email)

	fmt.Printf("Fill in the password.\n")
	passwordInput, _ := webDriver.FindElement(selenium.ByXPATH, `//input[@id="password"]`)
	passwordInput.SendKeys(credentials.Password)

	// Click the login button
	fmt.Printf("Click the login button.\n")
	loginButton, _ := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Login"]`)
	loginButton.Click()

	fmt.Printf("[+] Done login to the Webshop\n")
	sleepForShortDuration()
}

// ADD NEW ADDRESS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewAddress(webDriver selenium.WebDriver, url string, address utils.Address) {
	fmt.Printf("Navigating to %s\n", url+"/#/address/saved")
	if err := webDriver.Get(url + "/#/address/saved"); err != nil {
		fmt.Println("Error navigating to the address page:", err)
		return
	}
	sleepForShortDuration()

	fmt.Println("Clicking 'Add a new address'")
	clickButton(webDriver, `//button[@aria-label="Add a new address"]`)

	// Fill in the address details
	fmt.Println("Filling in the address details")
	fillInputField(webDriver, `//input[@data-placeholder="Please provide a country."]`, address.Country)
	fillInputField(webDriver, `//input[@data-placeholder="Please provide a name."]`, address.Name)
	fillInputField(webDriver, `//input[@data-placeholder="Please provide a mobile number."]`, address.Mobile)
	fillInputField(webDriver, `//input[@data-placeholder="Please provide a ZIP code."]`, address.ZipCode)
	fillInputField(webDriver, `//textarea[@id="address"]`, address.Address)
	fillInputField(webDriver, `//input[@data-placeholder="Please provide a city."]`, address.City)
	fmt.Printf("Country: %s, Name: %s, Mobile: %s, ZipCode: %s, Address: %s, City: %s\n",
		address.Country, address.Name, address.Mobile, address.ZipCode, address.Address, address.City)

	fmt.Println("Clicking 'Submit' to add the new address")
	clickButton(webDriver, `//button[@id="submitButton"]`)

	fmt.Println("[+] Done adding a new address to the Webshop")
}

// ADD NEW PAYMENT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewPayment(webDriver selenium.WebDriver, url string, payment utils.Payment) {
	data := map[string]string{
		"fullName": payment.Name,
		"cardNum":  payment.CardNr,
		"expMonth": payment.Month,
		"expYear":  payment.Year,
	}

	fmt.Printf("Full Name: %s, Card Number: %s, Expiry Month: %s, Expiry Year: %s\n",
		data["fullName"], data["cardNum"], data["expMonth"], data["expYear"])

	// Assuming the token is stored as a cookie after login
	cookies, _ := webDriver.GetCookies()
	var tokenValue string
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenValue = cookie.Value
			break
		}
	}

	// Prepare the request
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url+"/api/Cards", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+tokenValue)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Prepare the HTTP client with server certificate verification disabled
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] HTTP request failed: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == 201 {
		fmt.Println("Payment successfully added")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[-] ERROR Status: %d\n", resp.StatusCode)
		fmt.Printf("[-] Error Text:\n%s\n", string(body))
	}

	// Navigate to saved payment methods
	webDriver.Get(url + "/#/saved-payment-methods")
	fmt.Println("[+] Done adding a new payment method")
	sleepForShortDuration()
}

// ADD ITEMS TO SHOPPING CART ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addItemsToShoppingCart(webDriver selenium.WebDriver, url string) {
	// Navigate to the URL
	fmt.Println("Navigating to the URL:", url)
	if err := webDriver.Get(url); err != nil {
		fmt.Println("Error navigating to the URL:", err)
		return
	}

	sleepForMediumDuration()

	// Find all 'Add to Basket' buttons and click them
	fmt.Println("Finding all 'Add to Basket' buttons")
	addToBasketButtons, err := webDriver.FindElements(selenium.ByXPATH, `//button[@aria-label="Add to Basket"]`)
	if err != nil {
		fmt.Println("Error finding 'Add to Basket' buttons:", err)
		return
	}

	sleepForMediumDuration()

	for i, button := range addToBasketButtons {
		fmt.Printf("Adding item %d to the shopping cart\n", i+1)
		button.Click()
	}

	fmt.Println("[+] Done putting all available stuff into shopping cart")

}

// CHECK OUT SHOPPING CART ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkoutShoppingCart(webDriver selenium.WebDriver, url string) {
	// Click the button to show the shopping cart
	fmt.Println("Clicking the button to show the shopping cart")
	cartButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Show the shopping cart"]`)
	if err != nil {
		fmt.Println("Error finding the shopping cart button:", err)
		return
	}
	cartButton.Click()

	sleepForLongDuration() // longer sleep due to slow popup

	// Click the checkout button
	fmt.Println("Clicking the checkout button")
	checkoutButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@id="checkoutButton"]`)
	if err != nil {
		fmt.Println("Error finding the checkout button:", err)
		return
	}
	checkoutButton.Click()

	// Select the first address
	fmt.Println("Selecting the first address")
	sleepForMediumDuration()
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to payment selection
	sleepForMediumDuration()
	clickButton(webDriver, `//button[@aria-label="Proceed to payment selection"]`)

	// Select the first delivery speed
	sleepForMediumDuration()
	fmt.Println("Selecting the first delivery speed")
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to delivery method selection
	sleepForMediumDuration()
	clickButton(webDriver, `//button[@aria-label="Proceed to delivery method selection"]`)

	// Select the first payment method
	fmt.Println("Selecting the first payment method")
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to review
	clickButton(webDriver, `//button[@aria-label="Proceed to review"]`)

	// Complete the purchase
	fmt.Println("Completing the purchase")
	clickButton(webDriver, `//button[@aria-label="Complete your purchase"]`)

	fmt.Println("[+] Done placing order; selecting address, delivery & payment")

	// Navigate back to the main page
	webDriver.Get(url)
}

func logout(webDriver selenium.WebDriver) {
	// Click Account menu button
	fmt.Println("Clicking Account menu button")
	clickButton(webDriver, `//button[@aria-label="Show/hide account menu"]`)

	// Click the navbar login button
	fmt.Println("Clicking the navbar logout button")
	clickButton(webDriver, `//button[@id="navbarLogoutButton"]`)

	fmt.Println("[+] Successfully logged out")

}

// HELPER FUNCTIONS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func sleepForShortDuration() {
	// Random duration between 100 and 900 milliseconds
	duration := rand.Intn(800) + 100 // 100 to 900
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func sleepForMediumDuration() {
	// Random duration between 1 and 5 seconds
	duration := rand.Intn(4) + 1 // 1 to 5
	time.Sleep(time.Duration(duration) * time.Second)
}

func sleepForLongDuration() {
	// Random duration between 5 and 10 seconds
	duration := rand.Intn(5) + 5 // 5 to 10
	time.Sleep(time.Duration(duration) * time.Second)
}

func selectFirstOption(webDriver selenium.WebDriver, xpath string) {
	options, err := webDriver.FindElements(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Println("Error finding options:", err)
		return
	}
	options[0].Click()
	sleepForShortDuration()
}

func clickButton(webDriver selenium.WebDriver, xpath string) {
	sleepForShortDuration()
	button, err := webDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Println("Error finding button:", xpath, err)
		return
	}
	sleepForShortDuration()
	button.Click()
	sleepForShortDuration()
}

func scrollIntoView(webDriver selenium.WebDriver, element selenium.WebElement) {
	_, err := webDriver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{element})
	if err != nil {
		fmt.Println("Error scrolling element into view:", err)
	}
}

func fillInputField(webDriver selenium.WebDriver, xpath string, value string) {
	field, err := webDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Println("Error finding input field:", err)
		return
	}
	field.SendKeys(value)

	sleepForShortDuration()
}

func selectDropdownOption(webDriver selenium.WebDriver, dropdownStrategy string, dropdownValue string, optionStrategy string, optionValue string) {
	// Click the dropdown to open the list of options
	fmt.Println("Clicking the dropdown to show options")
	dropdown, err := webDriver.FindElement(dropdownStrategy, dropdownValue)
	if err != nil {
		fmt.Println("Error finding the dropdown:", err)
		return
	}
	dropdown.Click()
	sleepForShortDuration()

	// Click the specified option within the dropdown
	fmt.Println("Selecting an option from the dropdown")
	option, err := webDriver.FindElement(optionStrategy, optionValue)
	if err != nil {
		fmt.Println("Error finding the dropdown option:", err)
		return
	}
	option.Click()
	sleepForShortDuration()
}
