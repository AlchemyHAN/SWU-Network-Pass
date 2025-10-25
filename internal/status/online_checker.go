package status

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func VerifyOnline(ctx context.Context, logger *slog.Logger, client *http.Client, verificationURL *url.URL) (bool, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, verificationURL.String(), nil)
	if err != nil {
		return false, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return true, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, "", err
	}

	bodyStr := string(body)

	// 检查 HTML 中是否包含 <script>，代表被重定向到登录页
	if strings.Contains(bodyStr, "<script>") {
		logger.Warn("Network appears disconnected; Portal authentication required")
		return true, bodyStr, nil
	}

	logger.Info("Network connection is fine")
	return false, "", nil
}
