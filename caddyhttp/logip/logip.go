package logip

import (
	"os"
	"path"
	"strings"
	"net"
	"net/http"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

type LogIP struct {
	Next     httpserver.Handler
	Path	 string
}

func checkIPAddress(ip string) bool {
    if net.ParseIP(ip) == nil {
        return false
    } else {
        return true
    }
}

func (l LogIP) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var ip string
	if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" {
		ips := strings.Split(fwdFor, ", ")
		ip = ips[0]
	} else {
		// Otherwise, get the client ip from the request remote address.
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	if !checkIPAddress(ip) {
		return l.Next.ServeHTTP(w, r)
	}

	fPath := path.Join(l.Path, ip)
	_, err := os.Stat(fPath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(fPath)
		}
	}
	return l.Next.ServeHTTP(w, r)
}

