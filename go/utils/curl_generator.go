package utils

import (
	"fmt"
	"net/http"
)


func GenerateCurlCommand(req *http.Request, body []byte) string {
	curl := fmt.Sprintf("curl -k -X %s '%s'", req.Method, req.URL)
	for key, values := range req.Header {
		for _, value := range values {
			curl += fmt.Sprintf(" -H '%s: %s'", key, value)
		}
	}
	if len(body) > 0 {
		curl += fmt.Sprintf(" -d '%s'", string(body))
	}
	return curl
}