package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextKeyReqId ContextKey = "requestId"

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New()

		log.WithFields(log.Fields{
			"request_id": requestId,
			"method": r.Method,
			"url_path": r.URL.Path,
		}).Info("New request")

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextKeyReqId, requestId.String())))
	})
}