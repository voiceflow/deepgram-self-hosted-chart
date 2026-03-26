package main

// EngineStatus is the response from GET /v1/status/engine.
type EngineStatus struct {
	EngineConnectionStatus string `json:"engine_connection_status"`
}

// BatchResponse is the response from POST /v1/listen.
type BatchResponse struct {
	Results BatchResults `json:"results"`
}

type BatchResults struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string `json:"transcript"`
	Words      []Word `json:"words"`
}

type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// StreamingMessage is a WebSocket message from streaming or Flux STT.
type StreamingMessage struct {
	Type        string  `json:"type"`
	IsFinal     bool    `json:"is_final"`
	SpeechFinal bool    `json:"speech_final"`
	Channel     Channel `json:"channel"`
	Start       float64 `json:"start"`
	Duration    float64 `json:"duration"`
}

// FluxMessage is a WebSocket message from Flux v2 STT.
type FluxMessage struct {
	Type                 string     `json:"type"`
	Event                string     `json:"event"`
	TurnIndex            int        `json:"turn_index"`
	Transcript           string     `json:"transcript"`
	Words                []FluxWord `json:"words"`
	EndOfTurnConfidence  float64    `json:"end_of_turn_confidence"`
	SequenceID           int        `json:"sequence_id"`
}

type FluxWord struct {
	Word       string  `json:"word"`
	Confidence float64 `json:"confidence"`
}
