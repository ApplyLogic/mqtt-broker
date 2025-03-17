package middleware

import (
	"github.com/ApplyLogic/mqtt-broker/config"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Logger struct {
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func (l *Logger) Initialize(cfg *config.Config) {

}

func (l *Logger) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		r.Header.Set("X-Request-ID", requestId)
		log.WithFields(log.Fields{
			"message":      "Request received",
			"method":       r.Method,
			"url":          r.URL.Path,
			"X-Request-ID": requestId,
		})
		next.ServeHTTP(w, r)
	})
}

func (l *Logger) Log(key string, value string) {
	log.WithFields(log.Fields{
		key: value,
	})
}
