package handlers

import (
	"context"
	"darwin2/config"
	"darwin2/utils"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 1: Delete Policy
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeletePolicy(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	policyName := config.PolicyName

	result := ApiResult{
		TaskID:      "deletePolicy",
		Description: "Delete Policy",
	}

	resultBody, err := deletePolicy(c.Request().Context(), host, port, token, policyName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting Policy: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Policy"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted Policy"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deletePolicy(ctx context.Context, host, port, token, policyName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/policy?mkey=%s", host, port, url.QueryEscape(policyName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 2: Delete Inline Protection Profile
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteInlineProtection(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	profileName := config.CloneInlineProtectionProfileName

	result := ApiResult{
		TaskID:      "deleteInlineProtection",
		Description: "Delete Inline Protection Profile",
	}

	resultBody, err := deleteInlineProtection(c.Request().Context(), host, port, token, profileName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting Inline Protection Profile: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Inline Protection Profile"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted Inline Protection Profile"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteInlineProtection(ctx context.Context, host, port, token, profileName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s", host, port, url.QueryEscape(profileName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 3: Delete X-Forwarded-For Rule
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteXForwardedForRule(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	xff := config.XForwardedForName

	result := ApiResult{
		TaskID:      "deleteXForwardedForRule",
		Description: "Delete X-Forwarded-For Rule",
	}

	resultBody, err := deleteXForwardedForRule(c.Request().Context(), host, port, token, xff)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting X-Forwarded-For Rule: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete X-Forwarded-For Rule"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted X-Forwarded-For Rule"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteXForwardedForRule(ctx context.Context, host, port, token, ruleName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/x-forwarded-for?mkey=%s", host, port, url.QueryEscape(ruleName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 4: Delete Signature Protection
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteSignatureProtection(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	SignatureProtectionName := config.CloneSignatureProtectionName

	result := ApiResult{
		TaskID:      "deleteSignatureProtection",
		Description: "Delete Signature Protection",
	}

	resultBody, err := deleteSignatureProtection(c.Request().Context(), host, port, token, SignatureProtectionName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting Signature Protection: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Signature Protection"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted Signature Protection"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteSignatureProtection(ctx context.Context, host, port, token, signatureName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/waf/signature?mkey=%s", host, port, url.QueryEscape(signatureName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 5: DeleteVirtualServer
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteVirtualServer(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	virtualServerName := config.VirtualServerName

	result := ApiResult{
		TaskID:      "deleteVirtualServer",
		Description: "Delete Virtual Server",
	}

	resultBody, err := deleteVirtualServer(c.Request().Context(), host, port, token, virtualServerName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting Virtual Server: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Virtual Server"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted Virtual Server"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteVirtualServer(ctx context.Context, host, port, token, virtualServerName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/vserver?mkey=%s", host, port, url.QueryEscape(virtualServerName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 6 : Delete Server Pool
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteServerPool(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	poolName := config.PoolName

	result := ApiResult{
		TaskID:      "deleteServerPool",
		Description: "Delete Server Pool",
	}

	resultBody, err := deleteServerPool(c.Request().Context(), host, port, token, poolName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting Server Pool: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Server Pool"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted Server Pool"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteServerPool(ctx context.Context, host, port, token, poolName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/server-policy/server-pool?mkey=%s", host, port, url.QueryEscape(poolName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Step 7: Delete Virtual IP
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HandleDeleteVirtualIP(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	port := config.CurrentConfig.FWBMGTPORT
	token := utils.GenerateAPIToken()
	vipName := config.VipName

	result := ApiResult{
		TaskID:      "deleteVirtualIP",
		Description: "Delete Virtual IP",
	}

	resultBody, err := deleteVirtualIP(c.Request().Context(), host, port, token, vipName)
	if err != nil {
		result.Status = "failure"
		result.Message = fmt.Sprintf("Error deleting virtual IP: %v", err)
	} else if !utils.CheckOperationStatus(resultBody) {
		result.Status = "failure"
		result.Message = "Failed to delete Virtual IP"
	} else {
		result.Status = "success"
		result.Message = "Successfully deleted virtual IP"
	}

	fmt.Printf("Result: %+v\n", result)
	return c.JSON(http.StatusOK, result)
}

func deleteVirtualIP(ctx context.Context, host, port, token, vipName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s/api/v2.0/cmdb/system/vip?mkey=%s", host, port, url.QueryEscape(vipName))

	return utils.SendRequest(ctx, "DELETE", url, token, nil)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
