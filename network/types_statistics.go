package network

type IStatistics interface {
	Parse(json string) error
	JSON() (string, error)
	Update() error
}
