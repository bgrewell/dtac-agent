package api

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/bgrewell/dtac-agent/cmd/plugins/maas/maasplugin/structs"
)

// doRequest performs GET or DELETE requests to the MAAS server with OAuth authentication.
func doRequest(method, endpoint string, settings *structs.MAASSettings) ([]byte, error) {
	// Resolve MAAS server hostname
	ipList, err := net.LookupIP(settings.Server)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	url := fmt.Sprintf("http://%s:5240/MAAS/api/2.0/%s", ipList[0].String(), endpoint)

	// Create HTTP request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Build OAuth header
	nonce := strconv.Itoa(rand.Intn(1e10))
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	authHeader := fmt.Sprintf(
		"OAuth oauth_consumer_key=\"%s\",oauth_token=\"%s\",oauth_signature_method=\"PLAINTEXT\","+
			"oauth_timestamp=\"%s\",oauth_nonce=\"%s\",oauth_version=\"1.0\",oauth_signature=\"%%26%s\"",
		settings.ConsumerToken,
		settings.AuthToken,
		timestamp,
		nonce,
		settings.AuthSignature,
	)
	req.Header.Set("Authorization", authHeader)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and return the response body
	return io.ReadAll(resp.Body)
}

// doFormRequest performs POST or PUT requests with multipart/form-data.
func doFormRequest(method, endpoint string, fields map[string]string, settings *structs.MAASSettings) ([]byte, error) {
	// Resolve MAAS server hostname
	ipList, err := net.LookupIP(settings.Server)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	url := fmt.Sprintf("http://%s:5240/MAAS/api/2.0/%s", ipList[0].String(), endpoint)

	// Prepare multipart form data
	var bodyBuf bytes.Buffer
	writer := multipart.NewWriter(&bodyBuf)
	for key, value := range fields {
		_ = writer.WriteField(key, value)
	}
	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest(method, url, &bodyBuf)
	if err != nil {
		return nil, err
	}

	// Build OAuth header
	nonce := strconv.Itoa(rand.Intn(1e10))
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	authHeader := fmt.Sprintf(
		"OAuth oauth_consumer_key=\"%s\",oauth_token=\"%s\",oauth_signature_method=\"PLAINTEXT\","+
			"oauth_timestamp=\"%s\",oauth_nonce=\"%s\",oauth_version=\"1.0\",oauth_signature=\"%%26%s\"",
		settings.ConsumerToken,
		settings.AuthToken,
		timestamp,
		nonce,
		settings.AuthSignature,
	)
	req.Header.Set("Authorization", authHeader)

	// Set multipart content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and return the response body
	return io.ReadAll(resp.Body)
}

// Get performs a GET request to the specified MAAS endpoint.
func Get(endpoint string, settings *structs.MAASSettings) ([]byte, error) {
	return doRequest(http.MethodGet, endpoint, settings)
}

// Delete performs a DELETE request to the specified MAAS endpoint.
func Delete(endpoint string, settings *structs.MAASSettings) ([]byte, error) {
	return doRequest(http.MethodDelete, endpoint, settings)
}

// Post performs a multipart/form-data POST request to the specified MAAS endpoint.
func Post(endpoint string, fields map[string]string, settings *structs.MAASSettings) ([]byte, error) {
	return doFormRequest(http.MethodPost, endpoint, fields, settings)
}

// Put performs a multipart/form-data PUT request to the specified MAAS endpoint.
func Put(endpoint string, fields map[string]string, settings *structs.MAASSettings) ([]byte, error) {
	return doFormRequest(http.MethodPut, endpoint, fields, settings)
}
