package httpapi

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *statusWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: rw, status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(sw, r)

		dur := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sw.status, dur)
	})
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK) // status padrão antes de escrever o body
	}
	return w.ResponseWriter.Write(b)
}
