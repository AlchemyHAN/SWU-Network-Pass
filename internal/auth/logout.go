package auth

import (
	"context"
	"net/http"
)

func Logout(ctx context.Context, client *http.Client) error {
	logoutURL := "http://222.198.127.170/eportal/InterFace.do?method=logoutAndCancelUserMabInfo"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, logoutURL, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
