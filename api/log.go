package api

import (
	"log"
	"net/http"
	"time"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lw := &logWriter{ResponseWriter: w}

		next.ServeHTTP(lw, r)

		duration := time.Since(startTime)

		log.Printf("%s %s %d %s %s",
			r.Method,
			r.RequestURI,
			lw.statusCode,
			r.RemoteAddr,
			duration,
		)
	})
}

type logWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *logWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}
