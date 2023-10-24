package config

import "path"

const (
	GlobalConfigLocation  = "/etc/dtac/"
	LocalConfigLocation   = "$HOME/.dtac/"
	DefaultBinaryLocation = "/opt/dtac/bin/"
	DefaultPluginLocation = "/opt/dtac/plugins/"
)

var (
	GlobalCertLocation    = path.Join(GlobalConfigLocation, "certs/")
	LocalCertLocation     = path.Join(LocalConfigLocation, "certs/")
	GlobalDBLocation      = path.Join(GlobalConfigLocation, "db/")
	BinaryName            = path.Join(DefaultBinaryLocation, "dtac-agentd")
	DBName                = path.Join(GlobalDBLocation, "authn.db")
	DefaultTLSCACertName  = path.Join(GlobalCertLocation, "ca.crt")
	DefaultTLSCertName    = path.Join(GlobalCertLocation, "tls.crt")
	DefaultTLSKeyName     = path.Join(GlobalCertLocation, "tls.key")
	DefaultAuthModelName  = path.Join(GlobalConfigLocation, "auth_model.conf")
	DefaultAuthPolicyName = path.Join(GlobalConfigLocation, "auth_policy.csv")
)
