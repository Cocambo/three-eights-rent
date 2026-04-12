package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(target *url.URL, publicPrefix, targetPrefix string) gin.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalHost := req.Host
		originalDirector(req)

		req.URL.Path = rewritePath(req.URL.Path, publicPrefix, targetPrefix)
		req.URL.RawPath = rewritePath(req.URL.RawPath, publicPrefix, targetPrefix)
		req.Host = target.Host

		if clientIP := clientIPFromRemoteAddr(req.RemoteAddr); clientIP != "" {
			appendForwardHeader(req.Header, "X-Forwarded-For", clientIP)
		}
		req.Header.Set("X-Forwarded-Host", originalHost)
		if req.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
		} else {
			req.Header.Set("X-Forwarded-Proto", "http")
		}
	}

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"service unavailable","code":"PROXY_ERROR"}`))
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func rewritePath(path, publicPrefix, targetPrefix string) string {
	if path == "" {
		return targetPrefix
	}

	if !strings.HasPrefix(path, publicPrefix) {
		return path
	}

	rest := strings.TrimPrefix(path, publicPrefix)
	if rest == "" {
		return targetPrefix
	}

	targetPrefix = strings.TrimSuffix(targetPrefix, "/")

	if !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
	}

	return targetPrefix + rest
}

func clientIPFromRemoteAddr(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}

	return remoteAddr
}

func appendForwardHeader(header http.Header, key, value string) {
	if existing := header.Get(key); existing != "" {
		header.Set(key, existing+", "+value)
		return
	}

	header.Set(key, value)
}
