package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewResponseWriterWithLog(w http.ResponseWriter) ResponseWriterWithLog {
	return ResponseWriterWithLog{
		bodySize:   0,
		statusCode: 0,
		w:          w,
	}
}

type ResponseWriterWithLog struct {
	statusCode int
	bodySize   int
	w          http.ResponseWriter
}

func (rwl *ResponseWriterWithLog) Write(data []byte) (int, error) {
	size, err := rwl.w.Write(data)
	rwl.bodySize = size
	return size, err
}

func (rwl *ResponseWriterWithLog) Header() http.Header {
	return rwl.w.Header()
}

func (rwl *ResponseWriterWithLog) WriteHeader(statusCode int) {
	rwl.w.WriteHeader(statusCode)
	rwl.statusCode = statusCode
}

func (rwl *ResponseWriterWithLog) Status() int {
	return rwl.statusCode
}

func (rwl *ResponseWriterWithLog) BodySize() int {
	return rwl.bodySize
}

func WithLog(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wl := NewResponseWriterWithLog(w)
		defer func() {
			reqDuration := time.Since(startTime)
			log.Infof(" %s %s %dms %db status:%d",
				r.Method,
				r.RequestURI,
				reqDuration.Milliseconds(),
				wl.BodySize(),
				wl.Status())

		}()

		f(&wl, r)
	}
	return http.HandlerFunc(logFn)
}
