package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"math/big"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger *slog.Logger
var session *http.Client

type Account struct {
	Username string
	Password string
}

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	session = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "tcp4", addr) // 强制使用 IPV4, 避免 CERNET 网络环境下 IPV6 可以免认证访问互联网的问题
			},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 6 * time.Second,
		},
	}
}

func getAccount() Account {
	// 判断accounts.txt是否存在
	if _, err := os.Stat("accounts.txt"); os.IsNotExist(err) {
		logger.Error("File accounts.txt not found, please create it and put your accounts in it")
		os.Exit(1)
	}

	file, err := os.Open("accounts.txt")
	if err != nil {
		logger.Error("Error in opening file accounts.txt", "error", err)
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var accounts []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			accounts = append(accounts, line)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error in reading file accounts.txt", "error", err)
		panic(err)
	}

	randomIndex := int(math.Round(rand.Float64() * float64(len(accounts)-1)))
	account := strings.Fields(accounts[randomIndex])
	if len(account) != 2 {
		logger.Error("Invalid account format", "account", account)
		os.Exit(1)
	}
	return Account{Username: account[0], Password: account[1]}
}

func main() {
	logger.Info("Starting the application...")

	// 判断格式是否正确
	account := getAccount()
	if account.Username == "" || account.Password == "" {
		logger.Error("Invalid account format")
		os.Exit(1)
	}

	// heartbeat loop
	for {
		verificationUrl := url.URL{Scheme: "http", Host: "captive.apple.com"}
		success, redirectHtmlString := verifyNetworkStatus(verificationUrl)
		if !success && redirectHtmlString != "" {
			redirectUrl := getRedirectUrl(redirectHtmlString)
			if redirectUrl == (url.URL{}) {
				logger.Error("Failed to get redirect URL")
				time.Sleep(20 * time.Second)
				continue
			}
			account := getAccount()
			username := account.Username
			password := account.Password
			logger.Info("Trying to login...", "username", username)
			logger.Info("Trying to login...", "password", password)
			envNeedEncryption := os.Getenv("SWU_NEED_ENCRYPTION")
			needEncryption := envNeedEncryption == "" || envNeedEncryption == "true"
			if needEncryption {
				password = encryptPassword(redirectUrl, password)
			}
			login(username, password, redirectUrl, needEncryption)
		}
		time.Sleep(5 * time.Second)
	}
}

func login(username string, password string, redirectUrl url.URL, needEncryption bool) bool {
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

	req, err := http.NewRequest("POST", authenticationUrl, strings.NewReader(postData.Encode()))
	if err != nil {
		logger.Error("Failed to create request", "error", err)
		return false
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

	resp, err := session.Do(req)
	if err != nil {
		logger.Error("Failed to login", "error", err)
		return false
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", "error", err)
		return false
	}
	if strings.Contains(string(bodyBytes), "success") {
		logger.Info("Login success")
		return true
	} else {
		var resultJson map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &resultJson); err != nil {
			logger.Error("Failed to unmarshal JSON", "error", err)
		}
		message, ok := resultJson["message"].(string)
		if !ok {
			logger.Error("Fail reason message is missing")
		}
		logger.Error("Login failed", "reason", message)
		return false
	}
}

func verifyNetworkStatus(verificationUrl url.URL) (bool, string) {
	resp, err := session.Get(verificationUrl.String())

	if err != nil {
		logger.Error("Failed to verify network status", "error", err)
		return false, ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to verify network status", "error", err)
		return false, ""
	}
	// 如果返回的 HTML 中包含 <script> 标签, 则说明需要进行网络认证
	if strings.Contains(string(body), "<script>") {
		logger.Warn("Oops, Network connection is down!!!")
		return false, string(body)
	} else {
		logger.Info("Network connection is fine")
		return true, ""
	}
}

func getRedirectUrl(redirectHtmlString string) url.URL {
	start := strings.Index(redirectHtmlString, "href='")
	if start == -1 {
		logger.Error("href not found in the HTML string")
		return url.URL{}
	}
	start += len("href='")
	end := strings.Index(redirectHtmlString[start:], "'")
	if end == -1 {
		logger.Error("Referer end delimiter not found in the HTML string")
		return url.URL{}
	}
	redirectUrlString := redirectHtmlString[start : start+end]
	redirectUrl, err := url.Parse(redirectUrlString)
	if err != nil {
		logger.Error("Failed to parse redirect URL", "error", err)
		return url.URL{}
	}
	return *redirectUrl
}

func encryptPassword(redirectUrl url.URL, password string) string {
	pageInfoUrl := "http://" + redirectUrl.Host + "/eportal/InterFace.do?method=pageInfo&queryString=undefined"
	postData := url.Values{}
	queryString := redirectUrl.Query()
	postData.Add("queryString", queryString.Encode())
	req, err := http.NewRequest("GET", pageInfoUrl, nil)
	if err != nil {
		logger.Error("Failed to create request", "error", err)
		return ""
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

	resp, err := session.Do(req)
	if err != nil {
		logger.Error("Failed to get page info", "error", err)
		return ""
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Error("Failed to read response body", "error", err)
		return ""
	}
	var resultJson map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &resultJson); err != nil {
		logger.Error("Failed to unmarshal JSON", "error", err)
		return ""
	}

	publicKeyExponent, ok := resultJson["publicKeyExponent"].(string)
	if !ok {
		logger.Error("publicKeyExponent is not a string or is missing")
		return ""
	}

	publicKeyModulus, ok := resultJson["publicKeyModulus"].(string)
	if !ok {
		logger.Error("publicKeyModulus is not a string or is missing")
		return ""
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
	return encryptedHex
}
