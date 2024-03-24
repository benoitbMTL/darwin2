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
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	chromeDriverPath = "./selenium/chromedriver-linux64/chromedriver"
	chromePath       = "./selenium/chrome-linux64/chrome"
	port             = 4444
)

type requestParams struct {
	SelectedActions []string `json:"actions"`
	LoopCount       int      `json:"loopCount"`
	IsHeadless      bool     `json:"headless"`
}

// MAIN START ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleSelenium(c echo.Context) error {
	var reqParams requestParams
	var actionErrors []string // To record action errors

	if err := c.Bind(&reqParams); err != nil {
		fmt.Printf("Failed to bind request body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

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
	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr), // Direct logs to STDERR
	}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v\n", err)
		return err
	}

	defer service.Stop()

	// Set general capabilities for Selenium WebDriver
	fmt.Printf("Set general capabilities for Selenium WebDriver\n")
	caps := selenium.Capabilities{"browserName": "chrome"}

	// Define Chrome-specific capabilities with conditional headless argument
	chromeArgs := []string{
		"--no-sandbox",
		"--disable-gpu",
		"--ignore-certificate-errors",
		"--window-size=1280x800",
	}
	if reqParams.IsHeadless {
		chromeArgs = append(chromeArgs, "--headless") // Add headless if true
	}

	// Define Chrome-specific capabilities
	chromeCaps := chrome.Capabilities{
		Path: chromePath,
		Prefs: map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
			// 1 allows notifications (default)
			// 2 blocks all notifications
			// 0 asks user every time (prompt)
		},
		Args: chromeArgs,
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

	// Generate random data for each iteration
	fakeData, err := utils.FetchRandomData()
	if err != nil {
		return err
	}

	credentials := utils.RandomCredentials(*fakeData)
	address := utils.RandomAddress(*fakeData)
	payment := utils.RandomPayment(*fakeData)
	juiceshopUrl := config.CurrentConfig.JUICESHOPURL

	// Navigate to the URL
	fmt.Printf("\n\033[31m----------------------SCRIPT STARTING----------------------\033[0m\n")
	fmt.Printf("Navigate to the URL: %s\n", juiceshopUrl)
	if err := webDriver.Get(juiceshopUrl); err != nil {
		fmt.Println("Failed to load page:", err)
		return err
	}

	// Log the received actions, loop count, and headless mode
	fmt.Printf("\033[35m\n----------------------ACTIONS------------------------------\nLoop Count: %d\n", reqParams.LoopCount)
	fmt.Printf("Headless: %t\n", reqParams.IsHeadless) // Print headless mode status
	for _, action := range reqParams.SelectedActions {
		fmt.Printf("[+] %s\n", action)
	}
	fmt.Println("-----------------------------------------------------------\033[0m")

	// Allow time for the page to load
	sleepForShortDuration()

	// Maximize the Browser Window
	err = webDriver.MaximizeWindow("")
	if err != nil {
		fmt.Println("Failed to maximize window:", err)
	}

	// Allow time for the next action
	sleepForShortDuration()

	// Dismiss Welcome Page - always required
	fmt.Printf("\033[34m\nDismiss Welcome Page & Cookie Message\033[0m\n")
	if err := dismissWelcomePage(webDriver); err != nil {
		actionErrors = append(actionErrors, fmt.Sprintf("dismissWelcomePage: %v", err))
	}

	// Performing selected actions
	for _, action := range reqParams.SelectedActions {
		var err error
		switch action {
		case "clickAllProducts":
			fmt.Printf("\033[34m\nClicking all products on the landing page\033[0m\n")
			err = clickAllProducts(webDriver)
		case "createAccount":
			fmt.Printf("\033[34m\nCreating a new account\033[0m\n")
			err = createAccount(webDriver, credentials)
		case "login":
			fmt.Printf("\033[34m\nLogging in\033[0m\n")
			err = login(webDriver, juiceshopUrl, credentials)
		case "addNewAddress":
			fmt.Printf("\033[34m\nAdding a new address\033[0m\n")
			err = addNewAddress(webDriver, juiceshopUrl, address)
		case "addNewPayment":
			fmt.Printf("\033[34m\nAdding a new payment method\033[0m\n")
			err = addNewPayment(webDriver, juiceshopUrl, payment)
		case "addItemsToShoppingCart":
			fmt.Printf("\033[34m\nAdding items to the shopping cart\033[0m\n")
			err = addItemsToShoppingCart(webDriver, juiceshopUrl)
		case "checkoutShoppingCart":
			fmt.Printf("\033[34m\nChecking out the shopping cart\033[0m\n")
			err = checkoutShoppingCart(webDriver, juiceshopUrl)
		case "logout":
			fmt.Printf("\033[34m\nLogout\033[0m\n")
			err = logout(webDriver)
		}
		if err != nil {
			actionErrors = append(actionErrors, fmt.Sprintf("%s: %v", action, err))
			fmt.Printf("Error executing %s: %v\n", action, err)
		}
	}

	fmt.Printf("\033[31m----------------------END OF SCRIPT----------------------\033[0m\n")

	if len(actionErrors) > 0 {
		return c.JSON(http.StatusOK, "Actions Executed with Some Errors")
	}
	return c.JSON(http.StatusOK, "All Actions Successfully Executed")
}

// DISMISS WELCOME PAGE & COOKIE WARNING ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func dismissWelcomePage(webDriver selenium.WebDriver) error {
	sleepForShortDuration()

	// Click the 'Dismiss' button to close the welcome banner
	fmt.Println("Clicking the 'Dismiss' button to close the welcome banner")
	if err := clickButton(webDriver, `//button[@aria-label="Close Welcome Banner"]`); err != nil {
		return fmt.Errorf("failed to click dismiss button: %w", err)
	}

	// Click the 'Me want it!' link to dismiss any additional messages
	fmt.Println("Clicking the 'Me want it!' link")
	if err := clickButton(webDriver, `//a[contains(text(), 'Me want it!')]`); err != nil {
		return fmt.Errorf("failed to click 'Me want it!' link: %w", err)
	}

	fmt.Println("[+] Done Dismiss Welcome Page & Cookie Message")

	sleepForShortDuration()

	return nil
}

// CLICK ALL PRODUCTS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func clickAllProducts(webDriver selenium.WebDriver) error {
	sleepForShortDuration()

	// Find all product elements on the page
	products, err := webDriver.FindElements(selenium.ByXPATH, `//div[@aria-label="Click for more information about the product"]`)
	if err != nil {
		fmt.Println("Error finding product elements:", err)
		return fmt.Errorf("error finding product elements: %w", err)
	}

	for index, product := range products {

		// Scroll the product into view
		scrollIntoView(webDriver, product)

		fmt.Printf("Processing product %d\n", index+1)

		// Click the product
		err := product.Click()
		if err != nil {
			fmt.Println("Error clicking on product:", err)
			// Decide if you want to stop the loop and return the error, or just log and continue
			// For this example, we'll just log the error and continue with other products
			continue
		}

		// Attempt to click the close dialog button, if this fails, log and move on
		if err := clickButton(webDriver, `//button[@aria-label="Close Dialog"]`); err != nil {
			fmt.Printf("Error closing dialog for product %d: %v\n", index+1, err)
			// Continue with the next product
		}
	}

	fmt.Println("[+] Done clicking around on the Startpage")
	sleepForShortDuration()
	return nil // No errors encountered, or we've chosen not to return early for errors
}

// CREATE ACCOUNT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createAccount(webDriver selenium.WebDriver, credentials utils.Credentials) error {
	// Log actions for debugging
	fmt.Println("Starting account creation process")

	actions := []struct {
		description string
		action      func(selenium.WebDriver) error
	}{
		{"Clicking Account menu button", func(wd selenium.WebDriver) error {
			return clickButton(wd, `//button[@aria-label="Show/hide account menu"]`)
		}},
		{"Clicking the navbar login button", func(wd selenium.WebDriver) error { return clickButton(wd, `//button[@id="navbarLoginButton"]`) }},
		{"Clicking the 'Not yet a customer?' link", func(wd selenium.WebDriver) error {
			return clickButton(wd, `//a[contains(text(), 'Not yet a customer?')]`)
		}},
		{"Filling in the email", func(wd selenium.WebDriver) error {
			return fillInputField(wd, `//input[@id="emailControl"]`, credentials.Email)
		}},
		{"Filling in the password", func(wd selenium.WebDriver) error {
			return fillInputField(wd, `//input[@id="passwordControl"]`, credentials.Password)
		}},
		{"Filling in the repeat password", func(wd selenium.WebDriver) error {
			return fillInputField(wd, `//input[@id="repeatPasswordControl"]`, credentials.Password)
		}},
		{"Clicking the security question dropdown and selecting an option", func(wd selenium.WebDriver) error {
			return selectDropdownOption(wd, selenium.ByName, "securityQuestion", selenium.ByID, "mat-option-3")
		}},
		{"Filling in the security answer", func(wd selenium.WebDriver) error {
			return fillInputField(wd, `//input[@id="securityAnswerControl"]`, credentials.SecAnswer)
		}},
		{"Clicking on 'Register' button", func(wd selenium.WebDriver) error { return clickButton(wd, `//button[@id="registerButton"]`) }},
	}

	for _, act := range actions {
		fmt.Println(act.description)
		if err := act.action(webDriver); err != nil {
			return fmt.Errorf("%s failed: %w", act.description, err)
		}
	}

	fmt.Println("[+] Done Registering a new account")
	return nil
}

// LOGIN ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func login(webDriver selenium.WebDriver, url string, credentials utils.Credentials) error {
	// Navigate to the login page
	fmt.Println("Navigating to the login page.")
	loginURL := url + "/#/login"
	if err := webDriver.Get(loginURL); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Fill in the email and password
	emailInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="email"]`)
	if err != nil {
		return fmt.Errorf("failed to find email input: %w", err)
	}
	emailInput.SendKeys(credentials.Email)

	passwordInput, err := webDriver.FindElement(selenium.ByXPATH, `//input[@id="password"]`)
	if err != nil {
		return fmt.Errorf("failed to find password input: %w", err)
	}
	passwordInput.SendKeys(credentials.Password)

	// Click the login button
	loginButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Login"]`)
	if err != nil {
		return fmt.Errorf("failed to find login button: %w", err)
	}
	if err := loginButton.Click(); err != nil {
		return fmt.Errorf("failed to click login button: %w", err)
	}

	fmt.Println("[+] Successfully logged into the Webshop.")
	sleepForShortDuration()

	return nil
}

// ADD NEW ADDRESS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewAddress(webDriver selenium.WebDriver, url string, address utils.Address) error {
	fmt.Printf("Navigating to %s\n", url+"/#/address/saved")
	if err := webDriver.Get(url + "/#/address/saved"); err != nil {
		return fmt.Errorf("error navigating to the address page: %w", err)
	}
	sleepForShortDuration()

	if err := clickButton(webDriver, `//button[@aria-label="Add a new address"]`); err != nil {
		return fmt.Errorf("failed to click 'Add a new address': %w", err)
	}

	fields := map[string]string{
		`//input[@data-placeholder="Please provide a country."]`:       address.Country,
		`//input[@data-placeholder="Please provide a name."]`:          address.Name,
		`//input[@data-placeholder="Please provide a mobile number."]`: address.Mobile,
		`//input[@data-placeholder="Please provide a ZIP code."]`:      address.ZipCode,
		`//textarea[@id="address"]`:                                    address.Address,
		`//input[@data-placeholder="Please provide a city."]`:          address.City,
	}

	for xpath, value := range fields {
		if err := fillInputField(webDriver, xpath, value); err != nil {
			return fmt.Errorf("failed to fill in the field with xpath %s: %w", xpath, err)
		}
	}

	if err := clickButton(webDriver, `//button[@id="submitButton"]`); err != nil {
		return fmt.Errorf("failed to click 'Submit': %w", err)
	}

	fmt.Println("[+] Successfully added a new address to the Webshop.")
	return nil
}

// ADD NEW PAYMENT ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addNewPayment(webDriver selenium.WebDriver, url string, payment utils.Payment) error {
	data := map[string]string{
		"fullName": payment.Name,
		"cardNum":  payment.CardNr,
		"expMonth": payment.Month,
		"expYear":  payment.Year,
	}

	fmt.Printf("Full Name: %s, Card Number: %s, Expiry Month: %s, Expiry Year: %s\n",
		data["fullName"], data["cardNum"], data["expMonth"], data["expYear"])

	// Assuming the token is stored as a cookie after login
	cookies, err := webDriver.GetCookies()
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	var tokenValue string
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenValue = cookie.Value
			break
		}
	}
	if tokenValue == "" {
		return fmt.Errorf("token not found")
	}

	// Prepare the request
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %w", err)
	}

	req, err := http.NewRequest("POST", url+"/api/Cards", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
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
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("failed to add payment, status: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Payment successfully added")

	// Optional: Navigate to saved payment methods, consider error handling if necessary
	if err := webDriver.Get(url + "/#/saved-payment-methods"); err != nil {
		return fmt.Errorf("failed to navigate to saved payment methods: %w", err)
	}

	fmt.Println("[+] Done adding a new payment method")
	return nil
}

// ADD ITEMS TO SHOPPING CART ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addItemsToShoppingCart(webDriver selenium.WebDriver, url string) error {
	// Navigate to the URL
	fmt.Println("Navigating to the URL:", url)
	if err := webDriver.Get(url); err != nil {
		fmt.Printf("Error navigating to the URL: %v\n", err)
		return err
	}

	sleepForMediumDuration()

	// Find all 'Add to Basket' buttons and click them
	fmt.Println("Finding all 'Add to Basket' buttons")
	addToBasketButtons, err := webDriver.FindElements(selenium.ByXPATH, `//button[@aria-label="Add to Basket"]`)
	if err != nil {
		fmt.Printf("Error finding 'Add to Basket' buttons: %v\n", err)
		return err
	}

	sleepForMediumDuration()

	for i, button := range addToBasketButtons {
		fmt.Printf("Adding item %d to the shopping cart\n", i+1)
		if err := button.Click(); err != nil {
			fmt.Printf("Failed to add item %d to the shopping cart: %v\n", i+1, err)
			// Optionally log the error and continue, or return the error
			// return err
		}
	}

	fmt.Println("[+] Done putting all available stuff into shopping cart")
	return nil
}

// CHECK OUT SHOPPING CART ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkoutShoppingCart(webDriver selenium.WebDriver, url string) error {
	// Click the button to show the shopping cart
	fmt.Println("Clicking the button to show the shopping cart")
	cartButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@aria-label="Show the shopping cart"]`)
	if err != nil {
		fmt.Printf("Error finding the shopping cart button: %v\n", err)
		return err
	}
	if err := cartButton.Click(); err != nil {
		fmt.Printf("Error clicking the shopping cart button: %v\n", err)
		return err
	}

	sleepForLongDuration() // longer sleep due to slow popup

	// Click the checkout button
	fmt.Println("Clicking the checkout button")
	checkoutButton, err := webDriver.FindElement(selenium.ByXPATH, `//button[@id="checkoutButton"]`)
	if err != nil {
		fmt.Printf("Error finding the checkout button: %v\n", err)
		return err
	}
	if err := checkoutButton.Click(); err != nil {
		fmt.Printf("Error clicking the checkout button: %v\n", err)
		return err
	}

	// Select the first address
	fmt.Println("Selecting the first address")
	sleepForMediumDuration()
	if err := selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`); err != nil {
		return err
	}

	// Proceed to payment selection
	sleepForMediumDuration()
	if err := clickButton(webDriver, `//button[@aria-label="Proceed to payment selection"]`); err != nil {
		return err
	}

	// Select the first delivery speed
	sleepForMediumDuration()
	fmt.Println("Selecting the first delivery speed")
	if err := selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`); err != nil {
		return err
	}

	// Proceed to delivery method selection
	sleepForMediumDuration()
	if err := clickButton(webDriver, `//button[@aria-label="Proceed to delivery method selection"]`); err != nil {
		return err
	}

	// Select the first payment method
	fmt.Println("Selecting the first payment method")
	if err := selectFirstOption(webDriver, `//mat-radio-button[@class="mat-radio-button mat-accent"]`); err != nil {
		return err
	}

	// Proceed to review
	if err := clickButton(webDriver, `//button[@aria-label="Proceed to review"]`); err != nil {
		return err
	}

	// Complete the purchase
	fmt.Println("Completing the purchase")
	if err := clickButton(webDriver, `//button[@aria-label="Complete your purchase"]`); err != nil {
		return err
	}

	fmt.Println("[+] Done placing order; selecting address, delivery & payment")

	// Navigate back to the main page
	if err := webDriver.Get(url); err != nil {
		return fmt.Errorf("failed to navigate back to main page: %w", err)
	}

	return nil
}

func logout(webDriver selenium.WebDriver) error {
	// Click Account menu button
	fmt.Println("Clicking Account menu button")
	if err := clickButton(webDriver, `//button[@aria-label="Show/hide account menu"]`); err != nil {
		return err
	}

	// Click the navbar logout button
	fmt.Println("Clicking the navbar logout button")
	if err := clickButton(webDriver, `//button[@id="navbarLogoutButton"]`); err != nil {
		return err
	}

	fmt.Println("[+] Successfully logged out")
	return nil
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

func selectFirstOption(webDriver selenium.WebDriver, xpath string) error {
	options, err := webDriver.FindElements(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Printf("Error finding options with xpath %s: %v\n", xpath, err)
		return err
	}
	if len(options) == 0 {
		return fmt.Errorf("no options found with xpath %s", xpath)
	}
	if err := options[0].Click(); err != nil {
		fmt.Printf("Error clicking first option with xpath %s: %v\n", xpath, err)
		return err
	}
	sleepForShortDuration()
	return nil
}

func clickButton(webDriver selenium.WebDriver, xpath string) error {
	sleepForShortDuration()
	button, err := webDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Error finding button with xpath %s: %v\n", xpath, err)
		return fmt.Errorf("error finding button with xpath %s: %w", xpath, err)
	}
	sleepForShortDuration()

	// Attempt to click the button
	if err := button.Click(); err != nil {
		// Log the error and return it if the click fails
		fmt.Printf("Error clicking on button with xpath %s: %v\n", xpath, err)
		return fmt.Errorf("error clicking on button with xpath %s: %w", xpath, err)
	}

	sleepForShortDuration()
	return nil // Return nil if there was no error
}

func scrollIntoView(webDriver selenium.WebDriver, element selenium.WebElement) {
	_, err := webDriver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{element})
	if err != nil {
		fmt.Println("Error scrolling element into view:", err)
	}
}

func fillInputField(webDriver selenium.WebDriver, xpath string, value string) error {
	field, err := webDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		fmt.Printf("Error finding input field with xpath %s: %v\n", xpath, err)
		return err
	}
	if err := field.SendKeys(value); err != nil {
		fmt.Printf("Error sending keys to field with xpath %s: %v\n", xpath, err)
		return err
	}

	sleepForShortDuration()
	return nil
}

func selectDropdownOption(webDriver selenium.WebDriver, dropdownStrategy string, dropdownValue string, optionStrategy string, optionValue string) error {
	// Click the dropdown to open the list of options
	fmt.Println("Clicking the dropdown to show options")
	dropdown, err := webDriver.FindElement(dropdownStrategy, dropdownValue)
	if err != nil {
		fmt.Printf("Error finding the dropdown: %v\n", err)
		return err
	}
	if err := dropdown.Click(); err != nil {
		fmt.Printf("Error clicking the dropdown: %v\n", err)
		return err
	}
	sleepForShortDuration()

	// Click the specified option within the dropdown
	fmt.Println("Selecting an option from the dropdown")
	option, err := webDriver.FindElement(optionStrategy, optionValue)
	if err != nil {
		fmt.Printf("Error finding the dropdown option: %v\n", err)
		return err
	}
	if err := option.Click(); err != nil {
		fmt.Printf("Error clicking the dropdown option: %v\n", err)
		return err
	}
	sleepForShortDuration()

	return nil
}
