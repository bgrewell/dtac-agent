package common

import (
	log "github.com/sirupsen/logrus"
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