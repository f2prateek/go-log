package log

import (
	"errors"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

type RecorderHandler struct {
	entryC chan *Entry
}

func NewRecorderHandler() *RecorderHandler {
	return &RecorderHandler{
		entryC: make(chan *Entry, 1),
	}
}

func (r *RecorderHandler) Handle(e *Entry) {
	r.entryC <- e
}

type RecorderInterceptor struct {
	entryC chan *Entry
}

func NewRecorderInterceptor() *RecorderInterceptor {
	return &RecorderInterceptor{
		entryC: make(chan *Entry, 1),
	}
}

func (r *RecorderInterceptor) Intercept(e *Entry) bool {
	r.entryC <- e
	return true
}

func TestChain(t *testing.T) {
	clock = func() time.Time {
		return time.Unix(0, 0)
	}
	defer func() {
		clock = time.Now
	}()

	recorderHandler := NewRecorderHandler()
	recorderInterceptor := NewRecorderInterceptor()
	l := NewBuilder().AddInterceptor(recorderInterceptor).AddHandler(recorderHandler).Build()

	l.Error(map[string]interface{}{
		"key": "value",
	}, errors.New("error"))

	{
		e := <-recorderInterceptor.entryC
		assert.Equal(t, Error, e.Level)
		assert.Equal(t, "error", e.Message)
		assert.Equal(t, time.Unix(0, 0), e.Timestamp)
		assert.Equal(t, map[string]interface{}{
			"key": "value",
		}, e.Fields)
	}
	{
		e := <-recorderHandler.entryC
		assert.Equal(t, Error, e.Level)
		assert.Equal(t, "error", e.Message)
		assert.Equal(t, time.Unix(0, 0), e.Timestamp)
		assert.Equal(t, map[string]interface{}{
			"key": "value",
		}, e.Fields)
	}
}

func TestInterceptorFunc(t *testing.T) {
	entryC := make(chan *Entry, 1)
	f := func(e *Entry) bool {
		entryC <- e
		return true
	}

	InterceptorFunc(f).Intercept(&Entry{
		Level:   Debug,
		Message: "test",
	})

	e := <-entryC
	assert.Equal(t, Debug, e.Level)
	assert.Equal(t, "test", e.Message)
}

func TestHandlerFunc(t *testing.T) {
	entryC := make(chan *Entry, 1)
	f := func(e *Entry) {
		entryC <- e
	}

	HandlerFunc(f).Handle(&Entry{
		Level:   Debug,
		Message: "test",
	})

	e := <-entryC
	assert.Equal(t, Debug, e.Level)
	assert.Equal(t, "test", e.Message)
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		level LogLevel
		s     string
	}{
		{Debug, "debug"},
		{Info, "info"},
		{Warn, "warn"},
		{Error, "error"},
	}

	for _, c := range cases {
		assert.Equal(t, c.s, c.level.String())
	}
}

func TestInterceptorCanShortCircuit(t *testing.T) {
	recorderHandler := NewRecorderHandler()

	b := NewBuilder()
	b.AddInterceptor(InterceptorFunc(func(e *Entry) bool {
		return false
	}))
	b.AddHandler(recorderHandler)
	logger := b.Build()

	logger.Debug(nil, "foo")

	select {
	case <-recorderHandler.entryC:
		t.Error("handler should not have received any messages")
	default:
	}
}
