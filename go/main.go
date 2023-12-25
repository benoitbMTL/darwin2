package main

import (
	// "path/filepath"

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

	// Serve static files from the Vue dist directory
	e.Static("/", "../vue/dist")

	// Catch-all route to serve index.html for any non-static file request
	// e.GET("*", func(c echo.Context) error {
	// 	return c.File(filepath.Join("../vue/dist", "index.html"))
	// })

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
