package mods

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func UdpSendTimedPacket(target string, port int) (rtt float32, err error) {
	buff := make([]byte, 1024)
	log.Tracef("creating udp connection to %s:%d", target, port)
	addr := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"target": target,
			"port": port,
			"error": err,
		}).Error("failed to dial udp connection")
		return -1, err
	}
	defer conn.Close()

	sendtime := time.Now().UnixNano()
	log.Tracef("sending %d to remote udp socket", sendtime)
	fmt.Fprintf(conn, "%d", sendtime)

	log.Tracef("reading from udp socket")
	read, err := bufio.NewReader(conn).Read(buff)
	recvtime := time.Now().UnixNano()
	if err != nil {
		log.WithFields(log.Fields{
			"target": target,
			"port": port,
			"error": err,
		}).Error("failed to read from udp connection")
		return -1, err
	}
	log.Tracef("read %d bytes from udp socket", read)

	rtt = float32(recvtime - sendtime) / float32(time.Millisecond)
	log.Tracef("rtt was %f ms", rtt)
	return rtt, nil
}