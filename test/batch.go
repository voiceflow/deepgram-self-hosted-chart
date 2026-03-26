package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func testBatch(cfg *Config) error {
	url := fmt.Sprintf("http://%s/v1/listen?model=%s&language=%s&smart_format=true", cfg.APIURL, cfg.BatchModel, cfg.Language)
	req, err := http.NewRequest("POST", url, bytes.NewReader(cfg.Audio))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+cfg.APIKey)
	req.Header.Set("Content-Type", "audio/wav")

	client := &http.Client{Timeout: 60 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result BatchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if len(result.Results.Channels) == 0 {
		return fmt.Errorf("no channels in response")
	}
	if len(result.Results.Channels[0].Alternatives) == 0 {
		return fmt.Errorf("no alternatives in response")
	}

	transcript := result.Results.Channels[0].Alternatives[0].Transcript
	fmt.Printf("         Transcript: %q\n", transcript)
	fmt.Printf("         Responded in %s\n", elapsed.Round(time.Millisecond))
	return nil
}
