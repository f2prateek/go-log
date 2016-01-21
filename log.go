package log

import (
	"fmt"
	"sync"
)

type Log struct {
	interceptors []Interceptor
	handlers     []Handler
	fields       map[string]interface{}
	sync.Mutex
}

func (l *Log) Debug(fields map[string]interface{}, format string, a ...interface{}) {
	l.handle(Debug, fmt.Sprintf(format, a), fields)
}

func (l *Log) Info(fields map[string]interface{}, format string, a ...interface{}) {
	l.handle(Info, fmt.Sprintf(format, a), fields)
}

func (l *Log) Warn(fields map[string]interface{}, format string, a ...interface{}) {
	l.handle(Warn, fmt.Sprintf(format, a), fields)
}

func (l *Log) Error(fields map[string]interface{}, format string, a ...interface{}) {
	l.handle(Error, fmt.Sprintf(format, a), fields)
}

// merge returns a new map with elements from each of the provided maps.
func merge(a, b map[string]interface{}) map[string]interface{} {
	c := make(map[string]interface{})
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func (l *Log) handle(level LogLevel, msg string, fields map[string]interface{}) {
	l.Lock()
	defer l.Unlock()

	e := &Entry{
		Level:   level,
		Message: msg,
		Fields:  merge(l.fields, fields),
	}

	for _, interceptor := range l.interceptors {
		if stop := interceptor.Intercept(e); stop {
			return
		}
	}

	for _, handler := range l.handlers {
		handler.Handle(e)
	}
}
