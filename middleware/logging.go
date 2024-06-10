package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LogEntry struct {
	Level    string            `json:"level"`
	Url      string            `json:"url"`
	Method   string            `json:"method"`
	Header   map[string]string `json:"header"`
	Duration string            `json:"duration"`
	Request  string            `json:"request"`
	Response string            `json:"response"`
}

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 获取请求头
		requestHeader := make(map[string]string)
		for key, values := range r.Header {
			// 请求头的值是一个字符串切片，取第一个值
			if len(values) > 0 {
				requestHeader[key] = values[0]
			}
		}
		// 捕获请求体
		requestBody, _ := io.ReadAll(r.Body)
		// 重新赋值请求体，因为它已经被读取了
		r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		// 记录请求开始
		requestMap := map[string]interface{}{
			"header": requestHeader,
			"body":   string(requestBody),
		}
		rec := &responseRecorder{
			ResponseWriter: w,
			body:           new(bytes.Buffer),
		}
		next.ServeHTTP(rec, r)
		// 记录请求完成
		duration := time.Since(start).Milliseconds()
		responseString := rec.body.String()
		// 解析responseString
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(responseString), &responseMap); err == nil {
			if body, ok := responseMap["body"].(map[string]interface{}); ok {
				modifiedResponseBody, _ := json.Marshal(body)
				responseString = string(modifiedResponseBody)
			}
			var bodyMap map[string]interface{}
			if err := json.Unmarshal([]byte(responseString), &bodyMap); err == nil {
				if data, ok := bodyMap["data"].(map[string]interface{}); ok {
					modifiedResponseBody, _ := json.Marshal(data)
					responseString = string(modifiedResponseBody)
				}
			}
		}
		logEntry := LogEntry{
			Level:    "INFO",
			Url:      r.URL.String(),
			Method:   r.Method,
			Header:   requestHeader,
			Duration: fmt.Sprintf("%d ms", duration),
			Request:  requestMap["body"].(string),
			Response: responseString,
		}

		// 记录日志
		logger.Info("Request processed",
			zap.String("url", logEntry.Url),
			zap.String("method", logEntry.Method),
			zap.Any("header", logEntry.Header),
			zap.String("duration", logEntry.Duration),
			zap.String("request", logEntry.Request),
			zap.String("response", logEntry.Response),
		)

	})
}
