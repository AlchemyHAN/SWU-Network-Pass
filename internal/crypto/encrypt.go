package crypto

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
)

func EncryptPassword(ctx context.Context, client *http.Client, redirectUrl url.URL, password string) (string, error) {
	pageInfoUrl := "http://" + redirectUrl.Host + "/eportal/InterFace.do?method=pageInfo"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageInfoUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed in creating request: %w", err)
	}
	req.Header.Add("Host", redirectUrl.Host)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Referer", "http://123.123.123.123/")
	req.Header.Add("Origin", "http://"+redirectUrl.Host)

	queryString := redirectUrl.Query()
	values := url.Values{}
	values.Add("queryString", queryString.Encode())
	req.URL.RawQuery = values.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get page info: %w", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	var resultJson map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &resultJson); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	publicKeyExponent, ok := resultJson["publicKeyExponent"].(string)
	if !ok {
		return "", fmt.Errorf("publicKeyExponent is not a string or is missing")
	}

	publicKeyModulus, ok := resultJson["publicKeyModulus"].(string)
	if !ok {
		return "", fmt.Errorf("publicKeyModulus is not a string or is missing")
	}

	rsaE, _ := new(big.Int).SetString(publicKeyExponent, 16)
	rsaN, _ := new(big.Int).SetString(publicKeyModulus, 16)

	macAddr := redirectUrl.Query().Get("mac")
	secret := password + ">" + macAddr
	secretInt := new(big.Int)
	secretInt.SetBytes([]byte(secret))

	// 执行幂模运算
	encryptedBigInt := new(big.Int).Exp(secretInt, rsaE, rsaN)

	// 转换结果为十六进制字符串（去除前缀 "0x"，如果有）
	encryptedHex := hex.EncodeToString(encryptedBigInt.Bytes())
	return encryptedHex, nil
}
