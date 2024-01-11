package handlers

import (
	"darwin2/config"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// Define the structs
type Credentials struct {
	Email     string
	Password  string
	SecAnswer string
}

type Address struct {
	Country string
	Name    string
	Mobile  string
	ZipCode string
	Address string
	City    string
}

type Payment struct {
	Name   string
	CardNr string
	Month  string
	Year   string
}

// Initialize the structs with the given data
var credentials = Credentials{
	Email:     "kapler@netcourrier.com",
	Password:  "Forti123!",
	SecAnswer: "Roy",
}

var address = Address{
	Country: "USA",
	Name:    "Roy",
	Mobile:  "+1234567890",
	ZipCode: "12345",
	Address: "Street where the Bots live 1337",
	City:    "San Francisco",
}

var payment = Payment{
	Name:   "Roy",
	CardNr: "5519010974980551",
	Month:  "10",
	Year:   "2080",
}

// Define chromeDriverPath
const (
	chromeDriverPath = "./selenium/chromedriver" // Path to ChromeDriver
	port             = 4444                      // Port on which ChromeDriver will listen
)

// Define your variable for the number of seconds
var waitDurationInSeconds time.Duration = 1 // seconds

// MAIN START ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleSelenium(c echo.Context) error {
	// Start ChromeDriver service
	fmt.Printf("Start ChromeDriver service\n")
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
	if err := SeleniumActions(webDriver); err != nil {
		return err // Return any errors encountered during Selenium actions
	}

	return nil // No error occurred
}

// ACTIONS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SeleniumActions(webDriver selenium.WebDriver) error {

	// Define JuiceShop URL
	juiceshopUrl := config.CurrentConfig.JUICESHOPURL

	// Navigate to the URL
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

	fmt.Printf("Dismiss Welcome Page & Cookie Message\n")
	dismissWelcomePage(webDriver)

	// fmt.Printf("Clicking all products on the landing page\n")
	// clickAllProducts(webDriver)

	fmt.Printf("Creating a new account\n")
	createAccount(webDriver, credentials)

	fmt.Printf("Logging in\n")
	login(webDriver, juiceshopUrl, credentials)

	fmt.Printf("Adding a new address\n")
	addNewAddress(webDriver, juiceshopUrl, address)

	fmt.Printf("Adding a new payment method\n")
	addNewPayment(webDriver, juiceshopUrl, payment)

	fmt.Printf("Adding items to the shopping cart\n")
	addItemsToShoppingCart(webDriver, juiceshopUrl)

	fmt.Printf("Checking out the shopping cart\n")
	checkoutShoppingCart(webDriver, juiceshopUrl)

	time.Sleep(10 * time.Second)

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

	fmt.Println("### Done Dismiss Welcome Page & Cookie Message")

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
		fmt.Printf("Processing product %d\n", index+1)

		// Scroll the product into view
		scrollIntoView(webDriver, product)

		// Get product name (text)
		text, err := product.Text()
		if err != nil {
			fmt.Println("Error getting product text:", err)
			continue
		}
		fmt.Printf("Product name: %s\n", text)

		// Click the product
		err = product.Click()
		if err != nil {
			fmt.Println("Error clicking on product:", err)
			continue
		}

		// Click the close dialog button
		fmt.Println("Closing product dialog")
		clickButton(webDriver, `//button[@aria-label="Close Dialog"]`)
	}

	fmt.Println("### Done clicking around on the Startpage")
	sleepForShortDuration()
}

// CREATE ACCOUNT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createAccount(webDriver selenium.WebDriver, credentials Credentials) {
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

	fmt.Println("### Done Registering a new account")
}

// LOGIN ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func login(webDriver selenium.WebDriver, url string, credentials Credentials) {

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

	fmt.Printf("### Done login to the Webshop\n")
	sleepForShortDuration()
}

// ADD NEW ADDRESS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewAddress(webDriver selenium.WebDriver, url string, address Address) {
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

	fmt.Println("Clicking 'Submit' to add the new address")
	clickButton(webDriver, `//button[@id="submitButton"]`)

	fmt.Println("### Done adding a new address to the Webshop")
}

// ADD NEW PAYMENT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewPayment(webDriver selenium.WebDriver, url string, payment Payment) {
	// Navigating to the payment methods page
	fmt.Println("Navigating to the payment methods page:", url+"/#/saved-payment-methods")
	if err := webDriver.Get(url + "/#/saved-payment-methods"); err != nil {
		fmt.Println("Error navigating to the payment methods page:", err)
		return
	}
	sleepForShortDuration()

	// Click the button to add a new Card
	fmt.Println("Clicking the button to add a new Card")
	clickButton(webDriver, `//mat-expansion-panel-header[.//mat-panel-title[contains(text(), 'Add new card')]]`)

	// Fill in the payment method details
	fmt.Println("Filling in the payment method details")

/* 	fillInputField(webDriver, `//*[@id="mat-input-1"]`, payment.Name)
	fillInputField(webDriver, `//*[@id="mat-input-2"]`, payment.CardNr)
	selectDropdownOption(webDriver, selenium.ByID, "mat-input-3", selenium.ByXPATH, fmt.Sprintf(`//select[@id="mat-input-3"]/option[@value="%s"]`, payment.Month))
	selectDropdownOption(webDriver, selenium.ByID, "mat-input-4", selenium.ByXPATH, fmt.Sprintf(`//select[@id="mat-input-4"]/option[@value="%s"]`, payment.Year)) */



    // Fill in the Name field
    fillInputField(webDriver, `//input[@id="mat-input-1"]`, payment.Name)

    // Fill in the Card Number field
    fillInputField(webDriver, `//input[@id="mat-input-2"]`, payment.CardNr)

    // Select the Expiry Month
    selectDropdownOption(webDriver, selenium.ByID, "mat-input-3", selenium.ByValue, payment.Month)

    // Select the Expiry Year
    selectDropdownOption(webDriver, selenium.ByID, "mat-input-4", selenium.ByValue, payment.Year)




	// Submit the new payment method
	fmt.Println("Clicking 'Submit' to add the new payment method")
	clickButton(webDriver, `//button[@id="submitButton"]`)

	fmt.Println("### Done adding a new payment method")
}

// ADD ITEMS TO SHOPPING CART ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addItemsToShoppingCart(webDriver selenium.WebDriver, url string) {
	// Navigate to the URL
	fmt.Println("Navigating to the URL:", url)
	if err := webDriver.Get(url); err != nil {
		fmt.Println("Error navigating to the URL:", err)
		return
	}

	sleepForShortDuration()

	// Find all 'Add to Basket' buttons and click them
	fmt.Println("Finding all 'Add to Basket' buttons")
	addToBasketButtons, err := webDriver.FindElements(selenium.ByXPATH, `//button[@aria-label="Add to Basket"]`)
	if err != nil {
		fmt.Println("Error finding 'Add to Basket' buttons:", err)
		return
	}

	for i, button := range addToBasketButtons {
		fmt.Printf("Adding item %d to the shopping cart\n", i+1)
		button.Click()
		sleepForShortDuration()
	}

	fmt.Println("### Done putting all available stuff into shopping cart")
	sleepForShortDuration()
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

	sleepForShortDuration()

	// Select the first address
	fmt.Println("Selecting the first address")
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to payment selection
	clickButton(webDriver, `//button[@aria-label="Proceed to payment selection"]`)

	// Select the first delivery speed
	fmt.Println("Selecting the first delivery speed")
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to delivery method selection
	clickButton(webDriver, `//button[@aria-label="Proceed to delivery method selection"]`)

	// Select the first payment method
	fmt.Println("Selecting the first payment method")
	selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`)

	// Proceed to review
	clickButton(webDriver, `//button[@aria-label="Proceed to review"]`)

	// Complete the purchase
	fmt.Println("Completing the purchase")
	clickButton(webDriver, `//button[@aria-label="Complete your purchase"]`)

	fmt.Println("### Done placing order; selecting address, delivery & payment")

	// Navigate back to the main page
	webDriver.Get(url)
}

// HELPER FUNCTIONS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func sleepForShortDuration() {
	time.Sleep(1 * time.Second)
}

func sleepForLongDuration() {
	time.Sleep(5 * time.Second)
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
	button, err := webDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Println("Error finding button:", xpath, err)
		return
	}
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
