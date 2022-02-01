package remote

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type RemoteModule struct {
	Host  string `json:"host,omitempty" yaml:"host"`
	Port  int    `json:"port,omitempty" yaml:"port"`
	Token string `json:"token,omitempty" yaml:"token"`
}

func (p *RemoteModule) Register(config map[string]interface{}, r *gin.RouterGroup) error {

	// process config
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, p)
	if err != nil {
		return err
	}
	// Register any handlers here
	return nil
}

func (p *RemoteModule) Name() string {
	return "remote"
}
