package log

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Masker interface {
	Mask(fieldName string, s string) string
}
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
	//logFields["user"] = "Test"
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
		_, err := httputil.DumpRequest(r, true)
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if err != nil {
			logrus.Error("Error reading body: %v", err)
		} else {
			logFields[c.Request] = buf.String()
			r.Body = ioutil.NopCloser(buf)
		}
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

func BuildContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if fieldConfig.Map != nil && len(*fieldConfig.Map) > 0 {
			var v interface{}
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				m, ok := v.(map[string]interface{})
				if !ok {
					next.ServeHTTP(w, r)
				} else {
					var ctx context.Context
					ctx = r.Context()
					for k, e := range *fieldConfig.Map {
						if strings.Index(e, ".") >= 0 {
							v3 := ValueOf(v, e)
							if v3 != nil {
								s3, ok3 := v3.(string)
								if ok3 {
									if len(s3) > 0 {
										ctx = context.WithValue(ctx, k, s3)
									}
								} else {
									ctx = context.WithValue(ctx, k, v3)
								}
							}
						} else {
							x, ok2 := m[e]
							if ok2 {
								s3, ok3 := x.(string)
								if ok3 {
									if len(s3) > 0 {
										ctx = context.WithValue(ctx, k, s3)
									}
								} else {
									ctx = context.WithValue(ctx, k, x)
								}
							}
						}
					}
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func Include(vs []string, v string) bool {
	for _, s := range vs {
		if v == s {
			return true
		}
	}
	return false
}
func ValueOf(m interface{}, path string) interface{} {
	arr := strings.Split(path, ".")
	i := 0
	var c interface{}
	c = m
	l1 := len(arr) - 1
	for i < len(arr) {
		key := arr[i]
		m2, ok := c.(map[string]interface{})
		if ok {
			c = m2[key]
		}
		if !ok || i >= l1 {
			return c
		}
		i++
	}
	return c
}
