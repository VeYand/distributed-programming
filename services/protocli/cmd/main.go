package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const serverURL = "http://protokey:8082"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: protocli [set|get|keys] [arguments]")
		os.Exit(1)
	}

	switch cmd := os.Args[1]; cmd {
	case "set":
		if len(os.Args) != 4 {
			fmt.Println("Usage: protocli set <key> <value>")
			os.Exit(1)
		}
		setKey(os.Args[2], os.Args[3])

	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: protocli get <key>")
			os.Exit(1)
		}
		getKey(os.Args[2])

	case "keys":
		if len(os.Args) != 3 {
			fmt.Println("Usage: protocli keys <prefix>")
			os.Exit(1)
		}
		getKeys(os.Args[2])

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func setKey(key, value string) {
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Value must be a valid 32-bit integer")
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"key":   key,
		"value": v,
	}
	body, status, err := sendRequest("POST", "/set", payload)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	fmt.Println("OK")
}

func getKey(key string) {
	body, status, err := sendRequest("GET", "/get?key="+key, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	var resp struct {
		Value int `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("Invalid response format:", err)
		os.Exit(1)
	}

	fmt.Println(resp.Value)
}

func getKeys(prefix string) {
	body, status, err := sendRequest("GET", "/keys?prefix="+prefix, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	var resp struct {
		Keys []string `json:"keys"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("Invalid response format:", err)
		os.Exit(1)
	}

	for _, k := range resp.Keys {
		fmt.Println(k)
	}
}

func sendRequest(method, endpoint string, payload interface{}) ([]byte, int, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, fmt.Errorf("json.Marshal: %w", err)
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, serverURL+endpoint, body)
	if err != nil {
		return nil, 0, fmt.Errorf("http.NewRequest: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}
