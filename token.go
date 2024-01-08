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

curl -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJkZGFhMGY3YS1kZTBmLTRkYWQtYjM1ZS05MjljYjBiMTY3YjgifQ.eyJleHAiOjE3MDE4NDA4NTAsImlhdCI6MTcwMTgzMzY1MCwianRpIjoiZjBiNjY2YTctMDdhYS00NTExLThhOTgtZmU5MTVlMDE5OWY3IiwiaXNzIjoiaHR0cHM6Ly9rZXljbG9hay5kYm1pLmNsb3VkL2F1dGgvcmVhbG1zL1RDSUEiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiZjowMTliNTYzNC1kYWJkLTQyMTEtYTQxZC03MjNjNDRhZmNmZmQ6bmJpYV9ndWVzdCIsInR5cCI6IkJlYXJlciIsImF6cCI6Im5iaWEiLCJzZXNzaW9uX3N0YXRlIjoiZTRlNWFkMWItZmM2ZS00NjhlLWJkMDgtZWU1ZWY3YzFlYjZmIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL3NlcnZpY2VzLmNhbmNlcmltYWdpbmdlYXJjaGl2ZS5uZXQiLCJodHRwczovL25iaWEuY2FuY2VyaW1hZ2luZ2VhcmNoaXZlLm5ldCIsImh0dHBzOi8vd3d3LmNhbmNlcmltYWdpbmdlYXJjaGl2ZS5uZXQiLCIqIiwiaHR0cDovL3RjaWEtbmJpYS0yLmFkLnVhbXMuZWR1OjQ1MjEwIiwiaHR0cHM6Ly9jYW5jZXJpbWFnaW5nZWFyY2hpdmUubmV0IiwiaHR0cDovL3RjaWEtbmJpYS0xLmFkLnVhbXMuZWR1OjQ1MjEwIiwiaHR0cHM6Ly9wdWJsaWMuY2FuY2VyaW1hZ2luZ2VhcmNoaXZlLm5ldCJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJkZWZhdWx0LXJvbGVzLXRjaWEiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJlNGU1YWQxYi1mYzZlLTQ2OGUtYmQwOC1lZTVlZjdjMWViNmYiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6Ik5CSUEgR3Vlc3QiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJuYmlhX2d1ZXN0IiwiZ2l2ZW5fbmFtZSI6Ik5CSUEiLCJmYW1pbHlfbmFtZSI6Ikd1ZXN0IiwiZW1haWwiOiJuYmlhX2d1ZXN0QGNhbmNlcmltYWdpbmdhcmNoaXZlLm5ldCJ9.YwkblULAl_9gw_Wv3IkDjMQXUYEK39AAuR7RVl_X9W4" -k "https://services.cancerimagingarchive.net/nbia-api/services/getDicomTags?SeriesUID=1.3.6.1.4.1.14519.5.2.1.6834.5010.100089621274100103247029607723"

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

	content, err := RespContent(req)
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
