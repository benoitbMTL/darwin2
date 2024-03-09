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
	// CurrentConfig holds the current application configuration.
	CurrentConfig AppConfig
	// DefaultConfig holds the default application configuration.
	DefaultConfig AppConfig
	// configMutex is used to handle concurrent access to the configuration.
	configMutex sync.RWMutex
	// configsMap stores configurations by name.
	configsMap = make(map[string]AppConfig)
)

// Initialize sets up the default values for the application configuration
func Initialize() {
	configMutex.Lock()
	defer configMutex.Unlock()

    configsMap["Default"] = AppConfig{
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

    configsMap["AzureConfig"] = AppConfig{
		DVWAURL:			"https://dvwa.canadaeast.cloudapp.azure.com/",
		BANKURL:			"https://bank.canadaeast.cloudapp.azure.com/bank.html",
		JUICESHOPURL:		"https://juiceshop.canadaeast.cloudapp.azure.com",
		PETSTOREURL:		"https://petstore3.canadaeast.cloudapp.azure.com/api/v3/pet",
		SPEEDTESTURL:		"https://speedtest.canadaeast.cloudapp.azure.com",
		USERNAMEAPI:		"userapi",
		PASSWORDAPI:		"fortinet123!",
		VDOMAPI:			"root",
		FWBMGTIP:			"fortiweb.canadaeast.cloudapp.azure.com",
		FWBMGTPORT:  		"8443",
		MLPOLICY:     		"DVWA_POLICY",
		USERAGENT:    		"FortiWeb Demo Tool",
		FABRICLABSTORY:		"",    }

    configsMap["FabricLabConfig"] = AppConfig{
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
		FABRICLABSTORY:		"fortiweb",    }

    configsMap["FortiWebCloudConfig"] = AppConfig{
		DVWAURL:			"https://dvwa.96859.fortiwebcloud.net",
		BANKURL:			"https://bank.96859.fortiwebcloud.net/bank.html",
		JUICESHOPURL:		"https://juiceshop.96859.fortiwebcloud.net",
		PETSTOREURL:		"https://petstore3.96859.fortiwebcloud.net/api/v3/pet",
		SPEEDTESTURL:		"",
		USERNAMEAPI:		"",
		PASSWORDAPI:		"",
		VDOMAPI:			"",
		FWBMGTIP:			"",
		FWBMGTPORT:  		"",
		MLPOLICY:     		"",
		USERAGENT:    		"FortiWeb Demo Tool",
		FABRICLABSTORY:		"",    }
    
    // Set the current configuration to Default
    CurrentConfig = "Default"
}

// GetCurrentConfig returns the currently active AppConfig
func GetCurrentConfig() AppConfig {
    configMutex.RLock()
    defer configMutex.RUnlock()
    
    return configsMap[CurrentConfig]
}

// SetCurrentConfig changes the current configuration to the one specified by name
func SetCurrentConfig(name string) error {
    configMutex.Lock()
    defer configMutex.Unlock()
    
    if _, exists := configsMap[name]; !exists {
        return fmt.Errorf("configuration %s does not exist", name)
    }
    
    CurrentConfig = name
    return nil
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


func BackupConfigLocal(c echo.Context) error {
    name := "" // Extract this from request body
    if err := c.Bind(&name); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid name provided")
    }

    configMutex.Lock()
    configsMap[name] = CurrentConfig
    configMutex.Unlock()

    return c.JSON(http.StatusOK, echo.Map{"message": "Configuration backed up successfully"})
}

func RestoreConfigLocal(c echo.Context) error {
    name := "" // Extract this from request body
    if err := c.Bind(&name); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid name provided")
    }

    configMutex.Lock()
    config, exists := configsMap[name]
    if !exists {
        configMutex.Unlock()
        return echo.NewHTTPError(http.StatusNotFound, "Configuration not found")
    }

    CurrentConfig = config
    configMutex.Unlock()

    return c.JSON(http.StatusOK, CurrentConfig)
}

func DeleteConfigLocal(c echo.Context) error {
    // Assume we're receiving the name of the config to delete in the request body
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    // Check if the configuration exists
    if _, exists := configsMap[request.Name]; !exists {
        return echo.NewHTTPError(http.StatusNotFound, "Configuration not found")
    }

    // Delete the configuration from the map
    delete(configsMap, request.Name)

    return c.JSON(http.StatusOK, echo.Map{"message": "Configuration deleted successfully"})
}
