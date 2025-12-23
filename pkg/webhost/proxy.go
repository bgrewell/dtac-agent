package webhost

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxyHandler(target string, auth ProxyAuth, tp TokenProvider) (http.Handler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
	}

	if strings.EqualFold(auth.Type, "jwt_from_dtac") && tp != nil {
		proxy.Transport = &tokenRoundTripper{
			underlying: http.DefaultTransport,
			tp:         tp,
			scopes:     auth.Scopes,
		}
	}

	return proxy, nil
}

type tokenRoundTripper struct {
	underlying http.RoundTripper
	tp         TokenProvider
	scopes     []string
}

func (t *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, _, err := t.tp.GetToken(req.Context(), t.scopes, "")
	if err == nil && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return t.underlying.RoundTrip(req)
}
