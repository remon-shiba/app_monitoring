package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func AbsDiff(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TimestampToUnix(timestamp string) int64 {
	t, _ := time.Parse(time.RFC3339, timestamp)
	return t.Unix()
}

func SendRequest(baseURL string, method string, body []byte, headers map[string]string, timeout int) (interface{}, error) {
	reqBody := bytes.NewBuffer(body)

	// Create the request
	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return nil, err
	}

	// Set default content-type header if not provided
	if _, exists := headers["Content-Type"]; !exists {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		fmt.Printf("HEADER: %s: %s\n", key, value)
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle empty response
	if len(body) == 0 {
		return nil, nil
	}

	// Try to parse response as JSON object
	var jsonRespObject map[string]interface{}
	if err := json.Unmarshal(body, &jsonRespObject); err == nil {
		return jsonRespObject, nil
	}

	// If parsing as JSON object fails, try as JSON array
	var jsonRespArray []interface{}
	if err := json.Unmarshal(body, &jsonRespArray); err == nil {
		return jsonRespArray, nil
	}

	// If neither parsing works, return an error
	return nil, fmt.Errorf("response is neither a JSON object nor a JSON array: %s", string(body))
}

func SpeakNotif(message string) {
	cmd := exec.Command("espeak", message)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
