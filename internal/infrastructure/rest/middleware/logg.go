package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	CtxKeyStatus    = "status"
	CtxKeyErrorList = "error_list"
)

func LoggMiddleware(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Info(fmt.Sprintf("[START] %s %s | IP: %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		))

		next.ServeHTTP(w, r)

		ctx := r.Context()
		duration := time.Since(start)

		var statusDetails int = http.StatusOK
		if status := ctx.Value(CtxKeyStatus); status != nil {
			if s, ok := status.(int); ok {
				statusDetails = s
			}
		}
		errorList, _ := ctx.Value(CtxKeyErrorList).([]error)
		if len(errorList) > 0 {
			log.Error(fmt.Sprintf("[ERROR] %s %s | Status: %d | Duration: %v | ErrorList %v",
				r.Method,
				r.URL.Path,
				statusDetails,
				duration,
				errorList,
			))
		} else {
			log.Info(fmt.Sprintf("[ END ] %s %s | Status: %d | Duration: %v",
				r.Method,
				r.URL.Path,
				statusDetails,
				duration,
			))
		}
	}
}

func UpdateContext(ctx context.Context, r *http.Request, status int, errorList *[]error) {
	ctx = context.WithValue(ctx, CtxKeyStatus, status)
	ctx = context.WithValue(ctx, CtxKeyErrorList, *errorList)
	*r = *r.WithContext(ctx)
}
