package main

import (
	"darwin2/config"  // Import the config package
	"darwin2/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func main() {
	// Initialize the application configuration
	config.Initialize()

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method}, ${uri}, ${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	routes.Configure(e)

	// Vue Production Build
	e.Static("/", "vue/dist") // Adjust the path to your Vue app's dist directory as necessary

	// Catch-all Route to Serve index.html
	e.GET("/*", func(c echo.Context) error {
		indexPath := filepath.Join("vue/dist", "index.html") // Adjust the path as necessary
		indexFile, err := ioutil.ReadFile(indexPath)
		if err != nil {
			return err // Or return an appropriate error in your context
		}
		return c.HTMLBlob(http.StatusOK, indexFile)
	})

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
