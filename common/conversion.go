package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

func ConvertStringToUInt64or0(uintStr string) uint64 {
	value, err := strconv.ParseUint(uintStr, 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"uintstr": uintStr,
			"value":   value,
			"err":     err,
		}).Debug("error converting string. returning default of 0")
		return 0
	}
	return value
}

func ConvertToRateMbps(lastBytes uint64, currentBytes uint64, lastTime int64, currentTime int64) float32 {
	log.Printf("last: %d current: %d", lastTime, currentTime)
	period := float32(currentTime - lastTime) / float32(int64(time.Second))
	change := (currentBytes - lastBytes) * 8
	log.Printf("change: %d", change)
	log.Printf("period: %.4f", period)
	rate := float32(change) / period
	mbps := rate / 1000 / 1000
	return mbps
}

func Inet_aton(ip string) (ip_int uint32) {
	ip_byte := net.ParseIP(ip).To4()
	for i := 0; i < len(ip_byte); i++ {
		ip_int |= uint32(ip_byte[i])
		if i < 3 {
			ip_int <<= 8
		}
	}
	return
}

func Inet_ntoa(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func Inet_ntoha(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
}