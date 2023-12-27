package main

import (
	"darwin2/routes"
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
    e.Use(middleware.CORS())

    // Routes
    routes.Configure(e)

	// Démarrer le serveur sur le port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
