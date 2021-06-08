package mods

import (
	"bufio"
	"fmt"
	. "github.com/BGrewell/system-api/common"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type PingOverview struct {
	Results TimestampedFloatArray `json:"results"`
	Average float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

type UdpPingWorkerOptions struct {
	Target string `json:"target"`
	Port int `json:"port"`
	Interval int `json:"interval"`
	Timeout int `json:"timeout"`
}

type UdpPingWorker struct {
	Results TimestampedFloatArray
	target string
	port int
	interval int
	timeout int
	running bool
}

func (w *UdpPingWorker) SetOptions(options *UdpPingWorkerOptions) {
	w.target = options.Target
	w.port = options.Port
	w.interval = options.Interval
	w.timeout = options.Timeout
}

func (w *UdpPingWorker) Start() error {
	w.running = true
	go w.run()
	return nil
}

func (w *UdpPingWorker) Stop() error {
	w.running = false
	return nil
}

func (w *UdpPingWorker) run() {
	for w.running {
		next := time.Now().Add(time.Duration(w.interval) * time.Second)
		result, _ := UdpSendTimedPacket(w.target, w.port, w.timeout)
		w.Results.Add(result)
		for time.Now().Before(next) {
			time.Sleep(10 *time.Nanosecond)
		}
	}
}

func (w *UdpPingWorker) Average() float64 {
	return w.Results.Average()
}

func (w *UdpPingWorker) AveragePeriod(seconds int) float64 {
	return w.Results.AveragePeriod(seconds)
}

func (w *UdpPingWorker) StdDev() float64 {
	return w.Results.StdDev()
}

func (w *UdpPingWorker) StdDevPeriod(seconds int) float64 {
	return w.Results.StdDevPeriod(seconds)
}

func UdpSendTimedPacket(target string, port int, timeout int) (rtt float64, err error) {
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
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
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

	rtt = float64(recvtime - sendtime) / float64(time.Millisecond)
	log.Tracef("rtt was %f ms", rtt)
	return rtt, nil
}