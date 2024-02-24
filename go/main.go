package main

import (
	"darwin2/config"
	"darwin2/routes"
	"flag"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var debugMode bool

func init() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
}

func main() {
	flag.Parse() // Parse any command-line flags

	config.Initialize()

	e := echo.New()

	// Enabling debug mode based on the debug flag
	e.Debug = debugMode

	// Adjust logging based on debug mode
	if debugMode {
		e.Logger.Info("Debug mode enabled")
		e.Use(middleware.Logger())
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${remote_ip}, ${host}, ${method}, ${uri}, ${status}\n",
		}))
	}

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	routes.Configure(e)

	distPath := "../vue/dist"
	e.Static("/", distPath)

	e.GET("/*", func(c echo.Context) error {
		reqPath := c.Request().URL.Path
		filePath := filepath.Join(distPath, reqPath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			indexFilePath := filepath.Join(distPath, "index.html")
			return c.File(indexFilePath)
		}
		return c.File(filePath)
	})

	// Start HTTPS server
	go func() {
		e.AutoTLSManager.Cache = autocert.DirCache("/etc/letsencrypt/darwin2")
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	}()

	// Start HTTP server (which redirects to HTTPS)
	e.Logger.Fatal(e.Start(":80"))
}
