package structs

// MAASSettings is just a simple helper struct to encapsulate the MAAS server information
type MAASSettings struct {
	Server          string `json:"server"`
	ConsumerToken   string `json:"consumer_token"`
	AuthToken       string `json:"auth_token"`
	AuthSignature   string `json:"auth_signature"`
	MachinePollSecs int    `json:"machine_poll_secs"`
}
