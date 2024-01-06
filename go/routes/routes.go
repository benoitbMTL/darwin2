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
	e.POST("/known-bots", handlers.HandleKnownBots)


	// API Protection
	e.POST("/api-get", handlers.HandleApiGet)
	e.POST("/api-post", handlers.HandleApiPost)
	e.POST("/api-put", handlers.HandleApiPut)
	e.POST("/api-delete", handlers.HandleApiDelete)
	e.POST("/api-traffic-generation", handlers.HandleApiMachineLearning)

	// REST API CREATE
	e.POST("/create-virtual-ip", handlers.HandleCreateNewVirtualIP)
	e.POST("/create-server-pool", handlers.HandleCreateNewServerPool)
	e.POST("/create-member-pool", handlers.HandleCreateNewMemberPool)
	e.POST("/create-virtual-server", handlers.HandleCreateNewVirtualServer)
	e.POST("/assign-vip-to-virtual-server", handlers.HandleAssignVIPToVirtualServer)
	e.POST("/clone-signature-protection", handlers.HandleCloneSignatureProtection)
	e.POST("/clone-inline-protection", handlers.HandleCloneInlineProtection)
	e.POST("/create-x-forwarded-for-rule", handlers.HandleCreateNewXForwardedForRule)
	e.POST("/configure-protection-profile", handlers.HandleConfigureProtectionProfile)
	e.POST("/create-policy", handlers.HandleCreateNewPolicy)

	// REST API DELETE
	e.POST("/delete-virtual-ip", handlers.HandleDeleteVirtualIP)
	e.POST("/delete-server-pool", handlers.HandleDeleteServerPool)
	e.POST("/delete-virtual-server", handlers.HandleDeleteVirtualServer)
	e.POST("/delete-signature-protection", handlers.HandleDeleteSignatureProtection)
	e.POST("/delete-inline-protection", handlers.HandleDeleteInlineProtection)
	e.POST("/delete-x-forwarded-for-rule", handlers.HandleDeleteXForwardedForRule)
	e.POST("/delete-policy", handlers.HandleDeletePolicy)

	// System
	e.GET("/config", config.GetConfig)
	e.POST("/config", config.UpdateConfig)
	e.GET("/reset", config.ResetConfig)
	e.GET("/health-check", config.HandleHealthCheck)

}
