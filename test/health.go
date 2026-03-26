package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func testHealth(cfg *Config) error {
	url := fmt.Sprintf("http://%s/v1/status/engine", cfg.APIURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+cfg.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var status EngineStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if status.EngineConnectionStatus != "Connected" {
		return fmt.Errorf("engine status: %q (expected \"Connected\")", status.EngineConnectionStatus)
	}

	fmt.Printf("         engine_connection_status: %s\n", status.EngineConnectionStatus)
	return nil
}
