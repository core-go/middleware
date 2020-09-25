package log

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

func Standardize(config ChiLogConfig) ChiLogConfig {
	/*
		if len(config.Msg) == 0 {
			config.Msg = "msg"
		}
		if len(config.Timestamp) == 0 {
			config.Timestamp = "timestamp"
		}
		if len(config.TimeFormat) == 0 {
			config.TimeFormat = "2006-01-02T15:04:05-0700"
		}
	*/
	if len(config.Duration) == 0 {
		config.Duration = "duration"
	}
	return config
}

func Logger(c ChiLogConfig, logger *logrus.Logger, f Formatter) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dw := NewResponseWriter(w)
			ww := middleware.NewWrapResponseWriter(dw, r.ProtoMajor)
			startTime := time.Now()
			logFields := BuildLogFields(c, w, r)
			f.AppendFieldLog(logger, w, r, c, logFields)
			ctx := context.WithValue(r.Context(), Fields, logFields)
			if logrus.IsLevelEnabled(logrus.InfoLevel) {
				f.LogRequest(logger, r, c, logFields, c.Single)
				defer func() {
					if c.Single {
						f.LogResponse(logger, w, r, ww, c, startTime, dw.Body.String(), logFields, c.Single)
					} else {
						resLogFields := BuildLogFields(c, w, r)
						f.AppendFieldLog(logger, w, r, c, resLogFields)
						f.LogResponse(logger, w, r, ww, c, startTime, dw.Body.String(), resLogFields, c.Single)
					}
				}()
			}
			h.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
func BuildLogFields(c ChiLogConfig, w http.ResponseWriter, r *http.Request) logrus.Fields {
	logFields := logrus.Fields{}
	if !c.Build {
		return logFields
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	// logFields[c.Timestamp] = time.Now().UTC().Format(c.TimeFormat)
	if len(c.Uri) > 0 {
		logFields[c.Uri] = fmt.Sprintf("%s", r.RequestURI)
	}

	if len(c.ReqId) > 0 {
		if reqID := middleware.GetReqID(r.Context()); reqID != "" {
			logFields[c.ReqId] = reqID
		}
	}
	if len(c.Scheme) > 0 {
		logFields[c.Scheme] = scheme
	}
	if len(c.Proto) > 0 {
		logFields[c.Proto] = r.Proto
	}
	if len(c.UserAgent) > 0 {
		logFields[c.UserAgent] = r.UserAgent()
	}
	if len(c.RemoteAddr) > 0 {
		logFields[c.RemoteAddr] = r.RemoteAddr
	}
	if len(c.Method) > 0 {
		logFields[c.Method] = r.Method
	}
	if len(c.RemoteIp) > 0 {
		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteIP = r.RemoteAddr
		}
		logFields[c.RemoteIp] = remoteIP
	}
	return logFields
}
