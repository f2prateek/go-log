package log

type Interceptor interface {
	Intercept(*Entry) bool
}

type InterceptorFunc func(*Entry) bool

func (f InterceptorFunc) Intercept(e *Entry) bool {
	return f(e)
}
