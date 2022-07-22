package utils

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

// ClientIP get client ip
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ResponseCapture capture response and get response info
type ResponseCapture struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	body        *bytes.Buffer
}

func NewResponseCapture(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{
		ResponseWriter: w,
		wroteHeader:    false,
		body:           new(bytes.Buffer),
	}
}

func (c ResponseCapture) Bytes() []byte {
	return c.body.Bytes()
}

func (c ResponseCapture) StatusCode() int {
	return c.status
}
