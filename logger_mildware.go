package log

import (
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
)

func Standardize(config ChiLogConfig) ChiLogConfig {
	if len(config.Duration) == 0 {
		config.Duration = "duration"
	}
	return config
}
func InitializeFieldConfig(c ChiLogConfig) {
	if len(c.Duration) > 0 {
		fieldConfig.Duration = c.Duration
	} else {
		fieldConfig.Duration = "duration"
	}
	fieldConfig.Log = c.Log
	fieldConfig.Ip = c.Ip
	fieldConfig.Map = c.Map
	fieldConfig.Constants = c.Constants
	fieldConfig.FieldMap = c.FieldMap
	if len(c.Fields) > 0 {
		fields := strings.Split(c.Fields, ",")
		fieldConfig.Fields = &fields
	}
	if len(c.Masks) > 0 {
		fields := strings.Split(c.Masks, ",")
		fieldConfig.Masks = &fields
	}
	if len(c.Skips) > 0 {
		fields := strings.Split(c.Skips, ",")
		fieldConfig.Skips = &fields
	}
}
func Logger(c ChiLogConfig, logger *logrus.Logger, f Formatter) func(h http.Handler) http.Handler {
	InitializeFieldConfig(c)
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !fieldConfig.Log || !logrus.IsLevelEnabled(logrus.InfoLevel) || InSkipList(r, fieldConfig.Skips) {
				h.ServeHTTP(w, r)
			} else {
				dw := NewResponseWriter(w)
				ww := middleware.NewWrapResponseWriter(dw, r.ProtoMajor)
				startTime := time.Now()
				logFields := BuildLogFields(c, w, r)
				f.AppendFieldLog(logger, w, r, c, logFields)
				single := !c.Separate
				if r.Method == "GET" || r.Method == "DELETE" {
					single = true
				}
				f.LogRequest(logger, r, c, logFields, single)
				defer func() {
					if single {
						f.LogResponse(logger, w, r, ww, c, startTime, dw.Body.String(), logFields, single)
					} else {
						resLogFields := BuildLogFields(c, w, r)
						f.AppendFieldLog(logger, w, r, c, resLogFields)
						f.LogResponse(logger, w, r, ww, c, startTime, dw.Body.String(), resLogFields, single)
					}
				}()
				h.ServeHTTP(ww, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
func InSkipList(r *http.Request, skips *[]string) bool {
	if skips == nil || len(*skips) == 0 {
		return false
	}
	for _, s := range *skips {
		if strings.HasSuffix(s, r.RequestURI) {
			return true
		}
	}
	return false
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
	if len(c.Uri) > 0 {
		logFields[c.Uri] = r.RequestURI
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
		remoteIP := GetRemoteId(r)
		logFields[c.RemoteIp] = remoteIP
	}
	return logFields
}

func GetRemoteId(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
