package authentication

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AuthCheckResponse struct {
	Authenticated bool   `json:"authenticated"`
	UserID        string `json:"user_id,omitempty"`
}

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *Client) IsAuthenticatedFromRequest(r *http.Request) (string, bool, error) {
	url := c.BaseURL + "/internal/auth/check"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", false, err
	}

	for _, ck := range r.Cookies() {
		req.AddCookie(ck)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", false, errors.New("access forbidden")
	}
	if resp.StatusCode >= 500 {
		return "", false, errors.New("internal auth service error: " + resp.Status)
	}

	var out AuthCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", false, err
	}
	return out.UserID, out.Authenticated, nil
}
