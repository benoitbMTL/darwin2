package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}\t${uri}\t\t\t\t${status}\n",
	}))
	e.Use(middleware.Recover())

	// Serve static files
	e.Static("/", "web/public")


	// Route to handle version request
	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, "(c) Darwin 2.0 (2024)")
	})

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
