package config

var (
	VipName                             = "SPEEDTEST_VIP"
	VipIp                               = "10.163.7.34/32"
	Interface                           = "port1"
	PoolName                            = "SPEEDTEST_POOL"
	ServerBalance                       = "enable"
	HealthCheck                         = "HLTHCK_HTTP"
	PoolMemberIPs                       = []string{"10.163.3.78", "10.163.3.79", "10.163.3.80"}
	PoolMemberSSL                       = "disable"
	PoolMemberPort                      = 5000
	VirtualServerName                   = "SPEEDTEST_VS"
	VipStatus                           = "enable"
	XForwardedForName                   = "SPEEDTEST_X_FORWARDED_FOR"
	XForwardedForSupport                = "enable"
	OriginalSignatureProtectionName     = "Standard Protection"
	CloneSignatureProtectionName        = "SPEEDTEST_SIGNATURE_CLONE"
	OriginalInlineProtectionProfileName = "Inline Standard Protection"
	CloneInlineProtectionProfileName    = "SPEEDTEST_PROTECTION_CLONE"
	PolicyName                          = "SPEEDTEST_POLICY"
	PolicyDeploymentMode                = "server-pool"
	PolicyProtocol                      = "HTTP"
	PolicySSL                           = "enable"
	PolicyImplicitSSL                   = "enable"
	PolicyVirtualServer                 = VirtualServerName
	PolicyService                       = "HTTP"
	PolicyInlineProtectionProfile       = CloneInlineProtectionProfileName
	PolicyServerPool                    = PoolName
	PolicyTrafficLog                    = "enable"
	PolicyHTTPSService                  = "HTTPS"
	//PolicyCertificate                   = "speedtest.canadaeast.cloudapp.azure.com"
)

// Data Types Struct

type VirtualIPData struct {
	Name      string `json:"name,omitempty"`
	Vip       string `json:"vip,omitempty"`
	Interface string `json:"interface,omitempty"`
}

type ServerPoolData struct {
	Name          string `json:"name,omitempty"`
	ServerBalance string `json:"server-balance,omitempty"`
	Health        string `json:"health,omitempty"`
}

type MemberPoolData struct {
	IP   string `json:"ip,omitempty"`
	SSL  string `json:"ssl,omitempty"`
	Port int    `json:"port,omitempty"`
}

type VirtualServerData struct {
	Name string `json:"name,omitempty"`
}

type AssignVIPData struct {
	Interface string `json:"interface,omitempty"`
	Status    string `json:"status,omitempty"`
	VipName   string `json:"vip,omitempty"`
}

type Request struct {
	Data interface{} `json:"data"`
}

type XForwardedForData struct {
	Name                 string `json:"name,omitempty"`
	XForwardedForSupport string `json:"x-forwarded-for-support,omitempty"`
}

type ProtectionProfileData struct {
	SignatureRule     string `json:"signature-rule,omitempty"`
	XForwardedForRule string `json:"x-forwarded-for-rule,omitempty"`
}

type PolicyData struct {
	Name                    string `json:"name,omitempty"`
	DeploymentMode          string `json:"deployment-mode,omitempty"`
	Protocol                string `json:"protocol,omitempty"`
	Ssl                     string `json:"ssl,omitempty"`
	ImplicitSsl             string `json:"implicit_ssl,omitempty"`
	Vserver                 string `json:"vserver,omitempty"`
	Service                 string `json:"service,omitempty"`
	InlineProtectionProfile string `json:"web-protection-profile,omitempty"`
	ServerPool              string `json:"server-pool,omitempty"`
	TrafficLog              string `json:"tlog,omitempty"`
	HttpsService            string `json:"https-service,omitempty"`
	Certificate             string `json:"certificate,omitempty"`
}
