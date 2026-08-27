package handlers

import (
	"crypto/tls"
	"darwin2/config"
	"darwin2/jobs"
	"darwin2/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
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
	ctx := c.Request().Context()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("The Web Server (%s) is not responding: %s", targetURL, err.Error()))
	}
	response.Body.Close()

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
		if err := ctx.Err(); err != nil {
			return err
		}
		randomIP := utils.GenerateRandomPublicIP()
		randomTuning := generateRandomTuning()

		// Construct the nikto command
		cmd := exec.CommandContext(ctx,
			"perl", "nikto/program/nikto.pl",
			"-host", targetURL,
			"-ask", "no",
			"-followredirects",
			"-maxtime", "60s",
			"-nointeractive",
			"-404code", "404",
			"-timeout", "2",
			"-useragent", "Nikto Traffic Generator\r\nX-Forwarded-For: "+randomIP,
			"-Tuning", randomTuning,
		)

		// Execute the nikto command
		_, commandErr := cmd.CombinedOutput()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if commandErr != nil {
			jobs.ReportError(ctx, fmt.Sprintf("scan %d", i+1), commandErr)
		}
		summary.WriteString(fmt.Sprintf("Executed Nikto Web Scan from IP %s with tuning %s\n", randomIP, randomTuning))
		jobs.ReportProgress(ctx, int64(i+1), loopCount, fmt.Sprintf("Completed scan %d of %d", i+1, loopCount))
	}

	// Return the summary of actions
	return c.String(200, summary.String())
}

func generateRandomTuning() string {
	tuningOptions := "0123456789abcde"
	return string(tuningOptions[rand.Intn(len(tuningOptions))])
}
