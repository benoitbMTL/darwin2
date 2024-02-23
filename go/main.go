package main

import (
	"darwin2/config"
	"darwin2/routes"
	"flag"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		// When in debug mode, use the default logger for simplicity
		e.Logger.Info("Debug mode enabled")
		e.Use(middleware.Logger())
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, " +
				"latency=${latency_human}, latency_ns=${latency}, bytes_in=${bytes_in}, " +
				"bytes_out=${bytes_out}, remote_ip=${remote_ip}, error=${error}\n",
		}))
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	routes.Configure(e)

	distPath := "../vue/dist"
	e.Static("/", distPath)

	e.GET("/*", func(c echo.Context) error {
		// Try to find the file in the dist folder.
		reqPath := c.Request().URL.Path
		filePath := filepath.Join(distPath, reqPath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// If the file does not exist, serve index.html
			indexFilePath := filepath.Join(distPath, "index.html")
			return c.File(indexFilePath)
		}
		// For existing files, let the Static middleware handle them.
		return c.File(filePath)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
