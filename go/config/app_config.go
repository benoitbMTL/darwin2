package config

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// AppConfig defines the structure for your application configuration
type AppConfig struct {
	DVWAURL      string
	DVWAHOST     string
	JUICESHOPURL string
	BANKURL      string
	SPEEDTESTURL string
	PETSTOREURL  string
	USERNAMEAPI  string
	PASSWORDAPI  string
	VDOMAPI      string
	FWBMGTIP     string
	FWBMGTPORT   string
	MLPOLICY     string
	USERAGENT    string
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
		DVWAURL:      "https://dvwa.canadaeast.cloudapp.azure.com",
		DVWAHOST:     "dvwa.canadaeast.cloudapp.azure.com",
		JUICESHOPURL: "https://juiceshop.canadaeast.cloudapp.azure.com",
		BANKURL:      "https://dvwa.canadaeast.cloudapp.azure.com/bank.html",
		SPEEDTESTURL: "https://speedtest.canadaeast.cloudapp.azure.com",
		PETSTOREURL:  "https://petstore3.canadaeast.cloudapp.azure.com/api/v3/pet",
		USERNAMEAPI:  "userapi",
		PASSWORDAPI:  "fortinet123!",
		VDOMAPI:      "root",
		FWBMGTIP:     "dvwa.canadaeast.cloudapp.azure.com",
		FWBMGTPORT:   "8443",
		MLPOLICY:     "DVWA_POLICY",
		USERAGENT:    "FortiWeb Demo Tool",
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
