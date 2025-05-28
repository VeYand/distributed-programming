package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rankcalculator/pkg/app/auth"
	errors2 "rankcalculator/pkg/app/errors"
	"time"
)

type AuthCheckResponse struct {
	CanReadText bool `json:"can_read_text"`
}

func NewPermissionChecker(baseURL string) auth.PermissionChecker {
	return &client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

type client struct {
	BaseURL string
	HTTP    *http.Client
}

func (c *client) CanReadText(userID, textID string) (bool, error) {
	url := fmt.Sprintf(c.BaseURL+"/internal/permission/text?user_id=%s&text_id=%s", userID, textID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, errors2.ErrStatisticsNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("transport error: " + resp.Status)
	}

	var out AuthCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, err
	}
	return out.CanReadText, nil
}
