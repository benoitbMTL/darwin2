package config

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// AppConfig defines the structure for your application configuration
type AppConfig struct {
	NAME           string `json:"NAME"`
	DVWAURL        string `json:"DVWAURL"`
	BANKURL        string `json:"BANKURL"`
	JUICESHOPURL   string `json:"JUICESHOPURL"`
	PETSTOREURL    string `json:"PETSTOREURL"`
	SPEEDTESTURL   string `json:"SPEEDTESTURL"`
	USERNAMEAPI    string `json:"USERNAMEAPI"`
	PASSWORDAPI    string `json:"PASSWORDAPI"`
	VDOMAPI        string `json:"VDOMAPI"`
	FWBMGTIP       string `json:"FWBMGTIP"`
	FWBMGTPORT     string `json:"FWBMGTPORT"`
	MLPOLICY       string `json:"MLPOLICY"`
	USERAGENT      string `json:"USERAGENT"`
	FABRICLABSTORY string `json:"FABRICLABSTORY"`
}

var (
	CurrentConfig AppConfig                    // Holds the current application configuration
	DefaultConfig AppConfig                    // Holds the default application configuration
	configsMap    = make(map[string]AppConfig) // Stores configurations by name
	currentName   string                       // Holds the name of the current configuration
	configMutex   sync.RWMutex
)

// Initialize sets up the default values for the application configuration
func Initialize() {
	configMutex.Lock()
	defer configMutex.Unlock()

	defaultConfig := AppConfig{
		NAME:           "Default",
		DVWAURL:        "https://dvwa-xperts.labsec.ca",
		BANKURL:        "https://bank-xperts.labsec.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop-xperts.labsec.ca",
		PETSTOREURL:    "https://petstore3-xperts.labsec.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest-xperts.labsec.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "10.0.1.4",
		FWBMGTPORT:     "8443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: ""}

	configsMap["XPERTS 2026"] = AppConfig{
		NAME:           "XPERTS 2026",
		DVWAURL:        "https://dvwa-xperts.labsec.ca",
		BANKURL:        "https://bank-xperts.labsec.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop-xperts.labsec.ca",
		PETSTOREURL:    "https://petstore3-xperts.labsec.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest-xperts.labsec.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "10.0.1.4",
		FWBMGTPORT:     "8443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: ""}

	configsMap["Fabric Lab (fortiweb)"] = AppConfig{
		NAME:           "Fabric Lab (fortiweb)",
		DVWAURL:        "https://dvwa.corp.fabriclab.ca",
		BANKURL:        "https://bank.corp.fabriclab.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop.corp.fabriclab.ca",
		PETSTOREURL:    "https://petstore3.corp.fabriclab.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest.corp.fabriclab.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "10.163.7.21",
		FWBMGTPORT:     "443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: "fortiweb"}

	configsMap["Fabric Lab (fortiweb2)"] = AppConfig{
		NAME:           "Fabric Lab (fortiweb2)",
		DVWAURL:        "https://dvwa.corp.fabriclab.ca",
		BANKURL:        "https://bank.corp.fabriclab.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop.corp.fabriclab.ca",
		PETSTOREURL:    "https://petstore3.corp.fabriclab.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest.corp.fabriclab.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "10.163.7.40",
		FWBMGTPORT:     "443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: "fortiweb2"}

	configsMap["FortiAppSec Cloud"] = AppConfig{
		NAME:           "FortiWeb Cloud",
		DVWAURL:        "https://dvwa.labsec.ca",
		BANKURL:        "https://bank.labsec.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop.labsec.ca",
		PETSTOREURL:    "https://petstore3.labsec.ca/api/v3/pet",
		SPEEDTESTURL:   "",
		USERNAMEAPI:    "",
		PASSWORDAPI:    "",
		VDOMAPI:        "",
		FWBMGTIP:       "",
		FWBMGTPORT:     "",
		MLPOLICY:       "",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: ""}

	// Add your defaultConfig to the configsMap
	configsMap[defaultConfig.NAME] = defaultConfig

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

// SaveConfig handles the POST request to update the configuration
func SaveConfig(c echo.Context) error {
	var newConfig AppConfig
	if err := c.Bind(&newConfig); err != nil {
		log.Printf("Error binding new configuration: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	if newConfig.NAME == "Default" {
		log.Println("Error: Cannot overwrite the Default configuration.")
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot overwrite the Default configuration. Please use a different name.")
	}

	// Sauvegarder la nouvelle configuration sans affecter DefaultConfig
	configsMap[newConfig.NAME] = newConfig

	// Optionnel: définir cette nouvelle configuration comme la configuration actuelle
	currentName = newConfig.NAME
	CurrentConfig = newConfig

	log.Printf("Configuration '%s' updated and set as current.", newConfig.NAME)
	return c.JSON(http.StatusOK, newConfig)
}

func ResetConfig(c echo.Context) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Search for the default configuration in configsMap (optional if already done in Initialize)
	defaultConfig, exists := configsMap["Default"]
	if !exists {
		log.Println("Default configuration is missing.")
		return echo.NewHTTPError(http.StatusInternalServerError, "Default configuration is missing.")
	}

	// Resets the current configuration to the default configuration
	CurrentConfig = defaultConfig
	currentName = "Default" // Assurez-vous que currentName pointe aussi sur la configuration par défaut
	log.Println("Configuration reset to default.")

	// Optionally, you might want to update the configsMap with the reset configuration
	configsMap[currentName] = CurrentConfig
	return c.JSON(http.StatusOK, defaultConfig)
}

func ImportConfig(c echo.Context) error {
	var newConfig AppConfig
	if err := c.Bind(&newConfig); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	if newConfig.NAME == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Configuration name is required.")
	}

	configsMap[newConfig.NAME] = newConfig
	currentName = newConfig.NAME
	log.Printf("Configuration '%s' imported and set as current.", newConfig.NAME)
	return c.JSON(http.StatusOK, newConfig)
}

func CloneConfigLocal(c echo.Context) error {
	var request struct {
		SourceName string `json:"sourceName"`
		NewName    string `json:"newName"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	// Ensure the source configuration exists
	sourceConfig, exists := configsMap[request.SourceName]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Source configuration '%s' not found", request.SourceName))
	}

	// Ensure not to overwrite an existing configuration with the new name
	if _, exists := configsMap[request.NewName]; exists {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Configuration '%s' already exists", request.NewName))
	}

	// Clone the source configuration under the new name
	clonedConfig := sourceConfig
	clonedConfig.NAME = request.NewName        // Update the NAME field to reflect the new configuration name
	configsMap[request.NewName] = clonedConfig // Add the cloned configuration to the map under the new name
	currentName = request.NewName              // Set the cloned configuration as the current configuration

	log.Printf("Configuration '%s' cloned to '%s' and set as current.", request.SourceName, request.NewName)

	return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' cloned to '%s' and set as current", request.SourceName, request.NewName)})
}

func ApplyConfigLocal(c echo.Context) error {
	var request struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	newConfig, exists := configsMap[request.Name]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Configuration '%s' not found", request.Name))
	}

	// Update the current configuration name
	currentName = request.Name
	// Update the CurrentConfig global variable to the newly selected configuration
	CurrentConfig = newConfig

	// Log configuration
	log.Printf("Configuration '%s' applied successfully. Details: NAME: %s, DVWAURL: %s, BANKURL: %s",
		request.Name,
		newConfig.NAME,
		newConfig.DVWAURL,
		newConfig.BANKURL,
	)

	log.Printf("Configuration '%s' applied successfully.", request.Name)

	return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' applied successfully", request.Name)})
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

func ListConfigs(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var configNames []string
	for name := range configsMap {
		configNames = append(configNames, name)
	}

	return c.JSON(http.StatusOK, configNames)
}

func RenameConfig(c echo.Context) error {
	var request struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	if _, exists := configsMap[request.NewName]; exists {
		return echo.NewHTTPError(http.StatusConflict, `{"message": "Configuration name '`+request.NewName+`' already exists."}`)
	}

	config, exists := configsMap[request.OldName]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, `{"message": "Configuration '`+request.OldName+`' not found."}`)
	}

	delete(configsMap, request.OldName)  // Remove the old config
	config.NAME = request.NewName        // Update the config name
	configsMap[request.NewName] = config // Add the renamed config

	if currentName == request.OldName {
		currentName = request.NewName // Update current config if it was renamed
	}

	return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Configuration '%s' renamed to '%s'", request.OldName, request.NewName)})
}
