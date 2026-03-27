package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func testFlux(cfg *Config) error {
	url := fmt.Sprintf("ws://%s/v2/listen?model=%s",
		cfg.FluxURL, cfg.FluxModel)

	fmt.Printf("         URL: %s\n", url)

	header := http.Header{}
	header.Set("Authorization", "Token "+cfg.APIKey)

	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("websocket dial failed (HTTP %d): %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("websocket dial failed: %w", err)
	}
	defer conn.Close()

	// Send audio chunks in background
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		sendAudioChunks(conn, cfg.Audio)
	}()

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	messages := 0
	lastPrinted := 0
	gotTurnInfo := false
	lastTranscript := ""

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				break
			}
			if gotTurnInfo {
				break
			}
			return fmt.Errorf("reading message: %w", err)
		}

		var m FluxMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		if m.Type == "TurnInfo" {
			gotTurnInfo = true
			messages++
			lastTranscript = m.Transcript

			// Only print new characters when transcript grows
			if len(m.Transcript) > lastPrinted {
				fmt.Printf("%s", m.Transcript[lastPrinted:])
				lastPrinted = len(m.Transcript)
			}

			if m.Event == "EndOfTurn" {
				// Reset for next turn
				lastPrinted = 0
				fmt.Println()
			}
		}
	}

	<-sendDone

	fmt.Println()
	if lastTranscript != "" {
		fmt.Printf("         Final: %s\n", lastTranscript)
	}

	if !gotTurnInfo {
		return fmt.Errorf("no TurnInfo messages received")
	}

	fmt.Printf("         %d TurnInfo messages\n", messages)
	return nil
}
