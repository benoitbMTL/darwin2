package config

import (
	"encoding/json"
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
		DVWAURL:        "https://dvwa.corp.fabriclab.ca",
		BANKURL:        "https://bank.corp.fabriclab.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop.corp.fabriclab.ca",
		PETSTOREURL:    "https://petstore3.corp.fabriclab.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest.corp.fabriclab.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "fortiweb.corp.fabriclab.ca",
		FWBMGTPORT:     "443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: "fortiweb2",
	}

	configsMap["AzureConfig"] = AppConfig{
		NAME:           "Azure Lab",
		DVWAURL:        "https://dvwa.canadaeast.cloudapp.azure.com/",
		BANKURL:        "https://dvwa.canadaeast.cloudapp.azure.com/bank.html",
		JUICESHOPURL:   "https://juiceshop.canadaeast.cloudapp.azure.com",
		PETSTOREURL:    "https://petstore3.canadaeast.cloudapp.azure.com/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest.canadaeast.cloudapp.azure.com",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "dvwa.canadaeast.cloudapp.azure.com",
		FWBMGTPORT:     "8443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: ""}

	configsMap["FabricLabConfig"] = AppConfig{
		NAME:           "Fabric Lab",
		DVWAURL:        "https://dvwa.corp.fabriclab.ca",
		BANKURL:        "https://bank.corp.fabriclab.ca/bank.html",
		JUICESHOPURL:   "https://juiceshop.corp.fabriclab.ca",
		PETSTOREURL:    "https://petstore3.corp.fabriclab.ca/api/v3/pet",
		SPEEDTESTURL:   "https://speedtest.corp.fabriclab.ca",
		USERNAMEAPI:    "userapi",
		PASSWORDAPI:    "fortinet123!",
		VDOMAPI:        "root",
		FWBMGTIP:       "fortiweb.corp.fabriclab.ca",
		FWBMGTPORT:     "443",
		MLPOLICY:       "DVWA_POLICY",
		USERAGENT:      "FortiWeb Demo Tool",
		FABRICLABSTORY: "fortiweb"}

	configsMap["FortiWebCloudConfig"] = AppConfig{
		NAME:           "FortiWeb Cloud",
		DVWAURL:        "https://dvwa.96859.fortiwebcloud.net",
		BANKURL:        "https://bank.96859.fortiwebcloud.net/bank.html",
		JUICESHOPURL:   "https://juiceshop.96859.fortiwebcloud.net",
		PETSTOREURL:    "https://petstore3.96859.fortiwebcloud.net/api/v3/pet",
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

	// Recherche de la configuration par défaut dans configsMap (optionnel si déjà fait dans Initialize)
	defaultConfig, exists := configsMap["Default"]
	if !exists {
		log.Println("Default configuration is missing.")
		return echo.NewHTTPError(http.StatusInternalServerError, "Default configuration is missing.")
	}

	// Réinitialise la configuration actuelle avec la configuration par défaut
	CurrentConfig = defaultConfig
	currentName = "Default" // Assurez-vous que currentName pointe aussi sur la configuration par défaut
	log.Println("Configuration reset to default.")

	// Optionnellement, vous pourriez vouloir mettre à jour la map configsMap avec la configuration réinitialisée
	configsMap[currentName] = CurrentConfig

	return c.JSON(http.StatusOK, defaultConfig)
}

func ExportConfig(c echo.Context) error {
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
	log.Printf("Configuration '%s' restored and set as current.", newConfig.NAME)
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

func ListConfigs(c echo.Context) error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	var configNames []string
	for name := range configsMap {
		configNames = append(configNames, name)
	}

	return c.JSON(http.StatusOK, configNames)
}
