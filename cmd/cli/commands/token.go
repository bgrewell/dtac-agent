package commands

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/consts"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/spf13/cobra"
)

type tokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewTokenCmd returns a new instance of the token command.
func NewTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Get an API request token",
		Run: func(cmd *cobra.Command, args []string) {
			c := cmd.Context().Value(consts.KeyConfig)
			if c == nil {
				cmd.ErrOrStderr().Write([]byte("Error: config not found"))
				return
			}
			cfg := c.(*config.Configuration)

			body, done := getTokens(cmd, cfg)
			if done {
				return
			}

			var response helpers.ResponseWrapper
			err := json.Unmarshal(body, &response)
			if err != nil {
				cmd.ErrOrStderr().Write([]byte("Failed to unmarshal response: " + err.Error()))
				return
			}

			var tokens tokenDetails
			err = json.Unmarshal(response.Response, &tokens)
			if err != nil {
				cmd.ErrOrStderr().Write([]byte("Failed to unmarshal access token: " + err.Error()))
				return
			}

			cmd.OutOrStdout().Write([]byte(tokens.AccessToken))
		},
	}
}

func getTokens(cmd *cobra.Command, cfg *config.Configuration) ([]byte, bool) {
	scheme := "http"
	var transport *http.Transport
	if cfg.APIs.REST.TLS.Enabled {
		scheme = "https"

		profile := cfg.APIs.REST.TLS.Profile
		var ok bool
		var tlsCfg config.TLSConfigurationEntry

		if tlsCfg, ok = cfg.TLS[profile]; !ok {
			cmd.ErrOrStderr().Write([]byte("Error: TLS profile not found"))
			return nil, false
		}

		cert, err := tls.LoadX509KeyPair(tlsCfg.CertFile, tlsCfg.KeyFile)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
			return nil, false
		}

		caCert, err := os.ReadFile(tlsCfg.CAFile)
		if err != nil {
			return nil, false
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, false
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}

		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	port := cfg.APIs.REST.Port
	apiEndpoint := fmt.Sprintf("%s://localhost:%d/auth/login", scheme, port)
	data := map[string]string{
		"username": cfg.Auth.User,
		"password": cfg.Auth.Pass,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, false
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, false
	}
	return body, false
}
