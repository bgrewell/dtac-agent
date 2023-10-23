package commands

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/consts"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"
)

type tokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewConfigCmd returns a new instance of the config command for the dtac tool.
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
			cmd.OutOrStdout().Write([]byte("Token: " + string(body)))
		},
	}
}

func getTokens(cmd *cobra.Command, cfg *config.Configuration) ([]byte, bool) {
	scheme := "http"
	var transport *http.Transport
	if cfg.Listener.HTTPS.Enabled {
		scheme = "https"

		cert, err := tls.LoadX509KeyPair(cfg.Listener.HTTPS.CertFile, cfg.Listener.HTTPS.KeyFile)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
			return nil, true
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}

		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	port := cfg.Listener.Port
	apiEndpoint := fmt.Sprintf("%s://localhost:%d/auth/login", scheme, port)
	data := map[string]string{
		"username": cfg.Auth.User,
		"password": cfg.Auth.Pass,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, true
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, true
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, true
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
		return nil, true
	}
	return body, false
}
