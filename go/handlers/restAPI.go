package handlers

import (
	"bytes"
	"crypto/tls"
	"darwin2/config"
	"darwin2/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
)

// Virtual IP

func createNewVirtualIP(host, token string, data config.VirtualIPData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/system/vip", host)

	return sendRequest("POST", url, token, data)
}

func deleteVirtualIP(host, token, vipName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/system/vip?mkey=%s", host, url.QueryEscape(vipName))

	return sendRequest("DELETE", url, token, nil)
}

// Server Pool

func createNewServerPool(host, token string, data config.ServerPoolData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/server-pool", host)

	return sendRequest("POST", url, token, data)

}

func deleteServerPool(host, token, poolName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/server-pool?mkey=%s", host, url.QueryEscape(poolName))

	return sendRequest("DELETE", url, token, nil)
}

// Member Pool

func createNewMemberPool(host, token, poolName string, data config.MemberPoolData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/server-pool/pserver-list?mkey=%s", host, url.QueryEscape(poolName))

	return sendRequest("POST", url, token, data)
}

// Virtual Server

func createNewVirtualServer(host, token string, data config.VirtualServerData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/vserver", host)

	return sendRequest("POST", url, token, data)
}

func deleteVirtualServer(host, token, virtualServerName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/vserver?mkey=%s", host, url.QueryEscape(virtualServerName))

	return sendRequest("DELETE", url, token, nil)
}

// Assign VIP to Virtual Server

func assignVIPToVirtualServer(host, token, virtualServerName string, data config.AssignVIPData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/vserver/vip-list?mkey=%s", host, url.QueryEscape(virtualServerName))

	return sendRequest("POST", url, token, data)
}

// Signature Protection

func cloneSignatureProtection(host, token, originalKey, cloneKey string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/signature?mkey=%s&clone_mkey=%s", host, url.QueryEscape(originalKey), url.QueryEscape(cloneKey))

	return sendRequest("POST", url, token, nil)
}

func deleteSignatureProtection(host, token, signatureName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/signature?mkey=%s", host, url.QueryEscape(signatureName))

	return sendRequest("DELETE", url, token, nil)
}

// Inline Protection Profile

func cloneInlineProtection(host, token, originalKey, cloneKey string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s&clone_mkey=%s", host, url.QueryEscape(originalKey), url.QueryEscape(cloneKey))

	return sendRequest("POST", url, token, nil)
}

func deleteInlineProtection(host, token, profileName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s", host, url.QueryEscape(profileName))

	return sendRequest("DELETE", url, token, nil)
}

// X-Forwarded-For Rule

func createNewXForwardedForRule(host, token string, data config.XForwardedForData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/x-forwarded-for", host)

	return sendRequest("POST", url, token, data)
}

func deleteXForwardedForRule(host, token, ruleName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/x-forwarded-for?mkey=%s", host, url.QueryEscape(ruleName))

	return sendRequest("DELETE", url, token, nil)
}

// Protection Profile

func configureProtectionProfile(host, token, mkey string, data config.ProtectionProfileData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/waf/web-protection-profile.inline-protection?mkey=%s", host, url.QueryEscape(mkey))

	return sendRequest("PUT", url, token, data)
}

// Policy

func createNewPolicy(host, token string, data config.PolicyData) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/policy", host)

	return sendRequest("POST", url, token, data)
}

func deletePolicy(host, token, policyName string) ([]byte, error) {
	url := fmt.Sprintf("https://%s/api/v2.0/cmdb/server-policy/policy?mkey=%s", host, url.QueryEscape(policyName))

	return sendRequest("DELETE", url, token, nil)
}

///////////////////////////////////////////////////////////////////////////////
// SEND REQUEST                                                              //
///////////////////////////////////////////////////////////////////////////////

func sendRequest(method, url, token string, data interface{}) ([]byte, error) {
	// fmt.Println("Sending HTTP Request...")
	var req *http.Request
	var err error

	reqData := config.Request{
		Data: data,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("Error marshaling JSON data: %v\n", err)
		return nil, err
	}

	// Convert jsonData to string for comparison
	jsonDataStr := string(jsonData)

	if jsonDataStr != "" && jsonDataStr != `{"data":null}` {
		// Create a new request with JSON data
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating HTTP request: %v\n", err)
			return nil, err
		}
	} else {
		// Create a new request without data
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Printf("Error creating HTTP request: %v\n", err)
			return nil, err
		}
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	// Create a custom HTTP client with SSL/TLS certificate verification disabled
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// fmt.Println("Sending HTTP Request...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	// fmt.Println("Waiting for response...")
	time.Sleep(time.Duration(1000) * time.Millisecond)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	// fmt.Println("Received HTTP Response")
	return body, nil
}

// Checks JSON result to see if the operation was successful or not.
// It returns true if the operation was successful, and false otherwise.
func checkOperationStatus(result []byte) bool {
	var res map[string]interface{}
	json.Unmarshal(result, &res)

	// Check if the result map is empty
	if len(res) == 0 {
		return false
	}

	// Check if the "data" field is null
	if data, ok := res["data"]; ok && data == nil {
		return false
	}

	if _, ok := res["results"].(map[string]interface{})["errcode"]; ok {
		// The result contains an error code, so the operation failed
		return false
	}
	// The operation succeeded
	return true
}

///////////////////////////////////////////////////////////////////////////////////
// CREATE POLICY                                                                 //
///////////////////////////////////////////////////////////////////////////////////

func HandleCreatePolicy(c echo.Context) error {
	host := config.CurrentConfig.FWBMGTIP
	token := utils.GenerateAPIToken()

	vipData := config.VirtualIPData{
		Name:      config.VipName,
		Vip:       config.VipIp,
		Interface: config.Interface,
	}

	poolData := config.ServerPoolData{
		Name:          config.PoolName,
		ServerBalance: config.ServerBalance,
		Health:        config.HealthCheck,
	}

	poolMembers := make([]config.MemberPoolData, len(config.PoolMemberIPs))
	for i, ip := range config.PoolMemberIPs {
		poolMembers[i] = config.MemberPoolData{IP: ip, SSL: config.PoolMemberSSL, Port: config.PoolMemberPort}
	}

	vsData := config.VirtualServerData{
		Name: config.VirtualServerName,
	}

	assignVIPData := config.AssignVIPData{
		Interface: config.Interface,
		Status:    config.VipStatus,
		VipName:   config.VipName,
	}

	xffData := config.XForwardedForData{
		Name:                 config.XForwardedForName,
		XForwardedForSupport: config.XForwardedForSupport,
	}

	protectionProfileData := config.ProtectionProfileData{
		SignatureRule:     config.CloneSignatureProtectionName,
		XForwardedForRule: config.XForwardedForName,
	}

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
		//Certificate:          PolicyCertificate,
	}

	// Initialize a slice to store the statuses
	statuses := []map[string]string{}

	// Debug Statements
	fmt.Println("Debug: host =", host)
	fmt.Println("Debug: token =", token)
	fmt.Println("Debug: vipData =", vipData)
	fmt.Println("Debug: poolData =", poolData)
	fmt.Println("Debug: poolMembers =", poolMembers)
	fmt.Println("Debug: vsData =", vsData)
	fmt.Println("Debug: assignVIPData =", assignVIPData)
	fmt.Println("Debug: xffData =", xffData)
	fmt.Println("Debug: protectionProfileData =", protectionProfileData)
	fmt.Println("Debug: policyData =", policyData)

	// Step 1: Create new Virtual IP
	result, err := createNewVirtualIP(host, token, vipData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualIP",
			"status":      "failure",
			"description": "Create new Virtual IP",
			"message":     fmt.Sprintf("Error creating virtual IP: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualIP",
			"status":      "failure",
			"description": "Create new Virtual IP",
			"message":     "Failed to create virtual IP",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualIP",
			"status":      "success",
			"description": "Create new Virtual IP",
			"message":     "Successfully created virtual IP",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 2: Create new Server Pool
	result, err = createNewServerPool(host, token, poolData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewServerPool",
			"status":      "failure",
			"description": "Create new Server Pool",
			"message":     fmt.Sprintf("Error creating Server Pool: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewServerPool",
			"status":      "failure",
			"description": "Create new Server Pool",
			"message":     "Failed to create Server Pool",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewServerPool",
			"status":      "success",
			"description": "Create new Server Pool",
			"message":     "Successfully created Server Pool",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 3: Create new Member Pool
	for _, member := range poolMembers {
		result, err := createNewMemberPool(host, token, poolData.Name, member)
		if err != nil {
			statuses = append(statuses, map[string]string{
				"taskId":      "createNewMemberPool",
				"status":      "failure",
				"description": "Create new Member Pool",
				"message":     fmt.Sprintf("Error creating Member Pool: %v", err),
			})
		} else if !checkOperationStatus(result) {
			statuses = append(statuses, map[string]string{
				"taskId":      "createNewMemberPool",
				"status":      "failure",
				"description": "Create new Member Pool",
				"message":     "Failed to create Member Pool",
			})
		} else {
			statuses = append(statuses, map[string]string{
				"taskId":      "createNewMemberPool",
				"status":      "success",
				"description": "Create new Member Pool",
				"message":     "Successfully created Member Pool",
			})
		}
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 4: Create new Virtual Server
	result, err = createNewVirtualServer(host, token, vsData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualServer",
			"status":      "failure",
			"description": "Create new Virtual Server",
			"message":     fmt.Sprintf("Error creating Virtual Server: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualServer",
			"status":      "failure",
			"description": "Create new Virtual Server",
			"message":     "Failed to create Virtual Server",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewVirtualServer",
			"status":      "success",
			"description": "Create new Virtual Server",
			"message":     "Successfully created Virtual Server",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 5: Assign Virtual IP to Virtual Server
	result, err = assignVIPToVirtualServer(host, token, config.VirtualServerName, assignVIPData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "assignVIPToVirtualServer",
			"status":      "failure",
			"description": "Assign Virtual IP to Virtual Server",
			"message":     fmt.Sprintf("Error assigning VIP to Virtual Server: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "assignVIPToVirtualServer",
			"status":      "failure",
			"description": "Assign Virtual IP to Virtual Server",
			"message":     "Failed to assign VIP to Virtual Server",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "assignVIPToVirtualServer",
			"status":      "success",
			"description": "Assign Virtual IP to Virtual Server",
			"message":     "Successfully assigned VIP to Virtual Server",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 6: Clone Signature Protection
	result, err = cloneSignatureProtection(host, token, config.OriginalSignatureProtectionName, config.CloneSignatureProtectionName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneSignatureProtection",
			"status":      "failure",
			"description": "Clone Signature Protection",
			"message":     fmt.Sprintf("Error cloning Signature Protection: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneSignatureProtection",
			"status":      "failure",
			"description": "Clone Signature Protection",
			"message":     "Failed to clone Signature Protection",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneSignatureProtection",
			"status":      "success",
			"description": "Clone Signature Protection",
			"message":     "Successfully cloned Signature Protection",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 7: Clone Inline Protection
	result, err = cloneInlineProtection(host, token, config.OriginalInlineProtectionProfileName, config.CloneInlineProtectionProfileName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneInlineProtection",
			"status":      "failure",
			"description": "Clone Inline Protection",
			"message":     fmt.Sprintf("Error cloning Inline Protection: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneInlineProtection",
			"status":      "failure",
			"description": "Clone Inline Protection",
			"message":     "Failed to clone Inline Protection",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "cloneInlineProtection",
			"status":      "success",
			"description": "Clone Inline Protection",
			"message":     "Successfully cloned Inline Protection",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 8: Create new X-Forwarded-For Rule
	result, err = createNewXForwardedForRule(host, token, xffData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewXForwardedForRule",
			"status":      "failure",
			"description": "Create new X-Forwarded-For Rule",
			"message":     fmt.Sprintf("Error creating new X-Forwarded-For Rule: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewXForwardedForRule",
			"status":      "failure",
			"description": "Create new X-Forwarded-For Rule",
			"message":     "Failed to create new X-Forwarded-For Rule",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewXForwardedForRule",
			"status":      "success",
			"description": "Create new X-Forwarded-For Rule",
			"message":     "Successfully created new X-Forwarded-For Rule",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 9: Configure Protection Profile
	result, err = configureProtectionProfile(host, token, config.CloneInlineProtectionProfileName, protectionProfileData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "configureProtectionProfile",
			"status":      "failure",
			"description": "Configure Protection Profile",
			"message":     fmt.Sprintf("Error configuring Protection Profile: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "configureProtectionProfile",
			"status":      "failure",
			"description": "Configure Protection Profile",
			"message":     "Failed to configure Protection Profile",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "configureProtectionProfile",
			"status":      "success",
			"description": "Configure Protection Profile",
			"message":     "Successfully configured Protection Profile",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 10: Create new Policy
	result, err = createNewPolicy(host, token, policyData)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewPolicy",
			"status":      "failure",
			"description": "Create new Policy",
			"message":     fmt.Sprintf("Error creating new Policy: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewPolicy",
			"status":      "failure",
			"description": "Create new Policy",
			"message":     "Failed to create new Policy",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "createNewPolicy",
			"status":      "success",
			"description": "Create new Policy",
			"message":     "Successfully created new Policy",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Return a JSON response with the statuses of all steps
	return c.JSON(http.StatusOK, statuses)
}

///////////////////////////////////////////////////////////////////////////////////
// DELETE POLICY                                                                 //
///////////////////////////////////////////////////////////////////////////////////

func HandleDeletePolicy(c echo.Context) error {

			fmt.Printf("START HandleDeletePolicy")

	host := config.CurrentConfig.FWBMGTIP
	token := utils.GenerateAPIToken()

	// Initialize a slice to store the statuses
	statuses := []map[string]string{}

	// Step 1: Delete Policy
	result, err := deletePolicy(host, token, config.PolicyName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "deletePolicy",
			"status":      "failure",
			"description": "Delete Policy",
			"message":     fmt.Sprintf("Error deleting Policy: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "deletePolicy",
			"status":      "failure",
			"description": "Delete Policy",
			"message":     "Failed to delete Policy",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deletePolicy",
			"status":      "success",
			"description": "Delete Policy",
			"message":     "Successfully deleted Policy",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 2: Delete Inline Protection Profile
	result, err = deleteInlineProtection(host, token, config.CloneInlineProtectionProfileName)
	if err != nil {
		// ...
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteInlineProtection",
			"status":      "failure",
			"description": "Delete Inline Protection Profile",
			"message":     fmt.Sprintf("Error deleting Inline Protection Profile: %v", err),
		})
	} else if !checkOperationStatus(result) {
		// ...
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteInlineProtection",
			"status":      "failure",
			"description": "Delete Inline Protection Profile",
			"message":     "Failed to delete Inline Protection Profile",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteInlineProtection",
			"status":      "success",
			"description": "Delete Inline Protection Profile",
			"message":     "Successfully deleted Inline Protection Profile",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 3: Delete X-Forwarded-For Rule
	result, err = deleteXForwardedForRule(host, token, config.XForwardedForName)
	if err != nil {
		// ...
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteXForwardedForRule",
			"status":      "failure",
			"description": "Delete X-Forwarded-For Rule",
			"message":     fmt.Sprintf("Error deleting X-Forwarded-For Rule: %v", err),
		})
	} else if !checkOperationStatus(result) {
		// ...
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteXForwardedForRule",
			"status":      "failure",
			"description": "Delete X-Forwarded-For Rule",
			"message":     "Failed to delete X-Forwarded-For Rule",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteXForwardedForRule",
			"status":      "success",
			"description": "Delete X-Forwarded-For Rule",
			"message":     "Successfully deleted X-Forwarded-For Rule",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 4: Delete Signature Protection
	result, err = deleteSignatureProtection(host, token, config.CloneSignatureProtectionName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteSignatureProtection",
			"status":      "failure",
			"description": "Delete Signature Protection",
			"message":     fmt.Sprintf("Error deleting Signature Protection: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteSignatureProtection",
			"status":      "failure",
			"description": "Delete Signature Protection",
			"message":     "Failed to delete Signature Protection",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteSignatureProtection",
			"status":      "success",
			"description": "Delete Signature Protection",
			"message":     "Successfully deleted Signature Protection",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 5: DeleteVirtualServer
	result, err = deleteVirtualServer(host, token, config.VirtualServerName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualServer",
			"status":      "failure",
			"description": "Delete Virtual Server",
			"message":     fmt.Sprintf("Error deleting Virtual Server: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualServer",
			"status":      "failure",
			"description": "Delete Virtual Server",
			"message":     "Failed to delete Virtual Server",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualServer",
			"status":      "success",
			"description": "Delete Virtual Server",
			"message":     "Successfully deleted Virtual Server",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step : Delete Server Pool
	result, err = deleteServerPool(host, token, config.PoolName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteServerPool",
			"status":      "failure",
			"description": "Delete Server Pool",
			"message":     fmt.Sprintf("Error deleting Server Pool: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteServerPool",
			"status":      "failure",
			"description": "Delete Server Pool",
			"message":     "Failed to delete Server Pool",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteServerPool",
			"status":      "success",
			"description": "Delete Server Pool",
			"message":     "Successfully deleted Server Pool",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// Step 7: Delete Virtual IP
	result, err = deleteVirtualIP(host, token, config.VipName)
	if err != nil {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualIP",
			"status":      "failure",
			"description": "Delete Virtual IP",
			"message":     fmt.Sprintf("Error deleting virtual IP: %v", err),
		})
	} else if !checkOperationStatus(result) {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualIP",
			"status":      "failure",
			"description": "Delete Virtual IP",
			"message":     "Failed to delete Virtual IP",
		})
	} else {
		statuses = append(statuses, map[string]string{
			"taskId":      "deleteVirtualIP",
			"status":      "success",
			"description": "Delete Virtual IP",
			"message":     "Successfully deleted virtual IP",
		})
	}

	// Debug Statement
	printLastStatus(statuses)

	// log.Printf("End of deleteApplicationPolicy\n")
	// Return a JSON response with the statuses of all steps
	return c.JSON(http.StatusOK, statuses)
}

///////////////////////////////////////////////////////////////////////////////////
// DEBUG PRINT STATUSES                                                          //
///////////////////////////////////////////////////////////////////////////////////

func printLastStatus(statuses []map[string]string) {
	fmt.Println("Last Status:")
	if len(statuses) > 0 {
		lastStatus := statuses[len(statuses)-1]
		fmt.Printf("Task ID: %s\n", lastStatus["taskId"])
		fmt.Printf("Status: %s\n", lastStatus["status"])
		fmt.Printf("Description: %s\n", lastStatus["description"])
		fmt.Printf("Message: %s\n", lastStatus["message"])
		fmt.Println("------------------------")
	} else {
		fmt.Println("No statuses available.")
	}
}
