package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type wsResult struct {
	Messages       int
	Finals         int
	Transcript     string
	AllTranscripts []string
}

func runWebSocketTest(cfg *Config, wsURL string, opts ...bool) (*wsResult, error) {
	verbose := len(opts) > 0 && opts[0]
	header := http.Header{}
	header.Set("Authorization", "Token "+cfg.APIKey)

	fmt.Printf("         URL: %s\n", wsURL)
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("websocket dial failed (HTTP %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}
	defer conn.Close()

	// Send audio chunks in background
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		sendAudioChunks(conn, cfg.Audio)
	}()

	// Read messages
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	res := &wsResult{}
	gotMetadata := false
	gotResults := false

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			// If we already got results, a read error after close is expected
			if gotResults {
				break
			}
			return nil, fmt.Errorf("reading message: %w", err)
		}

		var m StreamingMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			if verbose {
				fmt.Printf("         [RAW] %s\n", string(msg))
			}
			continue
		}

		if verbose && m.Type != "Results" {
			fmt.Printf("         [%s] %s\n", m.Type, string(msg))
		}

		switch m.Type {
		case "Metadata":
			gotMetadata = true
		case "Results":
			gotResults = true
			res.Messages++
			if len(m.Channel.Alternatives) > 0 {
				t := m.Channel.Alternatives[0].Transcript
				if t != "" {
					if m.IsFinal {
						res.Finals++
						res.AllTranscripts = append(res.AllTranscripts, t)
						if res.Transcript == "" {
							res.Transcript = t
						}
						// Print final result in green
						fmt.Printf("\033[32m%s \033[0m", t)
					} else {
						// Print interim in gray
						fmt.Printf("\033[90m%s \033[0m", t)
					}
				}
			}
		default:
			// Try to extract transcript from unknown message types (e.g. Flux v2)
			var raw map[string]interface{}
			if err := json.Unmarshal(msg, &raw); err == nil {
				if ch, ok := raw["channel"].(map[string]interface{}); ok {
					if alts, ok := ch["alternatives"].([]interface{}); ok && len(alts) > 0 {
						if alt, ok := alts[0].(map[string]interface{}); ok {
							if t, ok := alt["transcript"].(string); ok && t != "" {
								gotResults = true
								res.Messages++
								res.Finals++
								res.AllTranscripts = append(res.AllTranscripts, t)
								if res.Transcript == "" {
									res.Transcript = t
								}
								fmt.Printf("\r         \033[32m%s\033[0m\n", t)
							}
						}
					}
				}
			}
		}
	}

	<-sendDone

	if !gotMetadata {
		return nil, fmt.Errorf("no Metadata message received")
	}
	if !gotResults {
		return nil, fmt.Errorf("no Results message received")
	}
	if res.Finals == 0 {
		return nil, fmt.Errorf("no final results received")
	}

	return res, nil
}

func sendAudioChunks(conn *websocket.Conn, audio []byte) {
	const chunkSize = 800
	const interval = 50 * time.Millisecond

	offset := 0
	for offset < len(audio) {
		end := offset + chunkSize
		if end > len(audio) {
			end = len(audio)
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, audio[offset:end]); err != nil {
			return
		}
		offset = end
		time.Sleep(interval)
	}

	// Signal end of audio
	conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"CloseStream"}`))
}
