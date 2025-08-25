package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stderr)
}

// レスポンス内容をキャプチャするラッパー
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// ミドルウェア本体
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
			body:           &bytes.Buffer{},
		}

		// リクエストボディを読み直せるように
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		next.ServeHTTP(lrw, r)

		log.Printf(
			"[%s] %s %s %d %s\nRequest: %s\nResponse: %s\n",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start),
			string(reqBody),
			lrw.body.String(),
		)
	})
}
