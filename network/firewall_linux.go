package network

import (
	"fmt"
	"github.com/BGrewell/iptables"
	"github.com/google/uuid"
)

var (
	managedRules map[string]*iptables.Rule
)

func init() {
	managedRules = make(map[string]*iptables.Rule)
}

type FirewallRule struct {
	Protocol        iptables.InvertableString `json:"protocol,omitempty" yaml:"protocol" xml:"protocol"`
	Source          iptables.InvertableString `json:"source,omitempty" yaml:"source" xml:"source"`
	Destination     iptables.InvertableString `json:"destination,omitempty" yaml:"destination" xml:"destination"`
	InputInterface  iptables.InvertableString `json:"input_interface,omitempty" yaml:"input_interface" xml:"input_interface"`
	OutputInterface iptables.InvertableString `json:"output_interface,omitempty" yaml:"output_interface" xml:"output_interface"`
}

type DNATRule struct {
	Protocol        iptables.InvertableString `json:"protocol,omitempty" yaml:"protocol" xml:"protocol"`
	Source          iptables.InvertableString `json:"source,omitempty" yaml:"source" xml:"source"`
	Destination     iptables.InvertableString `json:"destination,omitempty" yaml:"destination" xml:"destination"`
	InputInterface  iptables.InvertableString `json:"input_interface,omitempty" yaml:"input_interface" xml:"input_interface"`
	OutputInterface iptables.InvertableString `json:"output_interface,omitempty" yaml:"output_interface" xml:"output_interface"`
	Target          iptables.TargetDNat       `json:"target" yaml:"target" xml:"target"`
}

type SNATRule struct {
	Protocol        iptables.InvertableString `json:"protocol,omitempty" yaml:"protocol" xml:"protocol"`
	Source          iptables.InvertableString `json:"source,omitempty" yaml:"source" xml:"source"`
	Destination     iptables.InvertableString `json:"destination,omitempty" yaml:"destination" xml:"destination"`
	InputInterface  iptables.InvertableString `json:"input_interface,omitempty" yaml:"input_interface" xml:"input_interface"`
	OutputInterface iptables.InvertableString `json:"output_interface,omitempty" yaml:"output_interface" xml:"output_interface"`
	Target          iptables.TargetSNat       `json:"target" yaml:"target" xml:"target"`
}

func IptablesAddDNatRule(template *DNATRule) (id string, err error) {
	id = uuid.New().String()
	rule := &iptables.Rule{
		Protocol:        template.Protocol,
		Source:          template.Source,
		Destination:     template.Destination,
		InputInterface:  template.InputInterface,
		OutputInterface: template.OutputInterface,
	}
	rule.Id = fmt.Sprintf("system-api:%s", id)
	rule.Target = iptables.TargetDNat{
		DestinationIp:        template.Target.DestinationIp,
		DestinationIpRange:   template.Target.DestinationIpRange,
		DestinationPort:      template.Target.DestinationPort,
		DestinationPortRange: template.Target.DestinationPortRange,
	}
	rule.Table = iptables.TableNat
	rule.Chain = iptables.ChainPreRouting
	rule.Debug = true
	err = rule.Append()
	if err != nil {
		return "", err
	}
	managedRules[id] = rule
	return id, nil
}

func IptablesUpdateDNatRule(id string, template *DNATRule) (rule *DNATRule, err error) {
	if rule, ok := managedRules[id]; ok {
		if template.Source.Value != "" {
			rule.Source = template.Source
		}
		if template.Destination.Value != "" {
			rule.Destination = template.Destination
		}
		if template.Protocol.Value != "" {
			rule.Protocol = template.Protocol
		}
		if template.InputInterface.Value != "" {
			rule.InputInterface = template.InputInterface
		}
		if template.OutputInterface.Value != "" {
			rule.OutputInterface = template.OutputInterface
		}
		if template.Target.DestinationIp != "" || template.Target.DestinationIpRange != "" ||
			template.Target.DestinationPort != "" || template.Target.DestinationPortRange != "" {
			rule.Target = template.Target
		}
		rule.Replace()
		if err != nil {
			return nil, err
		}
		ret := &DNATRule{
			Protocol:        rule.Protocol,
			Source:          rule.Source,
			Destination:     rule.Destination,
			InputInterface:  rule.InputInterface,
			OutputInterface: rule.OutputInterface,
			Target:          rule.Target.(iptables.TargetDNat),
		}
		return ret, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}

func IptablesDelDNatRule(id string) (rule *DNATRule, err error) {
	if rule, ok := managedRules[id]; ok {
		err = rule.Delete()
		if err != nil {
			return nil, err
		}
		delete(managedRules, id)
		ret := &DNATRule{
			Protocol:        rule.Protocol,
			Source:          rule.Source,
			Destination:     rule.Destination,
			InputInterface:  rule.InputInterface,
			OutputInterface: rule.OutputInterface,
			Target:          rule.Target.(iptables.TargetDNat),
		}
		return ret, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}

func IptablesAddSNatRule(template *SNATRule) (id string, err error) {
	id = uuid.New().String()
	rule := &iptables.Rule{
		Protocol:        template.Protocol,
		Source:          template.Source,
		Destination:     template.Destination,
		InputInterface:  template.InputInterface,
		OutputInterface: template.OutputInterface,
	}
	rule.Id = fmt.Sprintf("system-api:%s", id)
	rule.Target = iptables.TargetSNat{
		SourceIp:        template.Target.SourceIp,
		SourceIpRange:   template.Target.SourceIpRange,
		SourcePort:      template.Target.SourcePort,
		SourcePortRange: template.Target.SourcePortRange,
	}
	rule.Table = iptables.TableNat
	rule.Chain = iptables.ChainPreRouting
	err = rule.Append()
	if err != nil {
		return "", err
	}
	managedRules[id] = rule
	return id, nil
}

func IptablesUpdateSNatRule(id string, template *SNATRule) (rule *SNATRule, err error) {
	if rule, ok := managedRules[id]; ok {
		if template.Source.Value != "" {
			rule.Source = template.Source
		}
		if template.Destination.Value != "" {
			rule.Destination = template.Destination
		}
		if template.Protocol.Value != "" {
			rule.Protocol = template.Protocol
		}
		if template.InputInterface.Value != "" {
			rule.InputInterface = template.InputInterface
		}
		if template.OutputInterface.Value != "" {
			rule.OutputInterface = template.OutputInterface
		}
		if template.Target.SourceIp != "" || template.Target.SourceIpRange != "" ||
			template.Target.SourcePort != "" || template.Target.SourcePortRange != "" {
			rule.Target = template.Target
		}
		rule.Replace()
		if err != nil {
			return nil, err
		}
		ret := &SNATRule{
			Protocol:        rule.Protocol,
			Source:          rule.Source,
			Destination:     rule.Destination,
			InputInterface:  rule.InputInterface,
			OutputInterface: rule.OutputInterface,
			Target:          rule.Target.(iptables.TargetSNat),
		}
		return ret, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}

func IptablesDelSNatRule(id string) (rule *SNATRule, err error) {
	if rule, ok := managedRules[id]; ok {
		err = rule.Delete()
		if err != nil {
			return nil, err
		}
		delete(managedRules, id)
		ret := &SNATRule{
			Protocol:        rule.Protocol,
			Source:          rule.Source,
			Destination:     rule.Destination,
			InputInterface:  rule.InputInterface,
			OutputInterface: rule.OutputInterface,
			Target:          rule.Target.(iptables.TargetSNat),
		}
		return ret, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}
