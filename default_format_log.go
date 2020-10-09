package log

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Formatter interface {
	AppendFieldLog(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, c ChiLogConfig, logFields logrus.Fields)
	LogRequest(logger *logrus.Logger, r *http.Request, c ChiLogConfig, logFields logrus.Fields, singleLog bool)
	LogResponse(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, ww middleware.WrapResponseWriter, c ChiLogConfig, startTime time.Time, response string, logFields logrus.Fields, singleLog bool)
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

var loggerStatic = logrus.Logger{}
var fieldConfig FieldConfig

func NewStructuredLogger(logger *logrus.Logger) *StructuredLogger {
	loggerStatic = *logger
	return &StructuredLogger{logger}
}

func (l *StructuredLogger) LogResponse(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, ww middleware.WrapResponseWriter,
	c ChiLogConfig, t1 time.Time, response string, logFields logrus.Fields, singleLog bool) {
	BuildResponseBody(ww, c, t1, response, logFields)
	if singleLog {
		logger.WithFields(logFields).Info(r.Method + " " + r.RequestURI)
	} else {
		logger.WithFields(logFields).Info("Response " + r.Method + " " + r.RequestURI)
	}
}

func (l *StructuredLogger) LogRequest(logger *logrus.Logger, r *http.Request, c ChiLogConfig, logFields logrus.Fields, singleLog bool) {
	if len(c.Request) > 0 && r.Method != "GET" && r.Method != "DELETE" {
		BuildRequestBody(r, c, logFields)
	}
	if !singleLog {
		logger.WithFields(logFields).Infof("Request %s %s", r.Method, r.RequestURI)
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
