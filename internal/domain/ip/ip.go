package ip

import "net/http"

func GetIp(req *http.Request) string {
	ip := req.Header.Get("X-FORWARDED-FOR")
	if ip != "" {
		return ip
	}
	return req.RemoteAddr
}