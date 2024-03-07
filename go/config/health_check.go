package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"darwin2/utils"

	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////////////////////////////////////
// HEALTH CHECK                                                                  //
///////////////////////////////////////////////////////////////////////////////////

func HandleHealthCheck(c echo.Context) error {
	urls := []string{CurrentConfig.DVWAURL, CurrentConfig.BANKURL, CurrentConfig.JUICESHOPURL, CurrentConfig.PETSTOREURL, CurrentConfig.SPEEDTESTURL, "https://www.google.com"}

	// Define a custom HTTP client with a redirect policy that returns an error
	client := &http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Start HTML Table with CSS
	result := `<style>
			table	{width: 100%; border-collapse: collapse; font-size: 1em; font-family: 'Courier New', monospace;}
			td		{border: 1px solid #ddd; padding: 8px;}
			th		{border: 1px solid #ddd; padding: 8px; text-align: left; line-height: normal;}
			th h3	{margin: 0;  line-height: normal;}
			.down	{color: red; font-weight: bold;}
			.up		{color: green; font-weight: bold;}
		</style><table>
		<tr>
			<th><h3>URL</h3></th>
			<th><h3>Result</h3></th>
			<th><h3>Code</h3></th>
			<th><h3>Error</h3></th>
		</tr>`

	// Loop over the URLs
	for _, url := range urls {
		res, err := client.Get(url)
		if err != nil {
			shortErr := strings.TrimPrefix(err.Error(), fmt.Sprintf(`Get "%s": `, url))
			result += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td class="down">Down</td>
				<td>N/A</td>
				<td>%s</td>
			</tr>`, url, shortErr)
		} else {
			result += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td class="up">Up</td>
				<td>%d</td>
				<td>N/A</td>
			</tr>`, url, res.StatusCode)
		}
	}

	// Handle FortiWeb Management IP/FQDN separately to add scheme and port
	fwbManagementURL := "https://" + CurrentConfig.FWBMGTIP + ":" + CurrentConfig.FWBMGTPORT
	res, err := client.Get(fwbManagementURL)
	if err != nil {
		shortErr := strings.TrimPrefix(err.Error(), fmt.Sprintf(`Get "%s": `, fwbManagementURL))
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="down">Down</td>
			<td>N/A</td>
			<td>%s</td>
		</tr>`, fwbManagementURL, shortErr)
	} else {
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="up">Up</td>
			<td>%d</td>
			<td>N/A</td>
		</tr>`, fwbManagementURL, res.StatusCode)
	}

	// API Test
	resultStatus, resultCode, resultMessage, err := TestAPI(CurrentConfig.FWBMGTIP, CurrentConfig.FWBMGTPORT, utils.GenerateAPIToken())
	if err != nil {
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="down">%s</td>
			<td>%s</td>
			<td>%v</td>
		</tr>`, fwbManagementURL+"/api/v2.0/cmdb/system/global", resultStatus, "N/A", err)
	} else {
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="%s">%s</td>
			<td>%d</td>
			<td>%s</td>
		</tr>`, fwbManagementURL+"/api/v2.0/cmdb/system/global", resultStatus, resultStatus, resultCode, resultMessage)
	}

	// End HTML Table
	result += `</table>`

	return c.HTML(http.StatusOK, result)
}


func TestAPI(host, port, token string) (string, int, string, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/system/global", host, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Down", 0, "", err // Port is Down
	}
	defer resp.Body.Close()

	var response struct {
		Results struct {
			Hostname string `json:"hostname"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "Up", resp.StatusCode, "", fmt.Errorf("The API configuration is incorrect") // The API is not working correctly
	}

	if response.Results.Hostname == "" {
		return "Up", resp.StatusCode, "", fmt.Errorf("The API configuration is incorrect") // No Hostname
	}

	return "Up", resp.StatusCode, response.Results.Hostname, nil // All Good
}
