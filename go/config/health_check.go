package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////////////////////////////////////
// HEALTH CHECK                                                                  //
///////////////////////////////////////////////////////////////////////////////////

func HandleHealthCheck(c echo.Context) error {
	urls := []string{CurrentConfig.DVWAURL, CurrentConfig.JUICESHOPURL, CurrentConfig.PETSTOREURL, CurrentConfig.BANKURL, CurrentConfig.SPEEDTESTURL, "https://www.google.com"}

	// Define a custom HTTP client with a redirect policy that returns an error
	client := &http.Client{
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

	// Handle FWB_MGT_IP separately because it's only an IP without a scheme
	ip := "https://" + CurrentConfig.FWBMGTIP
	res, err := client.Get(ip)
	if err != nil {
		shortErr := strings.TrimPrefix(err.Error(), fmt.Sprintf(`Get "%s": `, ip))
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="down">Down</td>
			<td>N/A</td>
			<td>%s</td>
		</tr>`, ip, shortErr)
	} else {
		result += fmt.Sprintf(`<tr>
			<td>%s</td>
			<td class="up">Up</td>
			<td>%d</td>
			<td>N/A</td>
		</tr>`, ip, res.StatusCode)
	}

	// End HTML Table
	result += `</table>`

	return c.HTML(http.StatusOK, result)
}
