package plugins

// LoggingLevel is a type for logging levels
type LoggingLevel int

const (
	// LevelDebug is the debug level
	LevelDebug LoggingLevel = iota
	// LevelInfo is the info level
	LevelInfo
	// LevelWarning is the warning level
	LevelWarning
	// LevelError is the error level
	LevelError
	// LevelFatal is the fatal level
	LevelFatal
)

// LogMessage is a struct for log messages
type LogMessage struct {
	Level   LoggingLevel
	Message string
	Fields  map[string]string
}
