package config

import "path"

const (
	GlobalConfigLocation = "c:\\Program Files\\Intel\\dtac-agent\\"
	LocalConfigLocation  = "$HOME/.dtac/"
)

var (
	GlobalCertLocation    = path.Join(GlobalConfigLocation, "certs\\")
	LocalCertLocation     = path.Join(LocalConfigLocation, "certs\\")
	DefaultBinaryLocation = path.Join(GlobalConfigLocation, "bin\\")
	DefaultPluginLocation = path.Join(GlobalConfigLocation, "plugins\\")
	BinaryName            = path.Join(DefaultBinaryLocation, "dtac-agentd")
	GlobalDBLocation      = path.Join(GlobalConfigLocation, "db\\")
	DBName                = path.Join(GlobalDBLocation, "authn.db")
	DefaultTLSCACertName  = path.Join(GlobalCertLocation, "ca.crt")
	DefaultTLSCertName    = path.Join(GlobalCertLocation, "tls.crt")
	DefaultTLSKeyName     = path.Join(GlobalCertLocation, "tls.key")
	DefaultAuthModelName  = path.Join(GlobalConfigLocation, "auth_model.conf")
	DefaultAuthPolicyName = path.Join(GlobalConfigLocation, "auth_policy.csv")
)
