package handlers

import (
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"crypto/tls"
)

func HandleTrafficGenerator(c echo.Context) error {
    var requestPayload struct {
        Target string `json:"target"`
    }
    if err := c.Bind(&requestPayload); err != nil {
        return c.String(http.StatusBadRequest, "Invalid request payload")
    }
    fmt.Println("Received target choice:", requestPayload.Target) // Debug log

    // Select the target URL based on the user's choice
    var targetURL string
    switch requestPayload.Target {
    case "DVWA":
        targetURL = config.CurrentConfig.DVWAURL
    case "Bank":
        targetURL = config.CurrentConfig.BANKURL
    case "JuiceShop":
        targetURL = config.CurrentConfig.JUICESHOPURL
    case "Petstore":
        targetURL = config.CurrentConfig.PETSTOREURL
    case "Speedtest":
        targetURL = config.CurrentConfig.SPEEDTESTURL
    default:
        return c.String(http.StatusBadRequest, "Invalid target choice")
    }
    fmt.Println("Using target URL:", targetURL) // Debug log

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	_, err := client.Get(targetURL)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("The Web Server (%s) is not responding: %s", targetURL, err.Error()))
	}

	_, err = exec.LookPath("perl")
	if err != nil {
		return c.String(200, "Perl is not installed on your system")
	}

	niktoScriptPath := "nikto/program/nikto.pl"

	// Check if the Nikto script exists
	if _, err := os.Stat(niktoScriptPath); os.IsNotExist(err) {
		return c.String(200, "Nikto is not installed on your system")
	}

	const loopCount = 100 // Number of iterations for each loop

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Traffic Generator executed %d rounds of attacks:\n", loopCount))

	for i := 0; i < loopCount; i++ {
		randomIP := utils.GenerateRandomPublicIP()
		randomTuning := generateRandomTuning()

		// Construct the nikto command
		cmd := exec.Command(
			"perl", "nikto/program/nikto.pl",
			"-host", targetURL,
			"-ask", "no",
			"-followredirects",
			"-maxtime", "60s",
			"-nointeractive",
			"-no404",
			"-timeout", "2",
			"-useragent", "Nikto Traffic Generator\r\nX-Forwarded-For: "+randomIP,
			"-T", randomTuning,
		)

		// Execute the nikto command
		cmd.CombinedOutput()
		summary.WriteString(fmt.Sprintf("Executed Nikto Web Scan from IP %s with tuning %s\n", randomIP, randomTuning))
	}

	// Return the summary of actions
	return c.String(200, summary.String())
}

func generateRandomTuning() string {
	tuningOptions := "0123456789abcde"
	return string(tuningOptions[rand.Intn(len(tuningOptions))])
}
