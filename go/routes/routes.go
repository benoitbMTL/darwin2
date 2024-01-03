package routes

import (
	"darwin2/config"
	"darwin2/handlers"

	"github.com/labstack/echo/v4"
)

func Configure(e *echo.Echo) {

	// Web Protection
	e.POST("/web-scan", handlers.HandleWebScan)
	e.POST("/web-attacks", handlers.HandleWebAttacks)
	e.POST("/traffic-generation", handlers.HandleTrafficGenerator)
	e.POST("/machine-learning", handlers.HandleMachineLearning)
	e.POST("/user-auth", handlers.HandleUserAuth)
	e.POST("/cookie-security", handlers.HandleCookieSecurityAttack)

	// Bot Mitigation
	// API Protection
	e.POST("/api-get", handlers.HandleApiGet)
	e.POST("/api-post", handlers.HandleApiPost)
	e.POST("/api-put", handlers.HandleApiPut)
	e.POST("/api-delete", handlers.HandleApiDelete)
	e.POST("/api-traffic-generation", handlers.HandleApiMachineLearning)

	// REST API
	e.POST("/create-policy", handlers.HandleCreatePolicy)
	e.POST("/delete-policy", handlers.HandleDeletePolicy)

	// System
	e.GET("/config", config.GetConfig)
	e.POST("/config", config.UpdateConfig)
	e.GET("/reset", config.ResetConfig)

}
