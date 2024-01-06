package ipAddress

import (
	"net/http"
	"strings"
)

func ReadUserIP(r *http.Request) string {
	ipPort := r.RemoteAddr
	ip := ipPort
	if idx := strings.LastIndex(ipPort, ":"); idx != -1 {
		ip = ipPort[:idx]
	}
	return ip
}
