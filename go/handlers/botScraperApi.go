package handlers

import (
	"math/rand"
	"crypto/tls"
	"darwin2/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Product represents the structure of a product in the response
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	DeluxePrice float64 `json:"deluxePrice"`
	Image       string  `json:"image"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	DeletedAt   *string `json:"deletedAt"`
}

// Quantity represents the structure of a product's quantity in the response
type Quantity struct {
	ProductId    int    `json:"ProductId"`
	ID           int    `json:"id"`
	Quantity     int    `json:"quantity"`
	LimitPerUser int    `json:"limitPerUser"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// QuantityResponse represents the structure of the quantity API response
type QuantityResponse struct {
	Status string     `json:"status"`
	Data   []Quantity `json:"data"`
}

// ProductResponse represents the structure of the API response
type ProductResponse struct {
	Status string    `json:"status"`
	Data   []Product `json:"data"`
}

// Review represents the structure of a review in the response
type Review struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

// ReviewsResponse represents the structure of the reviews API response
type ReviewsResponse struct {
	Status string   `json:"status"`
	Data   []Review `json:"data"`
}

func HandleScrapWithApi(c echo.Context) error {
	juiceshopUrl := config.CurrentConfig.JUICESHOPURL
	productURL := juiceshopUrl + "/api/products"
	quantityURL := juiceshopUrl + "/api/Quantitys"

	// Generate Random Traffic
	fetchRandomReviewsAndImages(juiceshopUrl)
	log.Printf("Fetching Random Reviews and Images")

	//Fetch robots.txt
	fetchRobotsTxt(juiceshopUrl)
	log.Printf("Fetching robots.txt")

	// GET PRODUCTS
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	productResp, err := client.Get(productURL)
	if err != nil {
		log.Printf("Error fetching product data: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch product data"})
	}
	defer productResp.Body.Close()

	// Example handling for reading response body
	body, err := io.ReadAll(productResp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error reading product response"})
	}

	// Unmarshal the JSON data into the ProductResponse struct
	var productResponse ProductResponse
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		log.Printf("Error unmarshaling product response: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error processing product data"})
	}

	// GET QUANTITIES
	customTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: customTransport}

	quantityResp, err := client.Get(quantityURL)
	if err != nil {
		log.Printf("Error fetching quantity data: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch quantity data"})
	}
	defer quantityResp.Body.Close()

	// Read the response body
	quantityBody, err := io.ReadAll(quantityResp.Body)
	if err != nil {
		log.Printf("Error reading quantity response: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error reading quantity data"})
	}

	// Unmarshal the JSON data into the QuantityResponse struct
	var quantityResponse QuantityResponse
	err = json.Unmarshal(quantityBody, &quantityResponse)
	if err != nil {
		log.Printf("Error unmarshaling quantity response: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error processing quantity data"})
	}

	// GET REVIEWS
	for i, product := range productResponse.Data {
		reviewsURL := fmt.Sprintf("%s/rest/products/%d/reviews", juiceshopUrl, product.ID)
		reviewsResp, err := client.Get(reviewsURL)
		if err != nil {
			log.Printf("Error fetching reviews: %v", err)
			continue // Skip this iteration and continue with the next product
		}
		defer reviewsResp.Body.Close()

		reviewsBody, err := io.ReadAll(reviewsResp.Body)
		if err != nil {
			log.Printf("Error reading reviews response: %v", err)
			// Continue to the next product instead of stopping the process
			continue
		}

		// Unmarshal the JSON data into the ReviewsResponse struct
		var reviewsResponse ReviewsResponse
		err = json.Unmarshal(reviewsBody, &reviewsResponse)
		if err != nil {
			log.Printf("Error unmarshaling reviews response: %v", err)
			// Continue to the next product instead of stopping the process
			continue
		}

		if len(reviewsResponse.Data) > 0 {
			firstReview := reviewsResponse.Data[0]
			productResponse.Data[i].Description += fmt.Sprintf(" - Review: \"%s\" by %s", firstReview.Message, firstReview.Author)
		}
	}

	// Combine product and quantity data
	productQuantityMap := make(map[int]int)
	for _, quantity := range quantityResponse.Data {
		productQuantityMap[quantity.ProductId] = quantity.Quantity
	}

	// Return HTML Table
	return c.String(http.StatusOK, buildHTMLTable(productResponse.Data, productQuantityMap))

}

func buildHTMLTable(products []Product, quantities map[int]int) string {
	var sb strings.Builder

	// Base URL for images
	var juiceshopUrl string
	if len(config.CurrentConfig.FABRICLABSTORY) > 0 {
		juiceshopUrl = "https://juiceshop." + config.CurrentConfig.FABRICLABSTORY + ".fabriclab.ca"
	} else {
		juiceshopUrl = config.CurrentConfig.JUICESHOPURL
	}

	baseURL := juiceshopUrl + "/assets/public/images/products/"

	// Initialize the HTML table
	sb.WriteString("<table style='border-collapse: collapse; width: 100%;'>")
	sb.WriteString("<thead>")
	sb.WriteString("<tr style='background-color: #f2f2f2;'>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>ID</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Image</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Name</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Price</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Quantity</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Description</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Reviews</th>")
	sb.WriteString("</tr>")
	sb.WriteString("</thead>")
	sb.WriteString("<tbody>")

	// Add rows for each product
	for _, product := range products {

		reviews, _ := fetchReviews(product.ID) 
		imageUrl := baseURL + product.Image
		quantity := quantities[product.ID]

		// Product details row
		sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", product.ID))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px; text-align: center;'><img src='%s' alt='%s' style='width:50px; height:50px;'></td>", imageUrl, product.Name))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%s</td>", product.Name))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%.2f</td>", product.Price))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", quantity))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%s</td>", product.Description))

		// Reviews column
		if len(reviews) > 0 {
			sb.WriteString("<td style='border: 1px solid #ddd; padding: 8px;'>")
			for _, review := range reviews {
				sb.WriteString(fmt.Sprintf("<div><strong>%s</strong>: %s</div>", review.Author, review.Message))
			}
			sb.WriteString("</td>")
		} else {
			sb.WriteString("<td style='border: 1px solid #ddd; padding: 8px;'>No reviews</td>")
		}
		sb.WriteString("</tr>")
	}

	sb.WriteString("</tbody></table>")
	return sb.String()
}

// fetchReviews sends a request to the API and returns a list of Reviews
func fetchReviews(productID int) ([]Review, error) {
	juiceshopUrl := config.CurrentConfig.JUICESHOPURL
	reviewsURL := fmt.Sprintf("%s/rest/products/%d/reviews", juiceshopUrl, productID)

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	resp, err := client.Get(reviewsURL)
	if err != nil {
		log.Printf("Error retrieving reviews for product %d: %v", productID, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for a non-200 HTTP status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 status code %d when fetching reviews for product %d", resp.StatusCode, productID)
		return nil, fmt.Errorf("received non-200 status code %d", resp.StatusCode)
	}

	reviewsBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading reviews response: %v", err)
		return nil, err
	}

	// Check for HTML response which starts with '<'
	if len(reviewsBody) > 0 && reviewsBody[0] == '<' {
		log.Printf("Expected JSON response but got HTML for product %d", productID)
		return nil, fmt.Errorf("expected JSON response but got HTML")
	}

	var reviewsResponse ReviewsResponse
	err = json.Unmarshal(reviewsBody, &reviewsResponse)
	if err != nil {
		log.Printf("Error unmarshaling reviews response: %v", err)
		return nil, err
	}

	return reviewsResponse.Data, nil
}


func fetchRandomReviewsAndImages(juiceshopUrl string) {

	// Custom transport for client to allow insecure certificates
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	for i := 0; i < 100; i++ {
		randomID := rand.Intn(100) + 1 // Assuming product IDs are between 1 and 100
		randomImageNumber := rand.Intn(100) + 1

		reviewsURL := fmt.Sprintf("%s/rest/products/%d/reviews", juiceshopUrl, randomID)
		imageURL := fmt.Sprintf("%s/assets/public/images/products/%d.jpg", juiceshopUrl, randomImageNumber)

		// Fetching reviews
		resp, err := client.Get(reviewsURL)
		if err != nil {
			log.Printf("Error fetching reviews for product %d: %v", randomID, err)
			continue
		}
		//log.Printf("Fetched reviews URL: %s, Response Code: %d", reviewsURL, resp.StatusCode)
		resp.Body.Close()

		// Fetching image
		resp, err = client.Get(imageURL)
		if err != nil {
			log.Printf("Error fetching image %d.jpg: %v", randomImageNumber, err)
			continue
		}
		// Reading and then closing the response body
		_, _ = io.ReadAll(resp.Body)
		//log.Printf("Fetched image URL: %s, Response Code: %d", imageURL, resp.StatusCode)
		resp.Body.Close()
	}
}

func fetchRobotsTxt(juiceshopUrl string) {
	// Custom transport for client to allow insecure certificates
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	robotsTxtURL := fmt.Sprintf("%s/robots.txt", juiceshopUrl)

	for i := 0; i < 5; i++ {
		// Fetching /robots.txt
		resp, err := client.Get(robotsTxtURL)
		if err != nil {
			log.Printf("Error fetching /robots.txt: %v", err)
			continue
		}
		//log.Printf("Fetched URL: %s, Response Code: %d", robotsTxtURL, resp.StatusCode)
		resp.Body.Close()
	}
}
