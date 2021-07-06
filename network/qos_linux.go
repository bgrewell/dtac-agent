package network

import (
	"errors"
	"fmt"
	"github.com/BGrewell/iptables"
	"github.com/google/uuid"
	"strings"
)

type DSCPRule struct {
	Chain           string                    `json:"chain" yaml:"chain" xml:"chain"`
	Protocol        iptables.InvertableString `json:"protocol,omitempty" yaml:"protocol" xml:"protocol"`
	Source          iptables.InvertableString `json:"source,omitempty" yaml:"source" xml:"source"`
	Destination     iptables.InvertableString `json:"destination,omitempty" yaml:"destination" xml:"destination"`
	SourcePort      iptables.InvertableString `json:"source_port,omitempty" yaml:"source_port" xml:"source_port"`
	DestinationPort iptables.InvertableString `json:"destination_port,omitempty" yaml:"destination_port" xml:"destination_port"`
	InputInterface  iptables.InvertableString `json:"input_interface,omitempty" yaml:"input_interface" xml:"input_interface"`
	OutputInterface iptables.InvertableString `json:"output_interface,omitempty" yaml:"output_interface" xml:"output_interface"`
	Target          iptables.TargetDSCP       `json:"target" yaml:"target" xml:"target"`
}

func IptablesAddDSCPRule(template *DSCPRule) (id string, err error) {
	id = uuid.New().String()
	if template.Chain == "" {
		return "", errors.New("required parameter 'chain' is missing")
	}
	rule := &iptables.Rule{
		Chain:           template.Chain,
		Protocol:        template.Protocol,
		Source:          template.Source,
		Destination:     template.Destination,
		SourcePort:      template.SourcePort,
		DestinationPort: template.DestinationPort,
		InputInterface:  template.InputInterface,
		OutputInterface: template.OutputInterface,
	}
	rule.Chain = strings.ToUpper(rule.Chain)
	rule.Id = id
	rule.App = app
	rule.Target = iptables.TargetDSCP{
		Value: template.Target.Value,
	}
	rule.Table = iptables.TableMangle
	rule.Debug = true
	err = rule.Append()
	if err != nil {
		return "", err
	}
	managedRules[id] = rule
	return id, nil
}

func IptablesDelDSCPRule(id string) (rule *DSCPRule, err error) {
	if rule, ok := managedRules[id]; ok {
		err = rule.Delete()
		if err != nil {
			return nil, err
		}
		delete(managedRules, id)
		ret := &DSCPRule{
			Chain:           rule.Chain,
			Protocol:        rule.Protocol,
			Source:          rule.Source,
			Destination:     rule.Destination,
			InputInterface:  rule.InputInterface,
			OutputInterface: rule.OutputInterface,
			Target:          rule.Target.(iptables.TargetDSCP),
		}
		return ret, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}
