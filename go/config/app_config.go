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
	CurrentConfig 	AppConfig                       // Holds the current application configuration
	DefaultConfig 	AppConfig                       // Holds the default application configuration
	configsMap    	= make(map[string]AppConfig)    // Stores configurations by name
	currentName   	string                          // Holds the name of the current configuration
	configMutex		sync.RWMutex
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
    
		// Add your defaultConfig to the configsMap
		configsMap[defaultConfig.Name] = defaultConfig

		currentName = "Default"

        // Ensure that CurrentConfig reflects the configuration pointed to by currentName.
        var exists bool
        CurrentConfig, exists = configsMap[currentName]
        if !exists {
            log.Fatalf("Fatal error: default configuration '%s' does not exist.", currentName)
            // Optionally, handle this more gracefully than just crashing.
        }
}

func GetCurrentConfig() AppConfig {
    configMutex.RLock()
    defer configMutex.RUnlock()

    // Use currentName to access the current configuration from configsMap
    currentConfig, exists := configsMap[currentName]
    if !exists {
        log.Printf("Warning: Configuration '%s' does not exist.", currentName)
        // This could happen if the configuration is deleted or not initialized properly
        // You may return a default configuration or handle this case as needed
        return DefaultConfig
    }
    return currentConfig
}

// GetConfig handles the GET request for the current configuration
func GetConfig(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Use currentName to fetch the current configuration from configsMap
	currentConfig, exists := configsMap[currentName]
	if !exists {
		log.Printf("Requested current configuration '%s' does not exist.", currentName)
		// Return an HTTP error with a message indicating the configuration does not exist
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration '%s' does not exist", currentName))
	}

	log.Printf("Returning current configuration: %s", currentName)
	// Return the current configuration as a JSON response
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

    // Assuming newConfig includes a 'Name' field to identify the configuration
    if newConfig.Name == "" {
        log.Println("Error: Configuration must have a name.")
        return echo.NewHTTPError(http.StatusBadRequest, "Configuration must have a name.")
    }

    // Update or add the new configuration in the map
    configsMap[newConfig.Name] = newConfig

    // Optionally, set this new configuration as the current one
    currentName = newConfig.Name

    log.Printf("Configuration '%s' updated and set as current.", newConfig.Name)
    return c.JSON(http.StatusOK, newConfig)
}


func ResetConfig(c echo.Context) error {
    configMutex.Lock()
    defer configMutex.Unlock()

    // Check if the default configuration exists
    defaultConfig, exists := configsMap["Default"]
    if !exists {
        log.Println("Default configuration is missing.")
        return echo.NewHTTPError(http.StatusInternalServerError, "Default configuration is missing.")
    }
    
    // Reset the current configuration to the default configuration
    currentName = "Default"
    log.Println("Configuration reset to default.")
    return c.JSON(http.StatusOK, defaultConfig)
}


func BackupConfig(c echo.Context) error {
    configMutex.RLock()
    defer configMutex.RUnlock()

    currentConfig, exists := configsMap[currentName]
    if !exists {
        return echo.NewHTTPError(http.StatusNotFound, "Current configuration not found.")
    }

    formattedJSON, err := json.MarshalIndent(currentConfig, "", "  ")
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="config_backup.json"`)
    c.Response().Header().Set(echo.HeaderContentType, "application/json")
    return c.Blob(http.StatusOK, "application/json", formattedJSON)
}



func RestoreConfig(c echo.Context) error {
    var newConfig AppConfig
    if err := c.Bind(&newConfig); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    if newConfig.Name == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "Configuration name is required.")
    }

    configsMap[newConfig.Name] = newConfig
    currentName = newConfig.Name
    log.Printf("Configuration '%s' restored and set as current.", newConfig.Name)
    return c.JSON(http.StatusOK, newConfig)
}

func BackupConfigLocal(c echo.Context) error {
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    // Ensure not to overwrite an existing configuration unless intended
    if _, exists := configsMap[request.Name]; exists {
        return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Configuration '%s' already exists", request.Name))
    }

    configsMap[request.Name] = GetCurrentConfig() // Use GetCurrentConfig to ensure consistency
    log.Printf("Configuration '%s' backed up successfully.", request.Name)

    return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' backed up successfully", request.Name)})
}

func RestoreConfigLocal(c echo.Context) error {
    var request struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&request); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    configMutex.Lock()
    defer configMutex.Unlock()

    _, exists := configsMap[request.Name]
    if !exists {
        return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration '%s' not found", request.Name))
    }

    currentName = request.Name // Correctly update the currentName
    log.Printf("Configuration '%s' restored successfully.", request.Name)

    return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' restored successfully", request.Name)})
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
        return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration '%s' not found", request.Name))
    }

    // Prevent deletion of the currently active configuration
    if request.Name == currentName {
        return echo.NewHTTPError(http.StatusBadRequest, "Cannot delete the currently active configuration")
    }

    delete(configsMap, request.Name)
    log.Printf("Configuration '%s' deleted successfully.", request.Name)

    return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' deleted successfully", request.Name)})
}

