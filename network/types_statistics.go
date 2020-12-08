package network

func init() {
	GetWindowsNetworkStats() //TODO Remove after testing
}

type IStatistics interface {
	Parse(json string) error
	JSON() (string, error)
	Update() error
}

type WindowsNetworkStats struct {
	StatsType                           string
	IpForwardingEnabled                 bool
	IpDefaultTTL                        int
	IpDatagramsReceived                 int
	IpDatagramsReceivedHeaderErrors     int
	IpDatagramsReceivedAddressErrors    int
	IpDatagramsForwarded                int
	IpDatagramsUnknownProtocol          int
	IpDatagramsReceivedDiscarded        int
	IpDatagramsDelivered                int
	IpDatagramsSent                     int
	IpDatagramsRoutingDiscards          int
	IpDatagramsSentDiscarded            int
	IpDatagramsSentNoRouteDiscarded     int
	IpDatagramFragmentReassemblyTimeout int
	IpDatagramsRequiredReassembled      int
	IpDatagramsReassembledOk            int
	IpDatagramsReassembledFail          int
	IpDatagramsFragmentOk               int
	IpDatagramsFragmentFail             int
	IpDatagramsFragmentCreated          int
	IpInterfaceCount                    int
	IpAddressCount                      int
	IpRouteCount                        int
	TcpRTOAlgo                          int
	TcpRTOMinValue                      int
	TcpRTOMaxValue                      int
	TcpConnMax                          int
	TcpConnActiveOpens                  int
	TcpConnPassiveOpens                 int
	TcpConnFailed                       int
	TcpConnEstablishedReset             int
	TcpConnCurrentEstablished           int
	TcpSegmentsIn                       int
	TcpSegmentsOut                      int
	TcpSegmentsRetrans                  int
	TcpSegmentsInErrors                 int
	TcpSegmentsOutReset                 int
	TcpConnCurrent                      int
}

func (s *WindowsNetworkStats) Parse(json string) (err error) {
	return nil
}

func (s *WindowsNetworkStats) JSON() (json string, err error) {
	return "", nil
}

func (s *WindowsNetworkStats) Update() (err error) {
	// Windows Function GetIpStatistics() should populate all Ip* fields
	// Windows Function GetTcpStatisticsEx2() should get all the Tcp* fields
	return GetIpStatistics()
}

func GetWindowsNetworkStats() (stats IStatistics, err error) {
	stats = &WindowsNetworkStats{
		StatsType: "WindowsNetworkStatistics",
	}
	err = stats.Update()
	return stats, err
}
