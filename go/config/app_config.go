package config

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// AppConfig defines the structure for your application configuration
type AppConfig struct {
	DVWAURL      		string
	BANKURL      		string
	JUICESHOPURL 		string
	PETSTOREURL  		string
	SPEEDTESTURL 		string
	USERNAMEAPI  		string
	PASSWORDAPI  		string
	VDOMAPI      		string
	FWBMGTIP     		string
	FWBMGTPORT   		string
	MLPOLICY     		string
	USERAGENT			string
	FABRICLABSTORY		string
}

var (
	// CurrentConfig holds the current application configuration
	CurrentConfig AppConfig
	DefaultConfig AppConfig

	// configMutex is used to handle concurrent access to the configuration
	configMutex sync.RWMutex
)

// Initialize sets up the default values for the application configuration
func Initialize() {
	configMutex.Lock()
	defer configMutex.Unlock()

	DefaultConfig = AppConfig{
		DVWAURL:			"https://dvwa.corp.fabriclab.ca",
		BANKURL:			"https://bank.corp.fabriclab.ca/bank.html",
		JUICESHOPURL:		"https://juiceshop.corp.fabriclab.ca",
		PETSTOREURL:		"https://petstore3.corp.fabriclab.ca/api/v3/pet",
		SPEEDTESTURL:		"https://speedtest.corp.fabriclab.ca",
		USERNAMEAPI:		"userapi",
		PASSWORDAPI:		"fortinet123!",
		VDOMAPI:			"root",
		FWBMGTIP:			"fortiweb.corp.fabriclab.ca",
		FWBMGTPORT:  		"443",
		MLPOLICY:     		"DVWA_POLICY",
		USERAGENT:    		"FortiWeb Demo Tool",
		FABRICLABSTORY:		"fortiweb2",
	}

	CurrentConfig = DefaultConfig
}

// GetConfig handles the GET request for the current configuration
func GetConfig(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return c.JSON(http.StatusOK, CurrentConfig)
}

// UpdateConfig handles the POST request to update the configuration
func UpdateConfig(c echo.Context) error {
	var newConfig AppConfig
	if err := c.Bind(&newConfig); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	CurrentConfig = newConfig
	return c.JSON(http.StatusOK, newConfig)
}

// ResetConfig handles the GET request to reset the configuration
func ResetConfig(c echo.Context) error {
	configMutex.Lock()
	defer configMutex.Unlock()
	CurrentConfig = DefaultConfig
	return c.JSON(http.StatusOK, CurrentConfig)
}

// Backup Config
func BackupConfig(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()
  
	// Serialize the CurrentConfig struct to JSON
    formattedJSON, err := json.MarshalIndent(CurrentConfig, "", "  ")
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="config_backup.json"`)
    c.Response().Header().Set(echo.HeaderContentType, "application/json")
    return c.Blob(http.StatusOK, "application/json", formattedJSON)
}


// Restore Config
func RestoreConfig(c echo.Context) error {
  var newConfig AppConfig
  if err := c.Bind(&newConfig); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }

  configMutex.Lock()
  defer configMutex.Unlock()
  CurrentConfig = newConfig
  return c.JSON(http.StatusOK, newConfig)
}
