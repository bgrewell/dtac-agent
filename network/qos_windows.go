// +build windows

package network

import (
	"github.com/BGrewell/go-netqospolicy"
)

func GetNetQosPolicies() (policies []*netqospolicy.NetQoSPolicy, err error) {
	return netqospolicy.GetAll()
}

func GetNetQosPolicy(name string) (policy *netqospolicy.NetQoSPolicy, err error) {
	return netqospolicy.Get(name)
}

func CreateNetQosPolicy(policy *netqospolicy.NetQoSPolicy) (err error) {
	return policy.Create()
}

func UpdateNetQosPolicy(policy *netqospolicy.NetQoSPolicy) (err error) {
	return policy.Update()
}

func DeleteNetQosPolicy(name string) (err error) {
	return netqospolicy.Remove(name)
}

func DeleteNetQosPolicies() (err error) {
	return netqospolicy.RemoveAll()
}

func GetUniversalQosRule(id string) (r *UniversalDSCPRule, err error) {
	policy, err := netqospolicy.Get(id)
	if policy != nil {
		return nil, err
	}
	return ConvertWindowsQosToUniversalDSCPRule(policy), nil
}

func GetUniversalQosRules() (r []*UniversalDSCPRule, err error) {
	policies, err := netqospolicy.GetAll()
	if err != nil {
		return nil, err
	}
	rules := make([]*UniversalDSCPRule,0)
	for _, policy := range policies {
		ur := ConvertWindowsQosToUniversalDSCPRule(policy)
		rules = append(rules, ur)
	}
	return rules, nil
}

func CreateUniversalQosRule(rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	policy := rule.ToWindowsQos()
	err = policy.Create()
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func UpdateUniversalQosRule(id string, rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	policy := rule.ToWindowsQos()
	err = policy.Update()
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func DeleteUniversalQosRule(id string) (err error) {
	policy, err := netqospolicy.Get(id)
	if err != nil {
		return err
	}
	return policy.Remove()
}