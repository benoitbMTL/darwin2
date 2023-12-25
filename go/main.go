package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}, ${uri}, ${status}\n",
	}))
	e.Use(middleware.Recover())

    // Serve static files from the Vue public directory
    e.Static("/", "../vue/dist")

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
