package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-admin/internal/utils"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New()

		fields := log.Fields{
			"request_id": requestId,
			"method": r.Method,
			"url_path": r.URL.Path,
		}

		if log.IsLevelEnabled(log.DebugLevel) {
			fields["url_values"] = fmt.Sprintf("%v", mux.Vars(r))

			defer r.Body.Close()
			byteBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("Error reading request body: %s", err)
				utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
				return
			}
			fields["body"] = string(byteBody)
		}

		log.WithFields(fields).Info("New request")

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), utils.ContextKeyReqId, requestId.String())))
	})
}