// +build darwin

package network

func GetInterfaceStats(name string) (stats *InterfaceStats, err error) {
	return nil, fmt.Errorf("this function has not been implemented on this operating system yet")
}
