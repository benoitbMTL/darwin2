package handlers

import (
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"crypto/tls"
)

func HandleTrafficGenerator(c echo.Context) error {

	dvwaHost := config.CurrentConfig.DVWAHOST

	// Test if dvwaHost is responding on HTTP port 443
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Test if dvwaHost is responding on HTTPS port 443
	_, err := client.Get("https://" + dvwaHost)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, fmt.Sprintf("DVWA Host (%s) is not responding on HTTPS port 443: %s", dvwaHost, err.Error()))
	}

	// Check if nikto is installed
	_, err = exec.LookPath("nikto")
	if err != nil {
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
			"nikto",
			"-host", dvwaHost,
			"-ask", "no",
			"-followredirects",
			"-maxtime", "60s",
			"-nointeractive",
			"-no404",
			"-timeout", "2",
			"-useragent", "Nikto Traffic Gen\r\nX-Forwarded-For: "+randomIP,
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
