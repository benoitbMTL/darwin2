package handlers

import (
	"darwin2/config"
	"fmt"
	"math/rand"
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
	sleepForDuration()

	// Maximize the Browser Window
	err := webDriver.MaximizeWindow("")
	if err != nil {
		fmt.Println("Failed to maximize window:", err)
	}

	// Allow time for the next action
	sleepForDuration()

	fmt.Printf("Dismiss Welcome Page & Cookie Message\n")
	dismissWelcomePage(webDriver)

	fmt.Printf("Clicking all products on the landing page\n")
    for i := 0; i < 100; i++ {
			fmt.Printf("clickAllProducts Iteration: %d\n", i)
        clickAllProducts(webDriver)
    }

	fmt.Printf("Creating a new account\n")
	// createAccount(webDriver, credentials)

	fmt.Printf("Logging in\n")
	// login(webDriver, juiceshopUrl, credentials)

	fmt.Printf("Adding a new address\n")
	// addNewAddress(webDriver, juiceshopUrl, address)

	fmt.Printf("Adding a new payment method\n")
	// addNewPayment(webDriver, juiceshopUrl, payment)

	fmt.Printf("Adding items to the shopping cart\n")
	// addItemsToShoppingCart(webDriver, juiceshopUrl)

	fmt.Printf("Checking out the shopping cart\n")
	// checkoutShoppingCart(webDriver, juiceshopUrl)

	time.Sleep(10 * time.Second)

	return nil
}

// DISMISS WELCOME PAGE & COOKIE WARNING ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func dismissWelcomePage(webDriver selenium.WebDriver) {

	sleepForDuration()

	// Find the 'Dismiss' button by its aria-label
	dismissButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Close Welcome Banner"]`)
	if err != nil {
		fmt.Println("Error finding the 'Dismiss' button:", err)
		return
	}
	dismissButton.Click()

	sleepForDuration()

	// Find the link with the text 'Me want it!' and click on it
	meWantItLink, err := webDriver.FindElement(selenium.ByXPATH, `//a[contains(text(), 'Me want it!')]`)
	if err != nil {
		fmt.Println("Error finding the 'Me want it!' link:", err)
		return
	}
	meWantItLink.Click()

	fmt.Println("[+] Done Dismiss Welcome Page & Cookie Message")

	sleepForDuration()
}

// CLICK ALL PRODUCTS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func clickAllProducts(webDriver selenium.WebDriver) {

	sleepForDuration()

	// Find all product elements on the page by their XPath
	products, err := webDriver.FindElements(selenium.ByXPATH, `//div[@aria-label="Click for more information about the product"]`)
	if err != nil {
		fmt.Println("Error finding product elements:", err)
		return
	}

	var productTexts []string
	for index, product := range products {
		fmt.Printf("Processing product %d\n", index+1) // Debug: print the index of the product being processed
		sleepForDuration()

		text, err := product.Text()
		if err != nil {
			fmt.Println("Error getting product text:", err)
			continue
		}

		fmt.Printf("Product name: %s\n", text) // Print the product name
		productTexts = append(productTexts, text)

		// Scroll the product into view using JavaScript
		_, err = webDriver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{product})
		if err != nil {
			fmt.Println("Error scrolling to product:", err)
			continue
		}

		// Click the product
		err = product.Click()
		if err != nil {
			fmt.Println("Error clicking on product:", err)
			continue
		}

		closeButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Close Dialog"]`)
		if err != nil {
			fmt.Println("Error finding the 'Close Dialog' button:", err)
			return
		} else {
			// Click the button to close the dialog
			if err := closeButton.Click(); err != nil {
				fmt.Println("Error clicking the 'Close Dialog' button:", err)
				return
			}
		}
	}

	// Find the 'All Products' title element
	allProductsTitle, err := webDriver.FindElement(selenium.ByXPATH, `//div[contains(text(), 'All Products')]`)
	if err != nil {
		fmt.Println("Error finding the 'All Products' title element:", err)
		return // or handle the error as needed
	}

	// Scroll to the 'All Products' title element
	_, err = webDriver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{allProductsTitle})
	if err != nil {
		fmt.Println("Error scrolling to the 'All Products' title element:", err)
		return // or handle the error as needed
	}

	fmt.Println("[+] Done clicking around on the Startpage")
	sleepForDuration()
}

// CREATE ACCOUNT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createAccount(webDriver selenium.WebDriver, credentials Credentials) {

	// Click Account menu button
	accountMenuButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Show/hide account menu"]`)
	if err != nil {
		fmt.Println("Error finding account menu button:", err)
		return
	}
	accountMenuButton.Click()
	sleepForDuration()

	// Find and click the navbar login button
	loginButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@id="navbarLoginButton"]`)
	if err != nil {
		fmt.Println("Error finding the navbar login button:", err)
		return
	}
	loginButton.Click()
	sleepForDuration()

	// Find and click the 'Not yet a customer?' link
	notCustomerLink, err := webDriver.FindElement(selenium.ByLinkText, "Not yet a customer?")
	if err != nil {
		fmt.Println("Error finding the 'Not yet a customer?' link:", err)
		return
	}
	notCustomerLink.Click()
	sleepForDuration()

	// Find the email input field and send the email from credentials
	emailInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="emailControl"]`)
	if err != nil {
		fmt.Println("Error finding the email input field:", err)
		return
	}
	emailInput.SendKeys(credentials.Email)
	sleepForDuration()

	// Find the password input field and send the password from credentials
	passwordInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="passwordControl"]`)
	if err != nil {
		fmt.Println("Error finding the password input field:", err)
		return
	}
	passwordInput.SendKeys(credentials.Password)
	sleepForDuration()

	// Find the repeat password input field and resend the password
	repeatPasswordInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="repeatPasswordControl"]`)
	if err != nil {
		fmt.Println("Error finding the repeat password input field:", err)
		return
	}
	repeatPasswordInput.SendKeys(credentials.Password)
	sleepForDuration()

	// Find and click the security question dropdown
	securityQuestionDropdown, err := webDriver.FindElement(selenium.ByName, "securityQuestion")
	if err != nil {
		fmt.Println("Error finding the security question dropdown:", err)
		return
	}
	securityQuestionDropdown.Click()
	sleepForDuration()

	// Find and click the specified option in the security question dropdown
	securityQuestionOption, err := webDriver.FindElement(selenium.ByID, "mat-option-3")
	if err != nil {
		fmt.Println("Error finding the security question option:", err)
		return
	}
	securityQuestionOption.Click()
	sleepForDuration()

	// Find the security answer input field and send the security answer from credentials
	securityAnswerInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="securityAnswerControl"]`)
	if err != nil {
		fmt.Println("Error finding the security answer input field:", err)
		return
	}
	securityAnswerInput.SendKeys(credentials.SecAnswer)
	sleepForDuration()

	// Find the 'Register' button by its ID and click on it
	registerButton, err := webDriver.FindElement(selenium.ByID, "registerButton")
	if err != nil {
		fmt.Println("Error finding the 'Register' button:", err)
		return // or handle the error appropriately
	}

	err = registerButton.Click()
	if err != nil {
		fmt.Println("Error clicking on the 'Register' button:", err)
		return // or handle the error appropriately
	}

	fmt.Println("[+] Done Registering a new account")
	sleepForDuration()
}

// LOGIN ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func login(webDriver selenium.WebDriver, url string, credentials Credentials) {

	// Navigate to the login page
	loginURL := url + "/#/login"
	webDriver.Get(loginURL)

	// Fill in the email and password
	emailInput, _ := webDriver.FindElement(selenium.ByXPATH, `//input[@id="email"]`)
	emailInput.SendKeys(credentials.Email)

	passwordInput, _ := webDriver.FindElement(selenium.ByXPATH, `//input[@id="password"]`)
	passwordInput.SendKeys(credentials.Password)

	// Click the login button
	loginButton, _ := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Login"]`)
	loginButton.Click()

	fmt.Println("[+] Done login to the Webshop")
	sleepForDuration()
}

// ADD NEW ADDRESS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewAddress(webDriver selenium.WebDriver, url string, address Address) {
}

// ADD NEW PAYMENT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewPayment(webDriver selenium.WebDriver, url string, payment Payment) {
}

// ADD ITEMS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addItemsToShoppingCart(webDriver selenium.WebDriver, url string) {}

// CHECK OUT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkoutShoppingCart(webDriver selenium.WebDriver, url string) {}

// SLEEP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func sleepForDuration() {

	// Seed the random number generator

	// Generate a random duration between 2 to 4 seconds
	randomDuration := time.Duration(rand.Intn(3)+2) * time.Second
	fmt.Println("Sleeping for", randomDuration)
	time.Sleep(randomDuration)
	fmt.Println("Woke up after", randomDuration)
}
