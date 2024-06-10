package middleware

import (
	"encoding/json"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = append(rw.body, data...)
	return len(data), nil
}

func JSONResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		// 包装成固定格式的JSON响应
		response := map[string]interface{}{
			"status": rw.statusCode,
			"body":   json.RawMessage(rw.body),
		}

		w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(rw.statusCode)
		json.NewEncoder(w).Encode(response)
	})
}
