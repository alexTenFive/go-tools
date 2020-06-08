package logger

import (
	"github.com/gocraft/health"
	"io"
	"sync"
)

var (
	stream      *health.Stream
	loggerLevel LogLevel
	mux sync.Mutex
	prefixes = map[LogLevel]string{
		Debug: "DEBUG",
		Info: "INFO",
		Warning: "WARNING",
		Error: "ERROR",
	}
)

type (
	//LogLevel - default log level for logger(can be debug, info, warn, error)
	LogLevel int
	//Logger - actually a logger
	Logger struct {
		job *health.Job
	}
)

const (
	// Debug - display all types of events
	Debug LogLevel = iota
	// Info - display all types of events except Debug
	Info
	// Warning - display all types of events except Info, Debug
	Warning
	// Error - display only Error events
	Error
)

// InitStream - create new stream function
func InitStream(w io.Writer) {
	if stream == nil {
		stream = health.NewStream()
		stream.AddSink(&health.WriterSink{Writer: w})
	}
}

// NewLogger - new instance of logger
func NewLogger(name string) *Logger {
	return &Logger{
		job: stream.NewJob(name),
	}
}

// SetLevel - set level of logger
func SetLevel(level LogLevel) {
	mux.Lock()
	defer mux.Unlock()
	loggerLevel = level
}

// Debug - debug job message
func (s *Logger) Debug(name string) {
	if loggerLevel != Debug {
		return
	}
	s.job.Event(getPrefix(Debug) + name)
}

// DebugKv - debug key-value job message
func (s *Logger) DebugKv(name string, kvs health.Kvs) {
	if loggerLevel != Debug {
		return
	}
	s.job.EventKv(getPrefix(Debug) + name, kvs)
}

// Info - information job message
func (s *Logger) Info(name string) {
	if loggerLevel <= Info {
		return
	}
	s.job.Event(getPrefix(Info) + name)
}

// InfoKv - information key-value job message
func (s *Logger) InfoKv(name string, kvs health.Kvs) {
	if loggerLevel <= Info {
		return
	}
	s.job.EventKv(getPrefix(Info) + name, kvs)
}

// Warn - warning job message
func (s *Logger) Warn(name string) {
	if loggerLevel <= Warning {
		return
	}
	s.job.Event(getPrefix(Warning) + name)
}

// WarnKv - warning key-value job message
func (s *Logger) WarnKv(name string, kvs health.Kvs) {
	if loggerLevel <= Warning {
		return
	}
	s.job.EventKv(getPrefix(Warning) + name, kvs)
}

// Err - error job message
func (s *Logger) Err(name string) {
	s.job.Event(getPrefix(Error) + name)
}

// ErrKv - error key-value message
func (s *Logger) ErrKv(name string, kvs health.Kvs) {
	s.job.EventKv(getPrefix(Error) + name, kvs)
}


func getPrefix(lvl LogLevel) string {
	return "[" + prefixes[lvl] + "] "
}