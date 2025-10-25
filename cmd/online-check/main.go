package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"swu.social/swu-network-pass/internal/account"
	"swu.social/swu-network-pass/internal/auth"
	"swu.social/swu-network-pass/internal/client"
	"swu.social/swu-network-pass/internal/crypto"
	"swu.social/swu-network-pass/internal/status"
)

const (
	heartbeatURL        = "http://captive.apple.com/"
	ipCheckURL          = "http://4.ipw.cn"
	iterationTimeout    = 15 * time.Second
	statusCheckTimeout  = 6 * time.Second
	loginTimeout        = 10 * time.Second
	statusCheckInterval = 5 * time.Second
)

var (
	accountFilePath = "accounts.txt"
	needEncryption  = false
	selectedISP     = ""
)

func main() {
	flag.StringVar(&accountFilePath, "accounts", "accounts.txt", "Path to the accounts file")
	flag.StringVar(&selectedISP, "isp", "", "Select ISP: CMCC, CT, CU, or leave empty for random")
	flag.BoolVar(&needEncryption, "encryption", false, "Whether password encryption is needed")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	httpClient := client.New()
	verificationURL, err := url.Parse(heartbeatURL)
	if err != nil {
		logger.Error("Failed to parse verification URL", "error", err)
		return
	}

	logger.Info("Starting SWU network heartbeat", "need_encryption", needEncryption)

	for {
		if err := runIteration(logger, httpClient, verificationURL, needEncryption); err != nil {
			logger.Error("Iteration failed", "error", err)
		}
		time.Sleep(statusCheckInterval)
	}
}

func runIteration(logger *slog.Logger, httpClient *http.Client, verificationURL *url.URL, needEncryption bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), iterationTimeout)
	defer cancel()

	verifyCtx, cancelVerify := context.WithTimeout(ctx, statusCheckTimeout)
	defer cancelVerify()

	requiresLogin, redirectHTML, err := status.VerifyOnline(verifyCtx, logger, httpClient, verificationURL)
	if err != nil {
		logger.Error("Failed to verify online status", "error", err)
	}

	if !requiresLogin {
		return nil
	}

	if err := performLogin(ctx, logger, httpClient, redirectHTML, needEncryption); err != nil {
		return fmt.Errorf("perform login: %w", err)
	}

	if selectedISP != "" {
		checkCtx, cancelCheck := context.WithTimeout(ctx, statusCheckTimeout)
		defer cancelCheck()
		isRightISP, err := status.CheckISP(checkCtx, httpClient, ipCheckURL, selectedISP)
		if err != nil {
			return fmt.Errorf("check ISP: %w", err)
		}
		if !isRightISP {
			logger.Warn("Connected to wrong ISP network", "expected_isp", selectedISP)
			auth.Logout(checkCtx, httpClient)
			logger.Info("Logged out due to ISP mismatch; will attempt re-login in next iteration")
		}
	}

	return nil
}

func performLogin(ctx context.Context, logger *slog.Logger, httpClient *http.Client, redirectHTML string, needEncryption bool) error {

	redirectURL, err := auth.GetRedirectUrl(redirectHTML)
	if err != nil {
		return fmt.Errorf("get redirect url: %w", err)
	}

	acct, err := account.GetRandomAccount(logger, accountFilePath)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	logger.Info("Attempting login", "username", acct.Username, "need_encryption", needEncryption)

	loginCtx, cancelLogin := context.WithTimeout(ctx, loginTimeout)
	defer cancelLogin()

	password := acct.Password
	if needEncryption {
		encrypted, err := crypto.EncryptPassword(loginCtx, httpClient, redirectURL, password)
		if err != nil {
			return fmt.Errorf("encrypt password: %w", err)
		}
		password = encrypted
	}

	if _, err := auth.Login(loginCtx, logger, httpClient, acct.Username, password, redirectURL, needEncryption); err != nil {
		return fmt.Errorf("login: %w", err)
	}
	return nil
}
