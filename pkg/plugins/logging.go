package plugins

type LoggingLevel int

const (
	LevelDebug LoggingLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

type LogMessage struct {
	Level   LoggingLevel
	Message string
	Fields  map[string]string
}
