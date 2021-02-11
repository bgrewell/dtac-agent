package network

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/BGrewell/iptables"
)

var (
	activeRules map[string]*iptables.Rule
)

func init() {
	activeRules = make(map[string]*iptables.Rule)
}

func AddIptablesDNatRule(rule *iptables.Rule) (id string, err error) {
	id = uuid.New().String()
	rule.Target = iptables.TargetDNat{}
	err = rule.Append()
	if err != nil {
		return "", err
	}
	activeRules[id] = rule
	return id, nil
}

func DelIptablesDNatRule(id string) (err error) {
	if rule, ok := activeRules[id]; ok {
		err = rule.Delete()
		if err != nil {
			return err
		}
		delete(activeRules, id)
		return nil
	} else {
		return fmt.Errorf("failed to find rule matching the supplied id: %s", id)
	}
}
