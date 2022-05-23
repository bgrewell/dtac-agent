package network

import (
	"fmt"
	"github.com/BGrewell/go-iptables"
	"github.com/google/uuid"
)

var (
)

func init() {
	// make sure all rules have an id
	iptables.LabelRules()
}

// TODO: Need to create SNAT/DNAT templates instead of taking in full iptables.Rules

func IptablesGetStatus() (status *iptables.Status, err error) {
	return iptables.GetStatus()
}

func IptablesGetDNatRules() (outRules []*iptables.Rule, err error) {
	t := iptables.TargetDNat{}
	return iptables.GetRulesByTarget(&t)
}

func IptablesGetDNatRule(id string) (outRule *iptables.Rule, err error) {
	return iptables.GetRuleById(id)
}

func IptablesAddDNatRule(inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	id := uuid.New().String()
	inRule.Id = id
	inRule.SetApp(app)
	inRule.Table = iptables.TableNat
	inRule.Chain = iptables.ChainPreRouting
	m := iptables.MarkerGeneric{}
	m.SetName("type")
	m.SetValue("dnat")
	inRule.AddMarker(&m)
	err = inRule.Append()
	if err != nil {
		return nil, err
	}
	return inRule, nil
}

func IptablesUpdateDNatRule(inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	if _, err = iptables.FindRuleById(inRule.Id); err == nil {
		r, err := iptables.GetRuleById(inRule.Id)
		if err != nil {
			return nil, err
		}
		r.Update(inRule)
		err = r.Replace()
		return inRule, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", inRule.Id)
	}
}

func IptablesDelDNatRules() (outRules []*iptables.Rule, err error) {
	return nil, fmt.Errorf("this function is not implemented yet")
}

func IptablesDelDNatRule(id string) (outRule *iptables.Rule, err error) {
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

func IptablesGetSNatRules() (outRules []*iptables.Rule, err error) {
	t := iptables.TargetSNat{}
	return iptables.GetRulesByTarget(&t)
}

func IptablesGetSNatRule(id string) (outRule *iptables.Rule, err error) {
	return iptables.GetRuleById(id)
}

func IptablesAddSNatRule(inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	id := uuid.New().String()
	inRule.Id = id
	inRule.SetApp(app)
	inRule.Table = iptables.TableNat
	inRule.Chain = iptables.ChainPreRouting
	m := iptables.MarkerGeneric{}
	m.SetName("type")
	m.SetValue("snat")
	inRule.AddMarker(&m)
	err = inRule.Append()
	if err != nil {
		return nil, err
	}
	return inRule, nil
}

func IptablesUpdateSNatRule(inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	if _, err = iptables.FindRuleById(inRule.Id); err == nil {
		r, err := iptables.GetRuleById(inRule.Id)
		if err != nil {
			return nil, err
		}
		r.Update(inRule)
		err = r.Replace()
		return inRule, nil
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", inRule.Id)
	}
}

func IptablesDelSNatRules() (outRules []*iptables.Rule, err error) {
	return nil, fmt.Errorf("this function is not implemented yet")
}

func IptablesDelSNatRule(id string) (outRule *iptables.Rule, err error) {
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

func IptablesDelRules() (err error) {
	return iptables.DeleteByComment(app)
}

func IptablesDelRule(id string) (outRule *iptables.Rule, err error) {
	r, err := iptables.GetRuleById(id)
	if err != nil {
		return nil, err
	}
	return r, iptables.DeleteById(id)
}

func IptablesGetRules() (outRules []*iptables.Rule, err error) {
	return iptables.Sync()
}

func IptablesGetRule(id string) (outRule *iptables.Rule, err error) {
	return iptables.GetRuleById(id)
}

func IptablesGetByTable(table string) (outRules []*iptables.Rule, err error) {
	return iptables.GetRulesByTable(table)
}

func IptablesGetByChain(table string, chain string) (outRules []*iptables.Rule, err error) {
	return iptables.GetRulesByChain(table, chain)
}

func IptablesUpdateRule(id string, inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	if _, err = iptables.FindRuleById(id); err == nil {
		r, err := iptables.GetRuleById(id)
		if err != nil {
			return nil, err
		}
		r.Update(inRule)
		err = r.Replace()
		return r, err
	} else {
		return nil, fmt.Errorf("failed to find rule matching the supplied id: %s", inRule.Id)
	}
}

func IptablesCreateRule(inRule *iptables.Rule) (outRule *iptables.Rule, err error) {
	id := uuid.New().String()
	inRule.Id = id
	inRule.SetApp(app)
	err = inRule.Append()
	if err != nil {
		return nil, err
	}
	return inRule, nil
}