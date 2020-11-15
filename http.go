package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HttpProducer struct {
	Client     *http.Client
	Url        string
	LogError   func(context.Context, string)
	Goroutines bool
	Retries    []time.Duration
}

func NewHttpProducer(client *http.Client, url string, logError func(context.Context, string), goroutines bool, retries ...time.Duration) *HttpProducer {
	return &HttpProducer{Client: client, Url: url, LogError: logError, Goroutines: goroutines, Retries: retries}
}
func (s *HttpProducer) Produce(ctx context.Context, data []byte, attributes *map[string]string) (string, error) {
	if s.Goroutines {
		go PostLog(ctx, s.Client, s.Url, data, nil, s.LogError, s.Retries...)
		return "", nil
	} else {
		err := PostLog(ctx, s.Client, s.Url, data, nil, s.LogError, s.Retries...)
		return "", err
	}
}
func PostLog(ctx context.Context, client *http.Client, url string, log []byte, headers *map[string]string, logError func(context.Context, string), retries ...time.Duration) error {
	l := len(retries)
	if l == 0 {
		_, err := Post(ctx, client, url, log, headers)
		return err
	} else {
		return PostWithRetries(ctx, client, url, log, headers, logError, retries)
	}
}
func PostWithRetries(ctx context.Context, client *http.Client, url string, log []byte, headers *map[string]string, logError func(context.Context, string), retries []time.Duration) error {
	_, er1 := Post(ctx, client, url, log, headers)
	if er1 == nil {
		return er1
	}
	i := 0
	err := Retry(ctx, retries, func() (err error) {
		i = i + 1
		_, er2 := Post(ctx, client, url, log, headers)
		s := string(log)
		if logError != nil {
			if er2 != nil {
				s2 := fmt.Sprintf("Fail to send log after %d retries %s", i, s)
				logError(ctx, s2)
			} else {
				s2 := fmt.Sprintf("Send log successfully after %d retries %s", i, s)
				logError(ctx, s2)
			}
		}
		return er2
	})
	if err != nil {
		if logError != nil {
			s := string(log)
			s2 := fmt.Sprintf("Failed to send log: %s. Error: %v.", s, err)
			logError(ctx, s2)
		}
	}
	return err
}
func GetString(ctx context.Context, key string) string {
	if len(key) > 0 {
		u := ctx.Value(key)
		if u != nil {
			s, ok := u.(string)
			if ok {
				return s
			} else {
				return ""
			}
		}
	}
	return ""
}
func Post(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*json.Decoder, error) {
	res, er1 := Do(ctx, client, url, "POST", body, headers)
	if er1 != nil {
		return nil, er1
	}
	if res.StatusCode == 503 {
		er2 := errors.New("503 Service Unavailable")
		return nil, er2
	}
	return json.NewDecoder(res.Body), nil
}
func Do(ctx context.Context, client *http.Client, url string, method string, body []byte, headers *map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return AddHeaderAndDo(client, req, headers)
}
func AddHeaderAndDo(client *http.Client, req *http.Request, headers *map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range *headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}

//Copy this code from https://stackoverflow.com/questions/47606761/repeat-code-if-an-error-occured
func Retry(ctx context.Context, sleeps []time.Duration, f func() error) (err error) {
	attempts := len(sleeps)
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleeps[i])
		//Infof(ctx, "Retrying %d of %d after error: %s", i+1, attempts, err.Error())
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
