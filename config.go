package dominos

import (
	"os"
	"syscall"
)

type configuration struct {
	listeners []Listener
	signals   []os.Signal
	logger    logger
}

func New(options ...option) ListenCloser {
	var config configuration
	Options.apply(options...)(&config)
	return newSignalWatcher(newListener(config), config)
}

var Options singleton

type singleton struct{}
type option func(*configuration)

func (singleton) AddListeners(value ...Listener) option {
	return func(this *configuration) { this.listeners = append(this.listeners, value...) }
}
func (singleton) WatchTerminateSignals() option {
	return Options.WatchSignals(syscall.SIGINT, syscall.SIGTERM)
}
func (singleton) WatchSignals(value ...os.Signal) option {
	return func(this *configuration) { this.signals = append(this.signals, value...) }
}
func (singleton) Logger(value logger) option {
	return func(this *configuration) { this.logger = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, option := range Options.defaults(options...) {
			option(this)
		}

		if len(this.listeners) == 0 {
			this.listeners = append(this.listeners, nop{})
		}
	}
}
func (singleton) defaults(options ...option) []option {
	var defaultLogger = nop{}

	return append([]option{
		Options.Logger(defaultLogger),
	}, options...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type nop struct{}

func (nop) Printf(_ string, _ ...interface{}) {}
func (nop) Println(_ ...interface{})          {}

func (nop) Listen() {}