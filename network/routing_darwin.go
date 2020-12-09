package network

type ForwardType uint32
type ForwardProtocol uint32

// GetRouteTable retrieves the full route table on the system
func GetRouteTable() (routes []RouteTableRow, err error) {
	return nil, fmt.Errorf("this method has not been implemented for this OS")
}

// UpdateRoute updates a given route on the system
func UpdateRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}

// CreateRoute creates a new route on the system
func CreateRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}

// DeleteRoute removes a route from the system
func DeleteRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}