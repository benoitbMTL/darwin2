package handlers

import (
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

func HandleScrapWithApi(c echo.Context) error {
	productURL := "http://juiceshop.canadaeast.cloudapp.azure.com/api/products"
	quantityURL := "http://juiceshop.canadaeast.cloudapp.azure.com/api/Quantitys"

	// GET PRODUCTS
	productResp, err := http.Get(productURL)
	if err != nil {
		log.Fatal("Error fetching data:", err)
	}
	defer productResp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(productResp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Unmarshal the JSON data into the ProductResponse struct
	var productResponse ProductResponse
	err = json.Unmarshal(body, &productResponse)
	if err != nil {
		log.Fatal("Error unmarshaling response:", err)
	}

	// GET QUANTITIES
	quantityResp, err := http.Get(quantityURL)
	if err != nil {
		log.Fatal("Error fetching quantity data:", err)
	}
	defer quantityResp.Body.Close()

	quantityBody, err := io.ReadAll(quantityResp.Body)
	if err != nil {
		log.Fatal("Error reading quantity response:", err)
	}

	var quantityResponse QuantityResponse
	err = json.Unmarshal(quantityBody, &quantityResponse)
	if err != nil {
		log.Fatal("Error unmarshaling quantity response:", err)
	}

	// Combine product and quantity data
	productQuantityMap := make(map[int]int)
	for _, quantity := range quantityResponse.Data {
		productQuantityMap[quantity.ProductId] = quantity.Quantity
	}

	return c.String(http.StatusOK, buildHTMLTable(productResponse.Data, productQuantityMap))

}

func buildHTMLTable(products []Product, quantities map[int]int) string {
	var sb strings.Builder

	// Base URL for images
	baseURL := "http://juiceshop.canadaeast.cloudapp.azure.com/assets/public/images/products/"

	// Start the HTML table
	sb.WriteString("<table style='border-collapse: collapse; width: 100%;'>")
	sb.WriteString("<tr style='background-color: #f2f2f2;'>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Image</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>ID</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Name</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Price</th>")
	sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Quantity</th>")
	sb.WriteString("</tr>")

	// Add rows for each product
	for _, product := range products {
		imageUrl := baseURL + product.Image
		quantity := quantities[product.ID]

		sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px; text-align: center;'><img src='%s' alt='%s' style='width:50px; height:50px;'></td>", imageUrl, product.Name))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", product.ID))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%s</td>", product.Name))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%.2f</td>", product.Price))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", quantity))
		sb.WriteString("</tr>")
	}

	// Close the table
	sb.WriteString("</table>")

	return sb.String()
}
