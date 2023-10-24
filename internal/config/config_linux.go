package config

import "path"

const (
	// GlobalConfigLocation System wide configuration location
	GlobalConfigLocation = "/etc/dtac/"
	// LocalConfigLocation User specific configuration location
	LocalConfigLocation = "$HOME/.dtac/"
	// DefaultBinaryLocation Default location of the dtac-agentd binary
	DefaultBinaryLocation = "/opt/dtac/bin/"
	// DefaultPluginLocation Default location of the dtac-agentd plugins
	DefaultPluginLocation = "/opt/dtac/plugins/"
)

var (
	// GlobalCertLocation System wide certificate location
	GlobalCertLocation = path.Join(GlobalConfigLocation, "certs/")
	// LocalCertLocation User specific certificate location
	LocalCertLocation = path.Join(LocalConfigLocation, "certs/")
	// GlobalDBLocation System wide database location
	GlobalDBLocation = path.Join(GlobalConfigLocation, "db/")
	// BinaryName is the name of the binary
	BinaryName = path.Join(DefaultBinaryLocation, "dtac-agentd")
	// DBName is the name of the database
	DBName = path.Join(GlobalDBLocation, "authn.db")
	// DefaultTLSCACertNAme Default TLS CA certificate name
	DefaultTLSCACertName = path.Join(GlobalCertLocation, "ca.crt")
	// DefaultTLSCertName Default TLS certificate name
	DefaultTLSCertName = path.Join(GlobalCertLocation, "tls.crt")
	// DefaultTLSKeyName Default TLS key name
	DefaultTLSKeyName = path.Join(GlobalCertLocation, "tls.key")
	// DefaultAuthModelName Default auth model name
	DefaultAuthModelName = path.Join(GlobalConfigLocation, "auth_model.conf")
	// DefaultAuthPolicyName Default auth policy name
	DefaultAuthPolicyName = path.Join(GlobalConfigLocation, "auth_policy.csv")
)
