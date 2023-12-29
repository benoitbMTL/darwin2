package routes

import (
    "darwin2/handlers"
    "darwin2/config"
    "github.com/labstack/echo/v4"
)

func Configure(e *echo.Echo) {

    // Web Protection
    e.POST("/web-scan", handlers.HandleWebScan)
    e.POST("/web-attacks", handlers.HandleWebAttacks)
    e.POST("/traffic-generation", handlers.HandleTrafficGenerator)


    // Bot Mitigation
    // API Protection
    // REST API
    // Health Check


    // Configuration
    e.GET("/config", config.GetConfig)
    e.POST("/config", config.UpdateConfig)
    e.GET("/reset", config.ResetConfig)

}
