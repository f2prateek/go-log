package log

type Handler interface {
	Handle(*Entry)
}

type HandlerFunc func(*Entry)

func (f HandlerFunc) Handle(e *Entry) {
	f(e)
}
