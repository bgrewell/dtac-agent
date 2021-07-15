// +build windows

package network

import "github.com/BGrewell/go-netqospolicy"

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
	return nil, errors.New("this method has not been implemented for windows yet")
}

func GetUniversalQosRules() (r []*UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for windows yet")
}

func CreateUniversalQosRule(rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for windows yet")
}

func UpdateUniversalQosRule(rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for windows yet")
}

func DeleteUniversalQosRule(id string) (err error) {
	return errors.New("this method has not been implemented for windows yet")
}