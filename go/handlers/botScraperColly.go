package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
)

func HandleColly(c echo.Context) error {

	// Colly collector with attached cookie jar
	collector := colly.NewCollector()

	// Define the URL you want to scrape
	url := "http://juiceshop.corp.fabriclab.ca"

	// On every `mat-card` element which has `product` class
	collector.OnHTML("div.product", func(e *colly.HTMLElement) {
		productName := e.ChildText("div.item-name")
		productPrice := e.ChildText("div.item-price span")

		// Clean up the strings
		productName = strings.TrimSpace(productName)
		productPrice = strings.TrimSpace(productPrice)

		fmt.Printf("Product: %s, Price: %s\n", productName, productPrice)
	})

	// Visit the URL and start scraping
	err := collector.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Colly Actions executed successfully")
}
