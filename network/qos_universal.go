package network

import (
	"fmt"
	"github.com/BGrewell/go-netqospolicy"
	"github.com/BGrewell/iptables"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type UniversalDSCPRule struct {
	Name       string             `json:"name"`
	DSCP       int                `json:"dscp"`
	Precedence *int                `json:"precedence,omitempty"`
	User       *string             `json:"user,omitempty"`
	Match      UniversalDSCPMatch `json:"match"`
}

type UniversalDSCPMatch struct {
	Protocol     *string `json:"protocol,omitempty"`
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
	r := &iptables.Rule{
		Chain: "OUTPUT",
		Target: iptables.TargetDSCP{Value: u.DSCP},
		Table: iptables.TableMangle,
		Name: u.Name,
		Id: uuid.New().String(),
		App: app,
	}

	if u.Precedence != nil {
		log.WithFields(log.Fields{
			"precedence": u.Precedence,
			"name": u.Name,
			"dscp": u.DSCP,
		}).Warn("precedence was used on Linux but has not yet been implemented")
	}

	if u.Match.SrcPort != nil {
		port := uint16(*u.Match.SrcPort)
		r.SourcePort = iptables.InvertableString{
			Value:    strconv.Itoa(int(port)),
			Inverted: false,
		}
	}

	if u.Match.DstPort != nil {
		port := uint16(*u.Match.DstPort)
		r.DestinationPort = iptables.InvertableString{
			Value:    strconv.Itoa(int(port)),
			Inverted: false,
		}
	}

	if u.Match.SrcPortStart != nil && u.Match.SrcPortEnd != nil {
		start := uint16(*u.Match.SrcPortStart)
		end := uint16(*u.Match.SrcPortEnd)
		r.SourcePort = iptables.InvertableString{
			Value:    fmt.Sprintf("%d:%d", start, end),
			Inverted: false,
		}
	}

	if u.Match.DstPortStart != nil && u.Match.DstPortEnd != nil {
		start := uint16(*u.Match.DstPortStart)
		end := uint16(*u.Match.DstPortEnd)
		r.DestinationPort = iptables.InvertableString{
			Value:    fmt.Sprintf("%d:%d", start, end),
			Inverted: false,
		}
	}

	if u.Match.Protocol != nil {
		r.Protocol = iptables.InvertableString{
			Value:    *u.Match.Protocol,
			Inverted: false,
		}
	} else {
		r.Protocol = iptables.InvertableString{
			Value:    "tcp",
			Inverted: false,
		}
		log.Warn("protocol not specified when creating qos rule. Assuming 'tcp' protocol")
	}

	if u.Match.SrcAddr != nil {
		r.Source = iptables.InvertableString{
			Value:    *u.Match.SrcAddr,
			Inverted: false,
		}
	}

	if u.Match.DstAddr != nil {
		r.Destination = iptables.InvertableString{
			Value:    *u.Match.DstAddr,
			Inverted: false,
		}
	}

	return r
}

func (u *UniversalDSCPRule) ToWindowsQos() *netqospolicy.NetQoSPolicy {

	dscp := int8(u.DSCP)
	precedence := uint32(127)
	if u.Precedence != nil {
		precedence = uint32(*u.Precedence)
	}

	p := &netqospolicy.NetQoSPolicy{
		Name: u.Name,
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

