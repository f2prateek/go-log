package log

import (
	"fmt"
	"time"
)

type Entry struct {
	Level     LogLevel
	Message   string
	Timestamp time.Time
	Fields    map[string]interface{}
}

func (e *Entry) String() string {
	return fmt.Sprintf("timestamp=%v level=%v message=%v %v", e.Timestamp, e.Level, e.Message, e.Fields)
}
