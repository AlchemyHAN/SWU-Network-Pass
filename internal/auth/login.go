package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Login(ctx context.Context, logger *slog.Logger, client *http.Client, username string, password string, redirectUrl url.URL, needEncryption bool) (bool, error) {
	authenticationUrl := "http://" + redirectUrl.Host + "/eportal/InterFace.do?method=login"

	queryString := redirectUrl.Query()

	postData := url.Values{}
	postData.Add("userId", username)
	postData.Add("password", password)
	postData.Add("service", "%E9%BB%98%E8%AE%A4")
	postData.Add("queryString", queryString.Encode())
	postData.Add("operatorPwd", "")
	postData.Add("operatorUserId", "")
	postData.Add("validcode", "")
	postData.Add("passwordEncrypt", strconv.FormatBool(needEncryption))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authenticationUrl, strings.NewReader(postData.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Host", redirectUrl.Host)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Referer", "http://123.123.123.123/")
	req.Header.Add("Origin", "http://"+redirectUrl.Host)

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}
	if strings.Contains(string(bodyBytes), "success") {
		logger.Info("Login successful", "username", username)
		return true, nil
	} else {
		var resultJson map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &resultJson); err != nil {
			return false, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		message, ok := resultJson["message"].(string)
		if !ok {
			return false, fmt.Errorf("login failed: fail reason message is missing")
		}
		return false, fmt.Errorf("login failed: %s", message)
	}
}

func GetRedirectUrl(redirectHtmlString string) (url.URL, error) {
	start := strings.Index(redirectHtmlString, "href='")
	if start == -1 {
		return url.URL{}, fmt.Errorf("href not found in the HTML string")
	}
	start += len("href='")
	end := strings.Index(redirectHtmlString[start:], "'")
	if end == -1 {
		return url.URL{}, fmt.Errorf("referer end delimiter not found in the HTML string")
	}
	redirectUrlString := redirectHtmlString[start : start+end]
	redirectUrl, err := url.Parse(redirectUrlString)
	if err != nil {
		return url.URL{}, fmt.Errorf("failed to parse redirect URL: %w", err)
	}
	return *redirectUrl, nil
}
