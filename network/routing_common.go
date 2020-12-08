package network

type RouteTableRow struct {
	Destination    string          `json:"destination"`
	Mask           string          `json:"mask"`
	NextHop        string          `json:"next_hop"`
	Policy         uint32          `json:"policy"`
	InterfaceIndex uint32          `json:"interface_index"`
	Type           ForwardType     `json:"type"`
	Protocol       ForwardProtocol `json:"protocol"`
	Age            uint32          `json:"age"`
	NextHopAs      uint32          `json:"next_hop_as"`
	Metric1        uint32          `json:"metric_1"`
	Metric2        uint32          `json:"metric_2"`
	Metric3        uint32          `json:"metric_3"`
	Metric4        uint32          `json:"metric_4"`
	Metric5        uint32          `json:"metric_5"`
}
