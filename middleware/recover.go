package middleware

import (
	"go_code/simplek8s/internal/utils"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func RecoverMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.ByteString("stack", debug.Stack()),
				)
				// 返回 500 错误码
				utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
