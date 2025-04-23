package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rankcalculator/pkg/app/service"
	"time"
)

func NewClient(apiURL, apiKey string) service.CentrifugoClient {
	return &client{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type client struct {
	apiURL string
	apiKey string
	client *http.Client
}

func (c *client) Publish(channel string, data interface{}) error {
	payload := struct {
		Method  string      `json:"method"`
		Channel string      `json:"channel"`
		Data    interface{} `json:"data"`
	}{
		Method:  "publish",
		Channel: channel,
		Data:    data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey "+c.apiKey)

	log.Printf("Publishing to centrifugo channel %s: %s", channel, string(body))
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Centrifugo response status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
