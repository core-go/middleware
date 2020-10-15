package log

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Producer interface {
	Produce(ctx context.Context, data []byte, attributes *map[string]string) (string, error)
}
type Formatter interface {
	AppendFieldLog(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, c ChiLogConfig, logFields logrus.Fields)
	LogRequest(logger *logrus.Logger, r *http.Request, c ChiLogConfig, logFields logrus.Fields, singleLog bool)
	LogResponse(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, ww middleware.WrapResponseWriter, c ChiLogConfig, startTime time.Time, response string, logFields logrus.Fields, singleLog bool)
}

type StructuredLogger struct {
	Logger     *logrus.Logger
	Producer   Producer
	Goroutines bool
}

var loggerStatic = logrus.Logger{}
var fieldConfig FieldConfig

func NewStructuredLogger(logger *logrus.Logger) *StructuredLogger {
	loggerStatic = *logger
	return &StructuredLogger{Logger: logger}
}
func NewStructuredLoggerWithProducer(logger *logrus.Logger, producer Producer, goroutines bool) *StructuredLogger {
	loggerStatic = *logger
	return &StructuredLogger{Logger: logger, Producer: producer, Goroutines: goroutines}
}

func (l *StructuredLogger) LogResponse(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, ww middleware.WrapResponseWriter,
	c ChiLogConfig, t1 time.Time, response string, logFields logrus.Fields, singleLog bool) {
	BuildResponseBody(ww, c, t1, response, logFields)
	var msg string
	if singleLog {
		msg = r.Method + " " + r.RequestURI
	} else {
		msg = "Response " + r.Method + " " + r.RequestURI
	}
	logger.WithFields(logFields).Info(msg)
	if l.Producer != nil {
		if l.Goroutines {
			go Produce(r.Context(), l.Producer, msg, logFields, fieldConfig.KeyMap)
		} else {
			Produce(r.Context(), l.Producer, msg, logFields, fieldConfig.KeyMap)
		}
	}
}
func Produce(ctx context.Context, producer Producer, msg string, logFields logrus.Fields, keyMap *map[string]string) {
	m := BuildMap(logFields)
	m2 := AddKeyFields(msg, m, keyMap)
	b, err := json.Marshal(m2)
	if err == nil {
		producer.Produce(ctx, b, nil)
	}
}
func (l *StructuredLogger) LogRequest(logger *logrus.Logger, r *http.Request, c ChiLogConfig, logFields logrus.Fields, singleLog bool) {
	if len(c.Request) > 0 && r.Method != "GET" && r.Method != "DELETE" {
		BuildRequestBody(r, c, logFields)
	}
	if !singleLog {
		msg := "Request " + r.Method + " " + r.RequestURI
		logger.WithFields(logFields).Info(msg)
		if l.Producer != nil {
			if l.Goroutines {
				go Produce(r.Context(), l.Producer, msg, logFields, fieldConfig.KeyMap)
			} else {
				Produce(r.Context(), l.Producer, msg, logFields, fieldConfig.KeyMap)
			}
		}
	}
}

// Add more fields middleware request and response
func (l *StructuredLogger) AppendFieldLog(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, c ChiLogConfig, logFields logrus.Fields) {
	AppendFields(r.Context(), logFields)
}

func BuildResponseBody(ww middleware.WrapResponseWriter, c ChiLogConfig, t1 time.Time, response string, logFields logrus.Fields) {
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
}
func BuildRequestBody(r *http.Request, c ChiLogConfig, logFields logrus.Fields) {
	if r.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		logFields[c.Request] = buf.String()
		r.Body = ioutil.NopCloser(buf)
	}
}

func AppendFields(ctx context.Context, fields logrus.Fields) logrus.Fields {
	if len(fieldConfig.FieldMap) > 0 {
		if logFields, ok := ctx.Value(fieldConfig.FieldMap).(map[string]interface{}); ok {
			for k, v := range logFields {
				fields[k] = v
			}
		}
	}
	if fieldConfig.Fields != nil {
		cfs := *fieldConfig.Fields
		for _, k2 := range cfs {
			if v2, ok := ctx.Value(k2).(string); ok && len(v2) > 0 {
				fields[k2] = v2
			}
		}
	}
	return fields
}
func AddKeyFields(message string, m map[string]interface{}, keys *map[string]string) map[string]interface{} {
	level := "level"
	t := "time"
	msg := "msg"
	if keys != nil {
		ks := *keys
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
func BuildMap(fields logrus.Fields) map[string]interface{} {
	m := make(map[string]interface{})
	for k, e := range fields {
		m[k] = e
	}
	return m
}
