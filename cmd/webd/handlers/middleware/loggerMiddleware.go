package middleware

import (
	"net/http"
	"time"
)

type MessageType int

const (
	REQUEST = iota
	RESPONSE
)

type logEntry struct {
	TransactionId uint64
	Type          MessageType
	Time          time.Time
	Method        string
	URL           string
	Headers       string
	Body          string

	ServerIP string

	Latency time.Duration
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
