package plugins

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Options is a struct that holds the startup options for a plugin
type Options struct {
	Encryption    bool
	EncryptionKey string
	TLSEnabled    bool
}

// ParseOptions parses the options string into an Options struct
func ParseOptions(optionsStr string) (options *Options, err error) {
	opts := &Options{}

	// Remove brackets
	optionsStr = strings.Trim(optionsStr, "[]")

	// Check for empty string
	if optionsStr == "" {
		return opts, nil
	}

	// Split the options
	pairs := strings.Split(optionsStr, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")

		// Check for flags (no "=" present)
		if len(kv) == 1 {
			switch kv[0] {
			// Supported flags
			default:
				return nil, errors.New("unsupported flag: " + kv[0])
			}
		} else if len(kv) == 2 {
			key := kv[0]
			value := kv[1]

			switch key {
			// Supported key=value pairs
			case "tls":
				enabled := false
				if strings.ToLower(value) == "true" {
					enabled = true
				}
				opts.TLSEnabled = enabled
			case "enc":
				opts.Encryption = true
				decodedKey, err := url.QueryUnescape(value)
				if err != nil {
					return nil, fmt.Errorf("failed to unescape encryption key: %s", err.Error())
				}
				opts.EncryptionKey = decodedKey
			default:
				return nil, errors.New("unsupported option key: " + key)
			}
		} else {
			return nil, errors.New("invalid option format")
		}
	}

	return opts, nil
}
