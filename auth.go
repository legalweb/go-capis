package capis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrAuthorizationFailed = errors.New("authorization failed")

type StaticToken string

type PasswordAuthentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ttl      time.Time
	token    string
}

func (a *PasswordAuthentication) Token(ctx context.Context) (string, error) {
	if a.ttl.After(time.Now()) {
		return a.token, nil
	}

	rb, _ := json.Marshal(a)
	req, _ := http.NewRequest("POST", DefaultBaseURL+"/auth", bytes.NewReader(rb))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to get response from auth %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrAuthorizationFailed
	}

	rb, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response from auth %w", err)
	}

	var data struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rb, &data); err != nil {
		return "", fmt.Errorf("malformed response from auth %w", err)
	}

	a.token = data.Token
	a.ttl = time.Now().Add(time.Hour)

	return a.token, nil
}

func (a *PasswordAuthentication) AuthorizeRequest(ctx context.Context, req *http.Request) error {
	v, err := a.Token(ctx)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+string(v))
	return nil
}

func (v StaticToken) AuthorizeRequest(ctx context.Context, req *http.Request) error {
	req.Header.Add("Authorization", "Bearer "+string(v))
	return nil
}
