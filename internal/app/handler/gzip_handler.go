package handler

import (
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WithGzip(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get(ContentEncodingHeader) {
		case "gzip":
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Error("Error create gzip reader: ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			//Необходимо сохранить контекст запроса
			req, err := http.NewRequestWithContext(r.Context(), r.Method, r.RequestURI, reader)
			if err != nil {
				log.Error("Error create decompressed Request: ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			req.URL = r.URL
			req.Proto = r.Proto
			req.ProtoMajor = r.ProtoMajor
			req.ProtoMinor = r.ProtoMinor
			req.Header = r.Header
			req.ContentLength = r.ContentLength
			req.Host = r.Host
			req.PostForm = r.PostForm
			req.RequestURI = r.RequestURI
			req.Close = r.Close

			h.ServeHTTP(w, req)
		default:
			h.ServeHTTP(w, r)
		}

	}
	return http.HandlerFunc(logFn)
}
