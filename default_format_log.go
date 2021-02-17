package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Formatter interface {
	LogRequest(log func(context.Context, string, map[string]interface{}), r *http.Request, c LogConfig, logFields map[string]interface{}, singleLog bool)
	LogResponse(log func(context.Context, string, map[string]interface{}), w http.ResponseWriter, r *http.Request, ww WrapResponseWriter, c LogConfig, startTime time.Time, response string, logFields map[string]interface{}, singleLog bool)
}
type StructuredLogger struct {
	Produce    func(ctx context.Context, data []byte, attributes map[string]string) (string, error)
	KeyMap     map[string]string
	Goroutines bool
}

var fieldConfig FieldConfig

func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{}
}
func NewStructuredLoggerWithProduce(produce func(context.Context, []byte, map[string]string) (string, error), goroutines bool, options...map[string]string) *StructuredLogger {
	var keyMap map[string]string
	if len(options) >= 1 {
		keyMap = options[0]
	}
	return &StructuredLogger{Produce: produce, Goroutines: goroutines, KeyMap: keyMap}
}
func (l *StructuredLogger) LogResponse(log func(ctx context.Context, msg string, fields map[string]interface{}), w http.ResponseWriter, r *http.Request, ww WrapResponseWriter,
	c LogConfig, t1 time.Time, response string, logFields map[string]interface{}, singleLog bool) {
	fs := BuildResponseBody(ww, c, t1, response, logFields)
	var msg string
	if singleLog {
		msg = r.Method + " " + r.RequestURI
	} else {
		msg = "Response " + r.Method + " " + r.RequestURI
	}
	log(r.Context(), msg, fs)
	if l.Produce != nil {
		if l.Goroutines {
			go Produce(r.Context(), l.Produce, msg, logFields, l.KeyMap)
		} else {
			Produce(r.Context(), l.Produce, msg, logFields, l.KeyMap)
		}
	}
}
func Produce(ctx context.Context, produce func(ctx context.Context, data []byte, attributes map[string]string) (string, error), msg string, logFields map[string]interface{}, keyMap map[string]string) {
	m2 := AddKeyFields(msg, logFields, keyMap)
	b, err := json.Marshal(m2)
	if err == nil {
		produce(ctx, b, nil)
	}
}
func (l *StructuredLogger) LogRequest(log func(ctx context.Context, msg string, fields map[string]interface{}), r *http.Request, c LogConfig, logFields map[string]interface{}, singleLog bool) {
	var fs map[string]interface{}
	fs = logFields
	if len(c.Request) > 0 && r.Method != "GET" && r.Method != "DELETE" {
		fs = BuildRequestBody(r, c, logFields)
	}
	if !singleLog {
		msg := "Request " + r.Method + " " + r.RequestURI
		log(r.Context(), msg, fs)
		if l.Produce != nil {
			if l.Goroutines {
				go Produce(r.Context(), l.Produce, msg, logFields, l.KeyMap)
			} else {
				Produce(r.Context(), l.Produce, msg, logFields, l.KeyMap)
			}
		}
	}
}

func BuildResponseBody(ww WrapResponseWriter, c LogConfig, t1 time.Time, response string, logFields map[string]interface{}) map[string]interface{} {
	if len(c.Response) > 0 {
		logFields[c.Response] = response
	}
	if len(c.ResponseStatus) > 0 {
		logFields[c.ResponseStatus] = ww.Status()
	}
	if len(c.Duration) > 0 {
		t2 := time.Now()
		duration := t2.Sub(t1)
		logFields[c.Duration] = duration.Milliseconds()
	}
	if len(c.Size) > 0 {
		logFields[c.Size] = ww.BytesWritten()
	}
	return logFields
}
func BuildRequestBody(r *http.Request, c LogConfig, logFields map[string]interface{}) map[string]interface{} {
	if r.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		logFields[c.Request] = buf.String()
		r.Body = ioutil.NopCloser(buf)
	}
	return logFields
}
func AddKeyFields(message string, m map[string]interface{}, keys map[string]string) map[string]interface{} {
	level := "level"
	t := "time"
	msg := "msg"
	if keys != nil {
		ks := keys
		v1, ok1 := ks[level]
		if ok1 && len(v1) > 0 {
			level = v1
		}
		v2, ok2 := ks[t]
		if ok2 && len(v2) > 0 {
			t = v2
		}
		v3, ok3 := ks[msg]
		if ok3 && len(v3) > 0 {
			msg = v3
		}
	}
	m[msg] = message
	m[level] = "info"
	m[t] = time.Now()
	return m
}
