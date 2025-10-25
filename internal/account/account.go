package account

import (
	"bufio"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strings"
)

type Account struct {
	Username string
	Password string
}

func GetRandomAccount(logger *slog.Logger, filepath string) (Account, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if logger != nil {
			logger.Error("Accounts file not found", "path", filepath)
		}
		return Account{}, fmt.Errorf("file %s not found, please create it and put your accounts in it", filepath)
	}

	file, err := os.Open(filepath)

	if err != nil {
		if logger != nil {
			logger.Error("Failed to open accounts file", "path", filepath, "error", err)
		}
		return Account{}, fmt.Errorf("error in opening file %s: %v", filepath, err)
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
		if logger != nil {
			logger.Error("Failed to read accounts file", "path", filepath, "error", err)
		}
		return Account{}, fmt.Errorf("error in reading file %s: %v", filepath, err)
	}

	if len(accounts) == 0 {
		if logger != nil {
			logger.Error("No accounts found in file", "path", filepath)
		}
		return Account{}, fmt.Errorf("no accounts found in file %s", filepath)
	}

	randomIndex := rand.IntN(len(accounts))
	account := strings.Fields(accounts[randomIndex])
	if len(account) != 2 {
		if logger != nil {
			logger.Error("Invalid account format", "path", filepath, "raw", accounts[randomIndex])
		}
		return Account{}, fmt.Errorf("invalid account format in file %s: %v", filepath, account)
	}
	return Account{Username: account[0], Password: account[1]}, nil
}
