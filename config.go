package logging

// Config defines the configuration of go-logging.
type Config struct {
	TimeFormat string
	Handlers   []*Handler
}

// LoadConfig load configuration and setting up the go-logging.
func LoadConfig(config Config) {
	newHandlers := make(map[string][]*Handler)
	for _, handler := range config.Handlers {
		for _, level := range handler.Levels {
			if newHandlers[level] == nil {
				newHandlers[level] = make([]*Handler, 0, 4)
			}
			newHandlers[level] = append(newHandlers[level], handler)
		}
	}
	handlers = newHandlers
	if config.TimeFormat != "" {
		timeFormat = config.TimeFormat
	}
}
