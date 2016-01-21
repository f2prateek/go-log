package log

type Interceptor interface {
	Intercept(*Entry) bool
}
