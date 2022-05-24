package mods

import (
	"bufio"
	"fmt"
	. "github.com/BGrewell/dtac-agent/common"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"time"
)

type TcpPingWorkerOptions struct {
	Target      string `json:"target"`
	Port        int    `json:"port"`
	Interval    int    `json:"interval_ms"`
	Timeout     int    `json:"timeout"`
	PayloadSize int    `json:"payload_size"`
}

type TcpPingWorker struct {
	Results   TimestampedFloatArray
	target    string
	port      int
	interval  int
	timeout   int
	size      int
	running   bool
	randbytes []byte
}

func (w *TcpPingWorker) SetOptions(options *TcpPingWorkerOptions) {
	w.target = options.Target
	w.port = options.Port
	w.interval = options.Interval
	w.timeout = options.Timeout
	w.size = options.PayloadSize
}

func (w *TcpPingWorker) Start() error {
	w.randbytes = make([]byte, 1000000)
	rand.Read(w.randbytes)
	w.running = true
	go w.run()
	return nil
}

func (w *TcpPingWorker) Stop() error {
	w.running = false
	return nil
}

func (w *TcpPingWorker) run() {
	for w.running {
		next := time.Now().Add(time.Duration(w.interval) * time.Second)
		go w.TcpSendTimedPacket()
		for time.Now().Before(next) {
			time.Sleep(10 * time.Nanosecond)
		}
	}
}

func (w *TcpPingWorker) Average() float64 {
	return w.Results.Average()
}

func (w *TcpPingWorker) AveragePeriod(seconds int) float64 {
	return w.Results.AveragePeriod(seconds)
}

func (w *TcpPingWorker) StdDev() float64 {
	return w.Results.StdDev()
}

func (w *TcpPingWorker) StdDevPeriod(seconds int) float64 {
	return w.Results.StdDevPeriod(seconds)
}

func (w *TcpPingWorker) TcpSendTimedPacket() {
	result, _ := TcpSendTimedPacket(w.target, w.port, w.timeout, w.size, &w.randbytes)
	w.Results.Add(result)
}

func TcpSendTimedPacket(target string, port int, timeout int, size int, r *[]byte) (rtt float64, err error) {
	buff := make([]byte, 1024)
	log.Tracef("creating tcp connection to %s:%d", target, port)
	addr := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"target": target,
			"port":   port,
			"error":  err,
		}).Error("failed to dial tcp connection")
		return -1, err
	}
	defer conn.Close()

	idx := rand.Intn(len(*r) - size - 1)
	sendtime := time.Now().UnixNano()
	log.Tracef("sending %d to remote udp socket", size)
	conn.Write((*r)[idx : idx+size])

	log.Tracef("reading from tcp socket")
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	read, err := bufio.NewReader(conn).Read(buff)
	recvtime := time.Now().UnixNano()
	if err != nil {
		log.WithFields(log.Fields{
			"target": target,
			"port":   port,
			"error":  err,
		}).Error("failed to read from tcp connection")
		return -1, err
	}
	log.Tracef("read %d bytes from tcp socket", read)

	rtt = float64(recvtime-sendtime) / float64(time.Millisecond)
	log.Tracef("rtt was %f ms", rtt)
	return rtt, nil
}
