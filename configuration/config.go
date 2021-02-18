package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type SourceEntry struct {
	Type  string `json:"type" yaml:"type" xml:"type"`
	Value string `json:"value" yaml:"value" xml:"value"`
}

type BlockingEntry struct {
	Trigger   string `json:"trigger" yaml:"trigger" xml:"trigger"`
	Detect    string `json:"detect" yaml:"detect" xml:"detect"`
	TimeoutMs int    `json:"timeout_ms" yaml:"timeout_ms" xml:"timeout_ms"`
}

type CustomEntry struct {
	Name     string        `json:"name" yaml:"name" xml:"name"`
	Source   *SourceEntry   `json:"source" yaml:"source" xml:"source"`
	Methods  []string      `json:"methods" yaml:"methods" xml:"methods"`
	Access   string        `json:"access" yaml:"access" xml:"access"`
	Blocking *BlockingEntry `json:"blocking" yaml:"blocking" xml:"blocking"`
	WrapJson bool          `json:"wrap_json" yaml:"wrap_json" xml:"wrap_json"`
}

type Config struct {
	ListenPort int                    `json:"listen_port" yaml:"listen_port" xml:"listen_port"`
	HTTPS      bool                   `json:"https" yaml:"https" xml:"https"`
	CertFile   string                 `json:"cert_file" yaml:"cert_file" xml:"cert_file"`
	KeyFile    string                 `json:"key_file" yaml:"key_file" xml:"key_file"`
	Custom     []map[string]*CustomEntry `json:"custom" yaml:"custom" xml:"custom"`
}

func Load(filename string) (config *Config, err error) {
	c := &Config{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	// populate names
	for _, entry := range c.Custom {
		for key, value := range entry {
			value.Name = key
		}
	}
	return c, nil
}
