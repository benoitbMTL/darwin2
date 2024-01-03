package utils

import (
	"darwin2/config"
	"encoding/base64"
	"fmt"
)

func GenerateAPIToken() string {
	tokenData := fmt.Sprintf(`{"username":"%s","password":"%s","vdom":"%s"}`, config.CurrentConfig.USERNAMEAPI, config.CurrentConfig.PASSWORDAPI, config.CurrentConfig.VDOMAPI)
	return base64.StdEncoding.EncodeToString([]byte(tokenData))
}