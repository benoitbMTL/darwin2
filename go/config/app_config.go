package config

import (
	"encoding/json"
	"net/http"
	"sync"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

// AppConfig defines the structure for your application configuration
type AppConfig struct {
	Name            string `json:"name"`
	DVWAURL         string `json:"dvwaUrl"`
	BANKURL         string `json:"bankUrl"`
	JUICESHOPURL    string `json:"juiceShopUrl"`
	PETSTOREURL     string `json:"petStoreUrl"`
	SPEEDTESTURL    string `json:"speedTestUrl"`
	USERNAMEAPI     string `json:"userNameApi"`
	PASSWORDAPI     string `json:"passwordApi"`
	VDOMAPI         string `json:"vdomApi"`
	FWBMGTIP        string `json:"fwbMgtIp"`
	FWBMGTPORT      string `json:"fwbMgtPort"`
	MLPOLICY        string `json:"mlPolicy"`
	USERAGENT       string `json:"userAgent"`
	FABRICLABSTORY  string `json:"fabricLabStory"`
}


var (
	// CurrentConfig holds the current application configuration.
	CurrentConfig string
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
	

	defaultConfig := AppConfig{
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
    
		configsMap["Default"] = defaultConfig
    	CurrentConfig = "Default"
}

// GetCurrentConfig function to retrieve the AppConfig object of the current config
func GetCurrentConfig() AppConfig {
    configMutex.RLock()
    defer configMutex.RUnlock()

    currentConfig, exists := configsMap[CurrentConfig]
    if !exists {
        // Log the error or handle the case where current config does not exist more robustly
        log.Printf("Warning: Configuration '%s' does not exist.", CurrentConfig)
        return AppConfig{} // Return an empty AppConfig as a fallback
    }
    return currentConfig
}

// SetCurrentConfig now correctly sets a string indicating the current configuration key
func SetCurrentConfig(name string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	if _, exists := configsMap[name]; !exists {
		log.Printf("Attempted to set non-existent configuration: %s", name)
		return fmt.Errorf("configuration %s does not exist", name)
	}

	CurrentConfig = name
	log.Printf("Current configuration set to: %s", name)
	return nil
}


// GetConfig handles the GET request for the current configuration
func GetConfig(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	currentConfig, exists := configsMap[CurrentConfig]
	if !exists {
		log.Printf("Requested current configuration '%s' does not exist.", CurrentConfig)
		// Depending on your error handling strategy, you might want to return an error or a default config
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration %s does not exist", CurrentConfig))
	}

	log.Printf("Returning current configuration: %s", CurrentConfig)
	return c.JSON(http.StatusOK, currentConfig)
}

// UpdateConfig handles the POST request to update the configuration
func UpdateConfig(c echo.Context) error {
	var newConfig AppConfig
	if err := c.Bind(&newConfig); err != nil {
		log.Printf("Error binding new configuration: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	// Validate newConfig has a Name
	if newConfig.Name == "" {
		log.Println("New configuration does not have a name")
		return echo.NewHTTPError(http.StatusBadRequest, "Configuration must have a name")
	}

	configsMap[newConfig.Name] = newConfig // Store/Update the configuration by its name
	CurrentConfig = newConfig.Name // Set this as the current active configuration

	log.Printf("Configuration '%s' updated and set as current.", newConfig.Name)
	return c.JSON(http.StatusOK, newConfig)
}


// ResetConfig handles the GET request to reset the configuration
func ResetConfig(c echo.Context) error {
    configMutex.Lock()
    defer configMutex.Unlock()

    // Assuming "Default" is the key for the default configuration in configsMap
    // Ensure that the "Default" configuration exists in configsMap
    if _, exists := configsMap["Default"]; !exists {
        // Handle the case where the default configuration is missing
        log.Printf("Default configuration is missing.")
        return echo.NewHTTPError(http.StatusInternalServerError, "Default configuration is missing.")
    }
    
    // Set CurrentConfig to the key/name of the default configuration
    CurrentConfig = "Default"

    // Optionally, if you want to reset the content of the current configuration to the default values
    configsMap[CurrentConfig] = configsMap["Default"]
    
    log.Printf("Configuration reset to default.")
    return c.JSON(http.StatusOK, configsMap[CurrentConfig])
}

// BackupConfig serializes and sends the current AppConfig struct to JSON
func BackupConfig(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	currentConfig, exists := configsMap[CurrentConfig]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Current configuration not found")
	}

	formattedJSON, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="config_backup.json"`)
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.Blob(http.StatusOK, "application/json", formattedJSON)
}



// RestoreConfig handles the POST request to replace the current configuration
func RestoreConfig(c echo.Context) error {
    var newConfig AppConfig
    if err := c.Bind(&newConfig); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    // Validate or generate a name for newConfig if it doesn't exist
    if newConfig.Name == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Configuration name is required")
    }

    configsMap[newConfig.Name] = newConfig // Update the configuration map with the new configuration
    CurrentConfig = newConfig.Name // Optionally set this new configuration as the current one

    return c.JSON(http.StatusOK, newConfig)
}

func BackupConfigLocal(c echo.Context) error {
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    name := request.Name
    configMutex.Lock()
    configsMap[name] = GetCurrentConfig() // Retrieve the current AppConfig using a correct helper function
    configMutex.Unlock()

    return c.JSON(http.StatusOK, echo.Map{"message": "Configuration '" + name + "' backed up successfully"})
}

func RestoreConfigLocal(c echo.Context) error {
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    name := request.Name

    configMutex.Lock()
    defer configMutex.Unlock()

    // Check if the requested config exists
    config, exists := configsMap[name]
    if !exists {
        return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration '%s' not found", name))
    }

    // Set the requested config as the current configuration
    CurrentConfig = name // Now `CurrentConfig` should be a string that holds the key/name of the current configuration

    // Log and inform the client about the successful operation
    log.Printf("Configuration '%s' restored successfully.", name)
    return c.JSON(http.StatusOK, config) // Return the restored configuration
}

func DeleteConfigLocal(c echo.Context) error {
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    if _, exists := configsMap[request.Name]; !exists {
        return echo.NewHTTPError(http.StatusNotFound, "Configuration '" + request.Name + "' not found")
    }

    delete(configsMap, request.Name) // Delete the configuration by its name/key

    return c.JSON(http.StatusOK, echo.Map{"message": "Configuration '" + request.Name + "' deleted successfully"})
}
