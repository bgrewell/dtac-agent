package webhost

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// ProxyHandler creates an HTTP handler that reverse-proxies requests to the target URL.
// It optionally injects JWT tokens obtained via the TokenProvider.
func ProxyHandler(config ProxyConfig, tokenProvider TokenProvider, logger Logger) http.HandlerFunc {
	// Parse the target URL
	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		logger.Error("Failed to parse proxy target URL", map[string]string{
			"proxy":  config.ID,
			"target": config.TargetURL,
			"error":  err.Error(),
		})
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Proxy configuration error", http.StatusInternalServerError)
		}
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the director to handle path stripping and header injection
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Strip prefix if configured
		if config.StripPrefix && config.FromPath != "/" {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, config.FromPath)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}

		// Inject additional headers
		for key, value := range config.Headers {
			req.Header.Set(key, value)
		}

		// Handle authentication
		if config.AuthType == "jwt_from_dtac" && tokenProvider != nil {
			// Request token from DTAC
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			token, _, err := tokenProvider.GetToken(ctx, config.Scopes, config.ID)
			if err != nil {
				logger.Warn("Failed to obtain token for proxy", map[string]string{
					"proxy": config.ID,
					"error": err.Error(),
				})
				// Continue without token - the upstream service will reject if auth is required
			} else {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}
		}

		// Preserve original host header or use target host
		req.Host = targetURL.Host
	}

	// Customize error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("Proxy error", map[string]string{
			"proxy":  config.ID,
			"path":   r.URL.Path,
			"error":  err.Error(),
			"remote": r.RemoteAddr,
		})
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Proxying request", map[string]string{
			"proxy":  config.ID,
			"path":   r.URL.Path,
			"method": r.Method,
			"remote": r.RemoteAddr,
		})
		proxy.ServeHTTP(w, r)
	}
}

// SetupProxies configures all proxy routes on the given HTTP mux.
func SetupProxies(mux *http.ServeMux, proxies []ProxyConfig, tokenProvider TokenProvider, logger Logger) {
	for _, proxyConfig := range proxies {
		handler := ProxyHandler(proxyConfig, tokenProvider, logger)
		mux.Handle(proxyConfig.FromPath, handler)

		logger.Info("Registered proxy route", map[string]string{
			"proxy":  proxyConfig.ID,
			"from":   proxyConfig.FromPath,
			"target": proxyConfig.TargetURL,
			"auth":   proxyConfig.AuthType,
		})
	}
}

// ProxyMiddleware wraps an HTTP handler with logging and error handling.
func ProxyMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Log incoming request
			logger.Debug("Incoming request", map[string]string{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			})

			// Call the next handler
			next.ServeHTTP(w, r)

			// Log completion
			duration := time.Since(start)
			logger.Debug("Request completed", map[string]string{
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": duration.String(),
			})
		})
	}
}
