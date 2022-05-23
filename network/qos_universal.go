package network

import (
	"fmt"
	"github.com/BGrewell/go-iptables"
	"github.com/BGrewell/go-netqospolicy"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type UniversalDSCPRule struct {
	Id         string             `json:"id,omitempty"`
	Name       string             `json:"name"`
	DSCP       int                `json:"dscp"`
	Precedence *int               `json:"precedence,omitempty"`
	User       *string            `json:"user,omitempty"`
	Match      UniversalDSCPMatch `json:"match"`
}

type UniversalDSCPMatch struct {
	Protocol *string `json:"protocol,omitempty"`
	//Port         *int    `json:"port,omitempty"`
	SrcAddr      *string `json:"src_addr,omitempty"`
	SrcPort      *int    `json:"src_port,omitempty"`
	SrcPortStart *int    `json:"src_port_start,omitempty"`
	SrcPortEnd   *int    `json:"src_port_end,omitempty"`
	DstAddr      *string `json:"dst_addr,omitempty"`
	DstPort      *int    `json:"dst_port,omitempty"`
	DstPortStart *int    `json:"dst_port_start,omitempty"`
	DstPortEnd   *int    `json:"dst_port_end,omitempty"`
}

func (u *UniversalDSCPRule) ToLinuxQos() *iptables.Rule {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	r := &iptables.Rule{
		Chain:  "OUTPUT",
		Target: &iptables.TargetDSCP{Value: u.DSCP},
		Table:  iptables.TableMangle,
		Name:   u.Name,
		Id:     u.Id,
	}
	marker := &iptables.MarkerGeneric{}
	marker.SetName("app")
	marker.SetValue(app)
	r.AddMarker(marker)

	if u.Precedence != nil {
		log.WithFields(log.Fields{
			"precedence": u.Precedence,
			"name":       u.Name,
			"dscp":       u.DSCP,
		}).Warn("precedence was used on Linux but has not yet been implemented")
	}

	if u.Match.SrcPort != nil {
		port := uint16(*u.Match.SrcPort)
		r.SourcePort = strconv.Itoa(int(port))
		r.SourcePortNegated = false
	}

	if u.Match.DstPort != nil {
		port := uint16(*u.Match.DstPort)
		r.DestinationPort = strconv.Itoa(int(port))
		r.DestinationPortNegated = false
	}

	if u.Match.SrcPortStart != nil && u.Match.SrcPortEnd != nil {
		start := uint16(*u.Match.SrcPortStart)
		end := uint16(*u.Match.SrcPortEnd)
		r.SourcePort = fmt.Sprintf("%d:%d", start, end)
		r.SourcePortNegated = false
	}

	if u.Match.DstPortStart != nil && u.Match.DstPortEnd != nil {
		start := uint16(*u.Match.DstPortStart)
		end := uint16(*u.Match.DstPortEnd)
		r.DestinationPort = fmt.Sprintf("%d:%d", start, end)
		r.DestinationPortNegated = false
	}

	if u.Match.Protocol != nil {
		r.Protocol = iptables.Protocol(*u.Match.Protocol)
	} else {
		r.Protocol = "tcp"
		log.Warn("protocol not specified when creating qos rule. Assuming 'tcp' protocol")
	}

	if u.Match.SrcAddr != nil {
		r.Source = *u.Match.SrcAddr
	}

	if u.Match.DstAddr != nil {
		r.Destination = *u.Match.DstAddr
	}

	return r
}

func ConvertLinuxQosToUniversalDSCPRule(rule *iptables.Rule) *UniversalDSCPRule {

	m := UniversalDSCPMatch{}

	if rule.Protocol != "" {
		protocol := string(rule.Protocol)
		m.Protocol = &protocol
	}

	if rule.Source != "" {
		m.SrcAddr = &rule.Source
	}

	if rule.Destination != "" {
		m.DstAddr = &rule.Destination
	}

	if rule.SourcePort != "" {
		p := strings.Split(rule.SourcePort, ":")
		if len(p) == 1 {
			port, _ := strconv.Atoi(p[0])
			m.SrcPort = &port
		} else {
			start, _ := strconv.Atoi(p[0])
			end, _ := strconv.Atoi(p[1])
			m.SrcPortStart = &start
			m.SrcPortEnd = &end
		}
	}

	if rule.DestinationPort != "" {
		p := strings.Split(rule.DestinationPort, ":")
		if len(p) == 1 {
			port, _ := strconv.Atoi(p[0])
			m.DstPort = &port
		} else {
			start, _ := strconv.Atoi(p[0])
			end, _ := strconv.Atoi(p[1])
			m.DstPortStart = &start
			m.DstPortEnd = &end
		}
	}

	r := &UniversalDSCPRule{
		Id:         rule.Id,
		Name:       rule.Name,
		DSCP:       rule.Target.(*iptables.TargetDSCP).Value,
		Precedence: nil,
		User:       nil,
		Match:      m,
	}

	return r
}

func (u *UniversalDSCPRule) ToWindowsQos() *netqospolicy.NetQoSPolicy {

	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	dscp := int8(u.DSCP)
	precedence := uint32(127)
	if u.Precedence != nil {
		precedence = uint32(*u.Precedence)
	}

	p := &netqospolicy.NetQoSPolicy{
		Name:       u.Name,
		DSCPAction: &dscp,
		Precedence: &precedence,
		Persistent: false,
	}

	if u.User != nil {
		p.UserMatchCondition = *u.User
	}

	if u.Match.SrcPort != nil {
		port := uint16(*u.Match.SrcPort)
		p.IPSrcPortMatchCondition = &port
	}

	if u.Match.DstPort != nil {
		port := uint16(*u.Match.DstPort)
		p.IPDstPortMatchCondition = &port
	}

	if u.Match.SrcPortStart != nil && u.Match.SrcPortEnd != nil {
		start := uint16(*u.Match.SrcPortStart)
		end := uint16(*u.Match.SrcPortEnd)
		p.IPSrcPortStartMatchCondition = &start
		p.IPSrcPortEndMatchCondition = &end
	}

	if u.Match.DstPortStart != nil && u.Match.DstPortEnd != nil {
		start := uint16(*u.Match.DstPortStart)
		end := uint16(*u.Match.DstPortEnd)
		p.IPDstPortStartMatchCondition = &start
		p.IPDstPortEndMatchCondition = &end
	}

	if u.Match.Protocol != nil {
		p.IPProtocolMatchCondition = netqospolicy.IPProtocol(*u.Match.Protocol)
	}

	if u.Match.SrcAddr != nil {
		p.IPSrcPrefixMatchCondition = *u.Match.SrcAddr
	}

	if u.Match.DstAddr != nil {
		p.IPDstPrefixMatchCondition = *u.Match.DstAddr
	}

	return p
}

func ConvertWindowsQosToUniversalDSCPRule(rule *netqospolicy.NetQoSPolicy) *UniversalDSCPRule {

	if rule == nil {
		log.Warn("attempt to convert a nil rule. this shouldn't happen! returning nil")
		return nil
	}

	var precedence *int
	if rule.Precedence != nil {
		p := int(*rule.Precedence)
		precedence = &p
	}

	m := UniversalDSCPMatch{}

	if rule.IPSrcPortMatchCondition != nil {
		p := int(*rule.IPSrcPortMatchCondition)
		m.SrcPort = &p
	}

	if rule.IPDstPortMatchCondition != nil {
		p := int(*rule.IPDstPortMatchCondition)
		m.DstPort = &p
	}

	if rule.IPSrcPortStartMatchCondition != nil && rule.IPSrcPortEndMatchCondition != nil {
		s := int(*rule.IPSrcPortStartMatchCondition)
		e := int(*rule.IPSrcPortEndMatchCondition)
		m.SrcPortStart = &s
		m.SrcPortEnd = &e
	}

	if rule.IPDstPortStartMatchCondition != nil && rule.IPDstPortEndMatchCondition != nil {
		s := int(*rule.IPDstPortStartMatchCondition)
		e := int(*rule.IPDstPortEndMatchCondition)
		m.DstPortStart = &s
		m.DstPortEnd = &e
	}

	if rule.IPProtocolMatchCondition != netqospolicy.IPProto_BOTH {
		p := string(rule.IPProtocolMatchCondition)
		m.Protocol = &p
	}

	if rule.IPSrcPrefixMatchCondition != "" {
		m.SrcAddr = &rule.IPSrcPrefixMatchCondition
	}

	if rule.IPDstPrefixMatchCondition != "" {
		m.DstAddr = &rule.IPDstPrefixMatchCondition
	}

	u := &UniversalDSCPRule{
		Id:         rule.Name,
		Name:       rule.Name,
		DSCP: int(*rule.DSCPAction),
		Precedence: precedence,
		User:       nil,
		Match:      m,
	}

	if rule.UserMatchCondition != "" {
		u.User = &rule.UserMatchCondition
	}

	return u
}