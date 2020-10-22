package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// refer to https://medium.com/@masnun/panic-recovery-middleware-for-go-http-handlers-51147c941f9 and  http://www.golangtraining.in/lessons/middleware/recovering-from-panic.html
func PanicHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("forcing a panic")
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if er := recover(); er != nil {
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "Internal Server Error",
				})
				f1 := logrus.Fields{}
				f2 := AppendFields(r.Context(), f1)
				s := GetError(er)
				logrus.WithFields(f2).Error(s)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func GetError(er interface{}) string {
	switch x := er.(type) {
	case string:
		return er.(string)
	case error:
		err := x
		return err.Error()
	default:
		return fmt.Sprintf("%v", er)
	}
}
