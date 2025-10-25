package client

import (
	"context"
	"net"
	"net/http"
	"time"
)

func New() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "tcp4", addr) // 强制使用 IPV4, 避免 CERNET 网络环境下 IPV6 可以免认证访问互联网的问题
			},
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		},
	}
}
