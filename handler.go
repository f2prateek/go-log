package log

type Handler interface {
	Handle(*Entry)
}
