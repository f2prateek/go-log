package log

import "sync"

type Builder struct {
	Interceptors []Interceptor
	Handlers     []Handler
	Fields       map[string]interface{}
}

func NewBuilder() *Builder {
	return &Builder{
		Interceptors: make([]Interceptor, 0),
		Handlers:     make([]Handler, 0),
		Fields:       make(map[string]interface{}),
	}
}

func (b *Builder) AddInterceptor(interceptor Interceptor) *Builder {
	b.Interceptors = append(b.Interceptors, interceptor)
	return b
}

func (b *Builder) AddHandler(handler Handler) *Builder {
	b.Handlers = append(b.Handlers, handler)
	return b
}

func (b *Builder) AddField(k string, v interface{}) *Builder {
	b.Fields[k] = v
	return b
}

func (b *Builder) Build() *Log {
	return &Log{
		interceptors: b.Interceptors,
		handlers:     b.Handlers,
		fields:       b.Fields,
	}
}

type Log struct {
	interceptors []Interceptor
	handlers     []Handler
	fields       map[string]interface{}
	sync.Mutex
}

func (l *Log) Debug(fields map[string]interface{}, msg string) {
	l.handle(Debug, msg, fields)
}

func (l *Log) Info(fields map[string]interface{}, msg string) {
	l.handle(Info, msg, fields)
}

func (l *Log) Warn(fields map[string]interface{}, msg string) {
	l.handle(Warn, msg, fields)
}

func (l *Log) Error(fields map[string]interface{}, err error) {
	l.handle(Error, err.Error(), fields)
}

func (l *Log) Fatal(fields map[string]interface{}, err error) {
	l.handle(Error, err.Error(), fields)
	panic(err)
}

func (l *Log) Clone() *Builder {
	return &Builder{
		Interceptors: append(make([]Interceptor, 0), l.interceptors...),
		Handlers:     append(make([]Handler, 0), l.handlers...),
		Fields:       merge(nil, l.fields),
	}
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
