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
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Route to handle the root request
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Welcome to the backend server!")
    })

    // Start the server on port 8080
    e.Logger.Fatal(e.Start(":8080"))
}
