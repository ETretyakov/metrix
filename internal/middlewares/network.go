package middlewares

import (
	"net"
	"net/http"
)

var subnet *net.IPNet

func InitSubnetMiddleware(s *net.IPNet) {
	subnet = s
}

func SubnetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subnet == nil {
			next.ServeHTTP(w, r)
		} else {
			reqIPAddr := r.Header.Get("X-Real-IP")
			if reqIPAddr == "" {
				http.Error(w, "empty x-real-ip header", http.StatusForbidden)
				return
			}

			ipAddr, _, err := net.ParseCIDR(reqIPAddr)
			if ipAddr == nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			if subnet.Contains(ipAddr) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "subnet not verified", http.StatusForbidden)
				return
			}
		}
	})
}
