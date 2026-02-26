package handlers

import (
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type ApiResult struct {
	TaskID      string
	Description string
	Status      string
	Message     string
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 1: Create new Virtual IP
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewVirtualIP(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()

	vipData := config.VirtualIPData{
		Name:      config.VipName,
		Vip:       config.VipIp,
		Interface: config.Interface,
	}

	result := ApiResult{
		TaskID:      "createNewVirtualIP",
		Description: "Create new Virtual IP",
	}

	resultBody, err := createNewVirtualIP(fwbmgtip, fwbmgtport, token, vipData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error creating virtual IP: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to create virtual IP"
	} else {
		result.Status = "success"
		result.Message = "Successfully created virtual IP"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewVirtualIP(host, port, token string, data config.VirtualIPData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/system/vip", host, port)

	return utils.SendRequest("POST", url, token, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 2: Create new Server Pool
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewServerPool(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()

	poolData := config.ServerPoolData{
		Name:          config.PoolName,
		ServerBalance: config.ServerBalance,
		Health:        config.HealthCheck,
	}

	result := ApiResult{
		TaskID:      "createNewServerPool",
		Description: "Create new Server Pool",
	}

	resultBody, err := createNewServerPool(fwbmgtip, fwbmgtport, token, poolData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error creating Server Pool: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to create Server Pool"
	} else {
		result.Status = "success"
		result.Message = "Successfully created Server Pool"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewServerPool(host, port, token string, data config.ServerPoolData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/server-pool", host, port)

	return utils.SendRequest("POST", url, token, data)

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 3: Create new Member Pool
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewMemberPool(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	poolName := config.PoolName

	poolMembers := make([]config.MemberPoolData, len(config.PoolMemberIPs))
	for i, ip := range config.PoolMemberIPs {
		poolMembers[i] = config.MemberPoolData{IP: ip, SSL: config.PoolMemberSSL, Port: config.PoolMemberPort}
	}

	result := ApiResult{
		TaskID:      "createNewMemberPool",
		Description: "Create new Member Pool",
	}

	for _, member := range poolMembers {
		resultBody, err := createNewMemberPool(fwbmgtip, fwbmgtport, token, poolName, member)
		if err != nil {
			result.Status = "failure"
			result.Message = fmt.Sprintf("Error creating Member Pool: %v", err)
		} else if !utils.CheckOperationStatus(resultBody) {
			result.Status = "failure"
			result.Message = "Failed to create Member Pool"
		} else {
			result.Status = "success"
			result.Message = "Successfully created Member Pool"
		}
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewMemberPool(host, port, token, poolName string, data config.MemberPoolData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/server-pool/pserver-list?mkey=%s", host, port, url.QueryEscape(poolName))

	return utils.SendRequest("POST", url, token, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 4: Create new Virtual Server
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewVirtualServer(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	vsData := config.VirtualServerData{
		Name: config.VirtualServerName,
	}

	result := ApiResult{
		TaskID:      "createNewVirtualServer",
		Description: "Create new Virtual Server",
	}

	resultBody, err := createNewVirtualServer(fwbmgtip, fwbmgtport, token, vsData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error creating Virtual Server: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to create Virtual Server"
	} else {
		result.Status = "success"
		result.Message = "Successfully created Virtual Server"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewVirtualServer(host, port, token string, data config.VirtualServerData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/vserver", host, port)

	return utils.SendRequest("POST", url, token, data)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 5: Assign Virtual IP to Virtual Server
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleAssignVIPToVirtualServer(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	virtualServerName := config.VirtualServerName

	assignVIPData := config.AssignVIPData{
		Interface: config.Interface,
		Status:    config.VipStatus,
		VipName:   config.VipName,
	}

	result := ApiResult{
		TaskID:      "assignVIPToVirtualServer",
		Description: "Assign Virtual IP to Virtual Server",
	}

	resultBody, err := assignVIPToVirtualServer(fwbmgtip, fwbmgtport, token, virtualServerName, assignVIPData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error assigning VIP to Virtual Server: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to assign VIP to Virtual Server"
	} else {
		result.Status = "success"
		result.Message = "Successfully assigned VIP to Virtual Server"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func assignVIPToVirtualServer(host, port, token, virtualServerName string, data config.AssignVIPData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/vserver/vip-list?mkey=%s", host, port, url.QueryEscape(virtualServerName))

	return utils.SendRequest("POST", url, token, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 6: Clone Signature Protection
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCloneSignatureProtection(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	originalKey := config.OriginalSignatureProtectionName
	cloneKey := config.CloneSignatureProtectionName

	result := ApiResult{
		TaskID:      "cloneSignatureProtection",
		Description: "Clone Signature Protection",
	}

	resultBody, err := cloneSignatureProtection(fwbmgtip, fwbmgtport, token, originalKey, cloneKey)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error cloning Signature Protection: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to clone Signature Protection"
	} else {
		result.Status = "success"
		result.Message = "Successfully cloned Signature Protection"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func cloneSignatureProtection(host, port, token, originalKey, cloneKey string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/signature?clone_mkey=%s&mkey=%s", host, port, url.QueryEscape(cloneKey), url.QueryEscape(originalKey))

	return utils.SendRequest("POST", url, token, nil)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 7: Clone Inline Protection
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCloneInlineProtection(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	originalKey := config.OriginalInlineProtectionProfileName
	cloneKey := config.CloneInlineProtectionProfileName

	result := ApiResult{
		TaskID:      "cloneInlineProtection",
		Description: "Clone Inline Protection",
	}

	resultBody, err := cloneInlineProtection(fwbmgtip, fwbmgtport, token, originalKey, cloneKey)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error cloning Inline Protection: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to clone Inline Protection"
	} else {
		result.Status = "success"
		result.Message = "Successfully cloned Inline Protection"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func cloneInlineProtection(host, port, token, originalKey, cloneKey string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s&clone_mkey=%s", host, port, url.QueryEscape(originalKey), url.QueryEscape(cloneKey))

	return utils.SendRequest("POST", url, token, nil)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 8: Create new X-Forwarded-For Rule
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewXForwardedForRule(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	xffData := config.XForwardedForData{
		Name:                 config.XForwardedForName,
		XForwardedForSupport: config.XForwardedForSupport,
	}

	result := ApiResult{
		TaskID:      "createNewXForwardedForRule",
		Description: "Create new X-Forwarded-For Rule",
	}

	resultBody, err := createNewXForwardedForRule(fwbmgtip, fwbmgtport, token, xffData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error creating new X-Forwarded-For Rule: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to create new X-Forwarded-For Rule"
	} else {
		result.Status = "success"
		result.Message = "Successfully created new X-Forwarded-For Rule"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewXForwardedForRule(host, port, token string, data config.XForwardedForData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/x-forwarded-for", host, port)

	return utils.SendRequest("POST", url, token, data)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 9: Configure Protection Profile
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleConfigureProtectionProfile(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	cloneKey := config.CloneInlineProtectionProfileName

	protectionProfileData := config.ProtectionProfileData{
		SignatureRule:     config.CloneSignatureProtectionName,
		XForwardedForRule: config.XForwardedForName,
	}

	result := ApiResult{
		TaskID:      "configureProtectionProfile",
		Description: "Configure Protection Profile",
	}

	resultBody, err := configureProtectionProfile(fwbmgtip, fwbmgtport, token, cloneKey, protectionProfileData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error configuring Protection Profile: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to configure Protection Profile"
	} else {
		result.Status = "success"
		result.Message = "Successfully configured Protection Profile"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func configureProtectionProfile(host, port, token, mkey string, data config.ProtectionProfileData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s", host, port, url.QueryEscape(mkey))

	return utils.SendRequest("PUT", url, token, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 10: Create new Policy
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleCreateNewPolicy(c echo.Context) error {
	fwbmgtip := config.CurrentConfig.FWBMGTIP
	fwbmgtport := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()

	policyData := config.PolicyData{
		Name:                    config.PolicyName,
		DeploymentMode:          config.PolicyDeploymentMode,
		Protocol:                config.PolicyProtocol,
		Ssl:                     config.PolicySSL,
		ImplicitSsl:             config.PolicyImplicitSSL,
		Vserver:                 config.PolicyVirtualServer,
		Service:                 config.PolicyService,
		InlineProtectionProfile: config.PolicyInlineProtectionProfile,
		ServerPool:              config.PolicyServerPool,
		TrafficLog:              config.PolicyTrafficLog,
		HttpsService:            config.PolicyHTTPSService,
	}

	result := ApiResult{
		TaskID:      "configureProtectionProfile",
		Description: "Configure Protection Profile",
	}

	resultBody, err := createNewPolicy(fwbmgtip, fwbmgtport, token, policyData)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error creating new Policy: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to create new Policy"
	} else {
		result.Status = "success"
		result.Message = "Successfully created new Policy"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func createNewPolicy(host, port, token string, data config.PolicyData) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/policy", host, port)

	return utils.SendRequest("POST", url, token, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
