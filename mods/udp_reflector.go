package mods

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type UdpReflector struct {
	port    int
	running bool
}

func (r *UdpReflector) Proto() string {
	return "udp"
}

func (r *UdpReflector) Port() int {
	return r.port
}

func (r *UdpReflector) SetPort(port int) {
	r.port = port
}

func (r *UdpReflector) Start() error {
	r.running = true
	go r.echo()
	return nil
}

func (r *UdpReflector) Stop() error {
	r.running = false //todo: this won't actually work as the thread will be blocked...
	return nil
}

func (r *UdpReflector) echo() {
	addr := net.UDPAddr{Port: r.port, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.WithFields(log.Fields{
			"addr": addr,
			"err":  err,
		}).Error("failed to start udp echo server")
	}

	b := make([]byte, 2048)

	for r.running {

		read, remote, rderr := conn.ReadFromUDP(b)
		if rderr != nil {
			log.WithFields(log.Fields{
				"remote": remote,
				"err":    rderr,
			}).Error("failed to read from udp socket")
		}

		_, wrerr := conn.WriteTo(b[0:read], remote)
		if wrerr != nil {
			log.WithFields(log.Fields{
				"remote": remote,
				"data":   b[:read],
				"err":    wrerr,
			}).Error("failed to write to udp socket")
		}
	}
}
