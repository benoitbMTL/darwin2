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
    ID           int     `json:"id"`
    Name         string  `json:"name"`
    Description  string  `json:"description"`
    Price        float64 `json:"price"`
    DeluxePrice  float64 `json:"deluxePrice"`
    Image        string  `json:"image"`
    CreatedAt    string  `json:"createdAt"`
    UpdatedAt    string  `json:"updatedAt"`
    DeletedAt    *string `json:"deletedAt"`
}

// APIResponse represents the structure of the API response
type APIResponse struct {
    Status string    `json:"status"`
    Data   []Product `json:"data"`
}

func HandleScrapWithApi(c echo.Context) error {
    url := "http://juiceshop.corp.fabriclab.ca/rest/products/search?q="

    // Send a GET request to the API
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal("Error fetching data:", err)
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("Error reading response:", err)
    }

    // Unmarshal the JSON data into the APIResponse struct
    var apiResponse APIResponse
    err = json.Unmarshal(body, &apiResponse)
    if err != nil {
        log.Fatal("Error unmarshaling response:", err)
    }

	return c.String(http.StatusOK, buildHTMLTable(apiResponse.Data))
}


func buildHTMLTable(products []Product) string {
    var sb strings.Builder

    // Base URL for images
    baseURL := "http://juiceshop.corp.fabriclab.ca/assets/public/images/products/"

    // Start the HTML table
    sb.WriteString("<table style='border-collapse: collapse; width: 100%;'>")
    sb.WriteString("<tr style='background-color: #f2f2f2;'>")
    sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Image</th>")
    sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>ID</th>")
    sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Name</th>")
    sb.WriteString("<th style='border: 1px solid #ddd; padding: 8px;'>Price</th>")
    sb.WriteString("</tr>")

    // Add rows for each product
    for _, product := range products {
        imageUrl := baseURL + product.Image

        sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px; text-align: center;'><img src='%s' alt='%s' style='width:50px; height:50px;'></td>", imageUrl, product.Name))
        sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", product.ID))
        sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%s</td>", product.Name))
        sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%.2f</td>", product.Price))
        sb.WriteString("</tr>")
    }

    // Close the table
    sb.WriteString("</table>")

    return sb.String()
}
