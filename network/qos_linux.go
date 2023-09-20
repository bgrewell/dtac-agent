package network

import (
	"errors"
	"fmt"
	. "github.com/intel-innersource/frameworks.automation.dtac.agent/common"
	"github.com/BGrewell/go-iptables"
	"github.com/google/uuid"
)

type DSCPTemplate struct {
	Name                   string              `json:"name,omitempty" yaml:"name" xml:"name"`
	Chain                  string              `json:"chain,omitempty" yaml:"chain" xml:"chain"`
	Protocol               string              `json:"protocol,omitempty" yaml:"protocol" xml:"protocol"`
	ProtocolNegated        bool                `json:"protocol_negated,omitempty" yaml:"protocol_negated" xml:"protocol_negated"`
	Source                 string              `json:"source,omitempty" yaml:"source" xml:"source"`
	SourceNegated          bool                `json:"source_negated,omitempty" yaml:"source_negated" xml:"source_negated"`
	Destination            string              `json:"destination,omitempty" yaml:"destination" xml:"destination"`
	DestinationNegated     bool                `json:"destination_negated,omitempty" yaml:"destination_negated" xml:"destination_negated"`
	SourcePort             string              `json:"source_port,omitempty" yaml:"source_port" xml:"source_port"`
	SourcePortNegated      bool                `json:"source_port_negated,omitempty" yaml:"source_port_negated" xml:"source_port_negated"`
	DestinationPort        string              `json:"destination_port,omitempty" yaml:"destination_port" xml:"destination_port"`
	DestinationPortNegated bool                `json:"destination_port_negated,omitempty" yaml:"destination_port_negated" xml:"destination_port_negated"`
	InputInterface         string              `json:"input_interface,omitempty" yaml:"input_interface" xml:"input_interface"`
	InputInterfaceNegated  bool                `json:"input_interface_negated,omitempty" yaml:"input_interface_negated" xml:"input_interface_negated"`
	OutputInterface        string              `json:"output_interface,omitempty" yaml:"output_interface" xml:"output_interface"`
	OutputInterfaceNegated bool                `json:"output_interface_negated,omitempty" yaml:"output_interface_negated" xml:"output_interface_negated"`
	Target                 iptables.TargetDSCP `json:"target,omitempty" yaml:"target" xml:"target"`
}

func IptablesGetDSCPRules() (outRules []*iptables.Rule, err error) {
	t := iptables.TargetDSCP{}
	return iptables.GetRulesByTarget(&t)
}

func IptablesGetDSCPRule(id string) (outRule *iptables.Rule, err error) {
	return iptables.GetRuleById(id)
}

func IptablesAddDSCPRule(template *DSCPTemplate) (outRule *iptables.Rule, err error) {
	id := uuid.New().String()
	if template.Chain == "" {
		return nil, errors.New("required parameter 'chain' is missing")
	}
	rule := &iptables.Rule{
		Name:                   template.Name,
		Chain:                  iptables.Chain(template.Chain),
		Protocol:               iptables.Protocol(template.Protocol),
		ProtocolNegated:        template.ProtocolNegated,
		Source:                 template.Source,
		SourceNegated:          template.SourceNegated,
		Destination:            template.Destination,
		DestinationNegated:     template.DestinationNegated,
		SourcePort:             template.SourcePort,
		SourcePortNegated:      template.SourcePortNegated,
		DestinationPort:        template.DestinationPort,
		DestinationPortNegated: template.DestinationPortNegated,
		Input:                  template.InputInterface,
		InputNegated:           template.InputInterfaceNegated,
		Output:                 template.OutputInterface,
		OutputNegated:          template.OutputInterfaceNegated,
		Target:                 &template.Target,
	}
	rule.Id = id
	rule.IpVersion = iptables.IPv4
	rule.Debug = true
	rule.SetApp(app)
	m := iptables.MarkerGeneric{}
	m.SetName("type")
	m.SetValue("dscp")
	rule.AddMarker(&m)
	rule.Table = iptables.TableMangle
	err = rule.Append()
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func IptablesUpdateDSCPRule(id string, template *DSCPTemplate) (outRule *iptables.Rule, err error) {
	rule := &iptables.Rule{
		Source:                 template.Source,
		SourceNegated:          template.SourceNegated,
		Destination:            template.Destination,
		DestinationNegated:     template.DestinationNegated,
		SourcePort:             template.SourcePort,
		SourcePortNegated:      template.SourcePortNegated,
		DestinationPort:        template.DestinationPort,
		DestinationPortNegated: template.DestinationPortNegated,
		Input:                  template.InputInterface,
		InputNegated:           template.InputInterfaceNegated,
		Output:                 template.OutputInterface,
		OutputNegated:          template.OutputInterfaceNegated,
	}
	if _, err = iptables.FindRuleById(id); err == nil {
		r, err := iptables.GetRuleById(id)
		if err != nil {
			return nil, err
		}
		r.Update(rule)
		err = r.Replace()
		return r, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}

func IptablesDelDSCPRule(id string) (rule *iptables.Rule, err error) {
	if _, err = iptables.FindRuleById(id); err == nil {
		r, err := iptables.GetRuleById(id)
		if err != nil {
			return nil, err
		}
		err = r.Delete()
		if err != nil {
			return nil, err
		}
		return r, err
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}

func GetUniversalQosRule(id string) (r *UniversalDSCPRule, err error) {
	var iptablesRule *iptables.Rule
	if !IsValidUUID(id) {
		iptablesRule, err = iptables.GetRuleByName(id)
	} else {
		iptablesRule, err = iptables.GetRuleById(id)
	}
	if err != nil {
		return nil, err
	}
	return ConvertLinuxQosToUniversalDSCPRule(iptablesRule), nil
}

func GetUniversalQosRules() (r []*UniversalDSCPRule, err error) {
	iptablesRules, err := iptables.GetRulesByTarget(&iptables.TargetDSCP{})
	if err != nil {
		return nil, err
	}
	rules := make([]*UniversalDSCPRule, 0)
	for _, rule := range iptablesRules {
		ur := ConvertLinuxQosToUniversalDSCPRule(rule)
		rules = append(rules, ur)
	}
	return rules, nil
}

func CreateUniversalQosRule(rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	iptablesRule := rule.ToLinuxQos()
	err = iptablesRule.Append()
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func UpdateUniversalQosRule(id string, rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	iptablesRule := rule.ToLinuxQos()
	var originalRule *iptables.Rule
	if !IsValidUUID(id) {
		originalRule, err = iptables.GetRuleByName(id)
	} else {
		originalRule, err = iptables.GetRuleById(id)
	}
	if err != nil {
		return nil, err
	}
	originalRule.Update(iptablesRule)
	t := &iptables.TargetDSCP{Value: rule.DSCP}
	originalRule.Target = t
	err = originalRule.Replace()
	if err != nil {
		return nil, err
	}
	return ConvertLinuxQosToUniversalDSCPRule(originalRule), nil
}

func DeleteUniversalQosRule(id string) (err error) {
	var iptablesRule *iptables.Rule
	if !IsValidUUID(id) {
		iptablesRule, err = iptables.GetRuleByName(id)
	} else {
		iptablesRule, err = iptables.GetRuleById(id)
	}

	if err != nil {
		return err
	}
	return iptablesRule.Delete()
}
