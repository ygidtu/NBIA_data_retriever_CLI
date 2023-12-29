package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Token is used to handle the NBIA official token request
/*
Official example be like:
curl -X -v -d "username=nbia_guest&password=&client_id=NBIA&grant_type=password" -X POST -k https://services.cancerimagingarchive.net/nbia-api/oauth/token
*/
type Token struct {
	AccessToken      string    `json:"access_token"`
	SessionState     string    `json:"session_state"`
	ExpiresIn        int       `json:"expires_in"`
	NotBeforePolicy  int       `json:"not-before-policy"`
	RefreshExpiresIn int       `json:"refresh_expires_in"`
	Scope            string    `json:"scope"`
	IdToken          string    `json:"id_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	ExpiredTime      time.Time `json:"expires_time"`
}

func makeURL(url_ string, values map[string]interface{}) (string, error) {
	u, err := url.Parse(url_)
	if err != nil {
		return url_, fmt.Errorf("failed to parse url: %v", err)
	}

	queries := u.Query()

	for k, v := range values {
		queries.Set(k, fmt.Sprintf("%v", v))
	}

	u.RawQuery = queries.Encode()
	return u.String(), nil
}

// NewToken create token from official NBIA API
func NewToken(username, passwd, path string) (*Token, error) {
	token := new(Token)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		logger.Infof("restore token from %v", path)
		err = token.Load(path)
		if err != nil {
			logger.Error(err)
			logger.Infof("create new token")
		}
		if token.ExpiredTime.Compare(time.Now()) <= 0 {
			logger.Warn("token expired, create new token")
		} else {
			return token, nil
		}
	}

	// format the token request url
	url_, err := makeURL(
		TokenUrl,
		map[string]interface{}{
			"username":   username,
			"password":   passwd,
			"client_id":  "NBIA",
			"grant_type": "password",
		})
	if err != nil {
		return nil, fmt.Errorf("error creating token: %v", err)
	}

	req, err := http.NewRequest("POST", url_, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("failed to do request: %v", err)
	}

	content, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		logger.Errorf("failed to read response data: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read token from response body: %v", err)
	}

	err = json.Unmarshal(content, token)

	token.ExpiredTime = time.Now().Local().Add(time.Second * time.Duration(token.ExpiresIn))
	return token, token.Dump(path)
}

// Dump is used to save token information
func (token *Token) Dump(path string) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open token json: %v", err)
	}

	content, err := json.MarshalIndent(token, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("failed to dump token: %v", err)
	}

	return f.Close()
}

// Load restore token from json
func (token *Token) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open token json: %v", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read token: %v", err)
	}
	err = json.Unmarshal(content, token)
	if err != nil {
		return fmt.Errorf("failed to unmarshal token: %v", err)
	}

	return f.Close()
}
