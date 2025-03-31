package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Printf("Error=%v Trace=%v", err, string(debug.Stack()))
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Printf("Status=%v Method=%v Path=%v Duration=%v", wrapped.Status(), r.Method, r.URL.EscapedPath(), time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{ResponseWriter: w}
}

func (w wrappedResponseWriter) Status() int {
	return w.status
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.status = code
	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
	return
}
