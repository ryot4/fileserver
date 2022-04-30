package main

import (
	"log"
	"net/http"
)

type loggingWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *loggingWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.bytesWritten += n
	return n, err
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type LoggingHandler struct {
	http.Handler
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggingWriter{ResponseWriter: w, statusCode: http.StatusOK}
	h.Handler.ServeHTTP(lw, r)
	log.Printf("%s \"%s %s %s\" %d %d \"%s\"",
		r.RemoteAddr,
		r.Method,
		r.URL.String(),
		r.Proto,
		lw.statusCode,
		lw.bytesWritten,
		r.UserAgent())
}
