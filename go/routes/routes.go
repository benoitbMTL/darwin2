package routes

import (
    "darwin2/handlers"
    "github.com/labstack/echo/v4"
)

func Configure(e *echo.Echo) {
    e.POST("/web-scan", handlers.HandleWebScan)
    e.POST("/web-attacks", handlers.HandleWebAttacks)
    e.POST("/traffic-generation", handlers.HandleTrafficGenerator)
}
