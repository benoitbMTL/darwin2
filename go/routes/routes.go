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
	e.POST("/reset-machine-learning", handlers.HandleResetMachineLearning)
	e.POST("/user-auth", handlers.HandleUserAuth)
	e.POST("/cookie-security", handlers.HandleCookieSecurityAttack)

	// Bot Mitigation
	e.POST("/known-bots", handlers.HandleKnownBots)
	e.POST("/selenium", handlers.HandleSelenium)
	e.POST("/bot-scraper-api", handlers.HandleScrapWithApi)
	e.POST("/bot-page-source", handlers.HandlePageSource)	
	e.POST("/bot-deception", handlers.HandleBotDeception)

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

	// SYSTEM
	e.GET("/config", config.GetConfig)
	e.GET("/list-configs", config.ListConfigs)

	// SAVE, RESET
	e.POST("/save-config", config.SaveConfig)
	e.GET("/reset-config", config.ResetConfig)

	// EXPORT, IMPORT
	e.POST("/export", config.ExportConfig)
	e.POST("/import", config.ImportConfig)

	// BACKUP, RSETORE, DELETE
	e.POST("/backup-local", config.BackupConfigLocal)
	e.POST("/restore-local", config.RestoreConfigLocal)
	e.POST("/delete-local", config.DeleteConfigLocal)

	// HEALTH CHECK
	e.GET("/run-health-check", handlers.HandleHealthCheck)
}
