package network

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	iphlp                 = windows.NewLazySystemDLL("Iphlpapi.dll")
	fGetIpStatistics      = iphlp.NewProc("GetIpStatistics")
	fGetIpForwardTable    = iphlp.NewProc("GetIpForwardTable")
	fSetIpForwardEntry    = iphlp.NewProc("SetIpForwardEntry")
	fCreateIpForwardEntry = iphlp.NewProc("CreateIpForwardEntry")
	fDeleteIpForwardEntry = iphlp.NewProc("DeleteIpForwardEntry")
)

type ForwardProtocol uint32

const (
	MIB_IPPROTO_OTHER             ForwardProtocol = 1
	MIB_IPPROTO_LOCAL             ForwardProtocol = 2
	MIB_IPPROTO_NETMGMG           ForwardProtocol = 3
	MIB_IPPROTO_ICMP              ForwardProtocol = 4
	MIB_IPPROTO_EGP               ForwardProtocol = 5
	MIB_IPPROTO_GGP               ForwardProtocol = 6
	MIB_IPPROTO_HELLO             ForwardProtocol = 7
	MIB_IPPROTO_RIP               ForwardProtocol = 8
	MIB_IPPROTO_IS_IS             ForwardProtocol = 9
	MIB_IPPROTO_ES_IS             ForwardProtocol = 10
	MIB_IPPROTO_CISCO             ForwardProtocol = 11
	MIB_IPPROTO_BBN               ForwardProtocol = 12
	MIB_IPPROTO_OSPF              ForwardProtocol = 13
	MIB_IPPROTO_BGP               ForwardProtocol = 14
	MIB_IPPROTO_NT_AUTOSTATIC     ForwardProtocol = 10002
	MIB_IPPROTO_NT_STATIC         ForwardProtocol = 10006
	MIB_IPPROTO_NT_STATIC_NON_DOD ForwardProtocol = 10007
)

type ForwardType uint32

const (
	MIB_IPROUTE_TYPE_OTHER    ForwardType = 1
	MIB_IPROUTE_TYPE_INVALID  ForwardType = 2
	MIB_IPROUTE_TYPE_DIRECT   ForwardType = 3
	MIB_IPROUTE_TYPE_INDIRECT ForwardType = 4
)

// MIB_IPFORWARDROW structure
type MIB_IPFORWARDTABLE struct {
	dwNumEntries uint32
	table        [10000]MIB_IPFORWARDROW
}

// Size returns the size of the MIB_IPFORWARDTABLE structure
func (m MIB_IPFORWARDTABLE) Size() uint32 {
	return uint32(unsafe.Sizeof(m))
}

// MIB_IPFORWARDROW structure
type MIB_IPFORWARDROW struct {
	dwForwardDest      uint32
	dwForwardMask      uint32
	dwForwardPolicy    uint32
	dwForwardNextHop   uint32
	dwForwardIfIndex   uint32
	dwForwardType      uint32
	dwForwardProto     uint32
	dwForwardAge       uint32
	dwForwardNextHopAs uint32
	dwForwardMetric1   uint32
	dwForwardMetric2   uint32
	dwForwardMetric3   uint32
	dwForwardMetric4   uint32
	dwForwardMetric5   uint32
}

// Size returns the size of the MIB_IPFORWARDROW structure
func (m MIB_IPFORWARDROW) Size() uint32 {
	return uint32(unsafe.Sizeof(m))
}

// MIB_IPSTATS_LH structure
type MIB_IP_STATS_LH struct {
	dwForwarding      uint32
	dwDefaultTTL      uint32
	dwInReceives      uint32
	dwInHdrErrors     uint32
	dwInAddrErrors    uint32
	dwForwDatagrams   uint32
	dwInUnknownProtos uint32
	dwInDiscards      uint32
	dwInDelivers      uint32
	dwOutRequests     uint32
	dwRoutingDiscards uint32
	dwOutDiscards     uint32
	dwOutNoRoutes     uint32
	dwReasmTimeout    uint32
	dwReasmReqds      uint32
	dwReasmOks        uint32
	dwReasmFails      uint32
	dwFragOks         uint32
	dwFragFails       uint32
	dwFragCreates     uint32
	dwNumIf           uint32
	dwNumAddr         uint32
	dwNumRoutes       uint32
}

// Size returns the size of the MIB_IP_STATS_LH structure
func (m *MIB_IP_STATS_LH) Size() uint32 {
	return uint32(unsafe.Sizeof(m))
}
