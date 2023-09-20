package api

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/plugin/maas/maas_plugin/structs"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func Get(endpoint string, settings *structs.MAASSettings) ([]byte, error) {
	ip, err := net.LookupIP(settings.Server)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:5240/MAAS/api/2.0/%s", ip[0].String(), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	nonce := strconv.Itoa(rand.Intn(10e9))
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	authHeader := fmt.Sprintf("OAuth oauth_consumer_key=\"%s\",oauth_token=\"%s\",oauth_signature_method=\"PLAINTEXT\",oauth_timestamp=\"%s\",oauth_nonce=\"%s\",oauth_version=\"1.0\",oauth_signature=\"%%26%s\"", settings.ConsumerToken, settings.AuthToken, timestamp, nonce, settings.AuthSignature)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
