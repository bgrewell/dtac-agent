package plugins

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Register(args RegisterArgs, reply *RegisterReply) error
	RootPath() string
}
