package utilities

import (
	"fmt"
	"time"
)

func ConvertEpochTimeToTimestamp(epochTime int64) (timestamp string) {
	t := time.Unix(epochTime, 0)
	return t.Format("2006-01-02 15:04:05")
}

func ConvertBytesToHumanReadable(bytes int64) (humanReadableBytes string) {
	const (
		KB int64 = 1 << 10 // 1024
		MB       = KB << 10
		GB       = MB << 10
		TB       = GB << 10
		PB       = TB << 10
	)

	var unit string
	var value float64

	switch {
	case bytes >= PB:
		unit = "PB"
		value = float64(bytes) / float64(PB)
	case bytes >= TB:
		unit = "TB"
		value = float64(bytes) / float64(TB)
	case bytes >= GB:
		unit = "GB"
		value = float64(bytes) / float64(GB)
	case bytes >= MB:
		unit = "MB"
		value = float64(bytes) / float64(MB)
	case bytes >= KB:
		unit = "KB"
		value = float64(bytes) / float64(KB)
	default:
		unit = "B"
		value = float64(bytes)
	}

	return fmt.Sprintf("%.2f%s", value, unit)
}
