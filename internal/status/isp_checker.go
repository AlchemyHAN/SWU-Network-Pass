package status

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	ISP_CMCC = "CMCC"
	ISP_CT   = "CT"
	ISP_CU   = "CU"

	CIDR_CMCC = "183.230.225.0/24"
	CIDR_CT   = "222.178.202.0/24"
	CIDR_CU   = "221.7.112.0/24"
)

func GetISPCIDR(isp string) string {
	switch isp {
	case ISP_CMCC:
		return CIDR_CMCC
	case ISP_CT:
		return CIDR_CT
	case ISP_CU:
		return CIDR_CU
	default:
		return ""
	}
}

func GetIP(ctx context.Context, client *http.Client, ipServiceURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipServiceURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func IsIPInCIDR(ip string, cidr string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, fmt.Errorf("invalid IP address: %s", ip)
	}

	return ipNet.Contains(parsedIP), nil
}

func CheckISP(ctx context.Context, client *http.Client, ipServiceURL string, expectedISP string) (bool, error) {
	ip, err := GetIP(ctx, client, ipServiceURL)
	if err != nil {
		return false, err
	}
	cidr := GetISPCIDR(expectedISP)
	if cidr == "" {
		return false, fmt.Errorf("unknown ISP: %s", expectedISP)
	}

	inCIDR, err := IsIPInCIDR(ip, cidr)
	if err != nil {
		return false, err
	}

	return inCIDR, nil
}
