package network

import (
	"errors"
)

func GetUniversalQosRule(id string) (r *UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for darwin yet")
}

func GetUniversalQosRules() (r []*UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for darwin yet")
}

func CreateUniversalQosRule(rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for darwin yet")
}

func UpdateUniversalQosRule(id string, rule *UniversalDSCPRule) (r *UniversalDSCPRule, err error) {
	return nil, errors.New("this method has not been implemented for darwin yet")
}

func DeleteUniversalQosRule(id string) (err error) {
	return errors.New("this method has not been implemented for darwin yet")
}
