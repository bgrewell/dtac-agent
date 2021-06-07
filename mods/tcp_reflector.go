package mods

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type TcpReflector struct {
	port    int
	running bool
}

func (r *TcpReflector) Proto() string {
	return "tcp"
}

func (r *TcpReflector) Port() int {
	return r.port
}

func (r *TcpReflector) SetPort(port int) {
	r.port = port
}

func (r *TcpReflector) Start() error {
	r.running = true
	go r.echo()
	return nil
}

func (r *TcpReflector) Stop() error {
	r.running = false //todo: this won't actually work as the thread will be blocked...
	return nil
}

func (r *TcpReflector) echo() {
	addr := net.TCPAddr{Port: r.port, IP: net.ParseIP("0.0.0.0")}
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.WithFields(log.Fields{
			"addr": addr,
			"err":  err,
		}).Error("failed to start tcp echo server")
	}
	defer l.Close()

	for r.running {
		conn, err := l.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"conn": conn,
				"err":  err,
			}).Error("failed to accept tcp client")
		}

		go r.handleRequest(conn)
	}
}

func (r *TcpReflector) handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		b := make([]byte, 1024)
		read, err := conn.Read(b)
		if err != nil {
			if err != io.EOF {
				log.WithFields(log.Fields{
					"conn": conn,
					"err":  err,
				}).Error("failed to read from tcp connection")
			}
			return
		}
		conn.Write(b[:read])
	}
}
