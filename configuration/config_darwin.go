package configuration

import "path"

const (
	GLOBAL_CONFIG_LOCATION  = "/etc/dtac/"
	LOCAL_CONFIG_LOCATION   = "$HOME/.dtac/"
	DEFAULT_BINARY_LOCATION = "/opt/dtac/bin/"
	DEFAULT_PLUGIN_LOCATION = "/opt/dtac/plugins/"
)

var (
	GLOBAL_CERT_LOCATION  = path.Join(GLOBAL_CONFIG_LOCATION, "certs/")
	LOCAL_CERT_LOCATION   = path.Join(LOCAL_CONFIG_LOCATION, "certs/")
	BINARY_NAME           = path.Join(DEFAULT_BINARY_LOCATION, "dtac-agentd")
	GLOBAL_DB_LOCATION    = path.Join(GLOBAL_CONFIG_LOCATION, "db/")
	DB_NAME               = path.Join(GLOBAL_DB_LOCATION, "auth.db")
	DEFAULT_TLS_CERT_NAME = path.Join(GLOBAL_CERT_LOCATION, "tls.crt")
	DEFAULT_TLS_KEY_NAME  = path.Join(GLOBAL_CERT_LOCATION, "tls.key")
)
