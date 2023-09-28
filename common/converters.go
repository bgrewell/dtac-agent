package common

import (
	"log"
	"time"
)

func DurationToSeconds(durStr string) int {
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		log.Fatalf("Failed to parse duration: %v", err)
	}
	return int(dur.Seconds())
}
