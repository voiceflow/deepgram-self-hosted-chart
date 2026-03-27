package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Config holds all settings for the validation run.
type Config struct {
	APIURL        string
	FluxURL       string
	APIKey        string
	Model         string
	BatchModel    string
	FluxModel     string
	Language      string
	Audio         []byte
	SmartFormat   bool
	SkipFlux      bool
	SkipBatch     bool
	SkipStreaming bool
}

func main() {
	apiURL := flag.String("url", "localhost:8080", "Deepgram API host:port")
	apiKey := flag.String("key", "", "Deepgram API key (or set DEEPGRAM_API_KEY env)")
	model := flag.String("model", "nova-2-general", "Model name for streaming/flux STT requests")
	batchModel := flag.String("batch-model", "", "Model for batch STT (defaults to -model value)")
	fluxModel := flag.String("flux-model", "", "Model for Flux STT (defaults to -model value)")
	fluxURL := flag.String("flux-url", "", "Flux API host:port if different from -url")
	language := flag.String("language", "en", "Language code (e.g. en, es, fr)")
	audio := flag.String("audio", "", "Path to audio file (default: ../benchmarking/audio.8k.wav)")
	smartFormat := flag.Bool("smart-format", true, "Enable smart formatting (requires NER model for Nova-3)")
	skipFlux := flag.Bool("skip-flux", false, "Skip Flux STT test")
	skipBatch := flag.Bool("skip-batch", false, "Skip Batch STT test")
	skipStreaming := flag.Bool("skip-streaming", false, "Skip Streaming STT test")
	flag.Parse()

	key := *apiKey
	if key == "" {
		key = os.Getenv("DEEPGRAM_API_KEY")
	}
	if key == "" {
		fmt.Println("ERROR: API key required. Use -key flag or set DEEPGRAM_API_KEY env var.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Deepgram Self-Hosted Validation")
	fmt.Println("================================")

	audioData, err := loadAudio(*audio)
	if err != nil {
		fmt.Printf("  ERROR %v\n", err)
		os.Exit(1)
	}

	bm := *batchModel
	if bm == "" {
		bm = *model
	}
	fm := *fluxModel
	if fm == "" {
		fm = *model
	}
	fu := *fluxURL
	if fu == "" {
		fu = *apiURL
	}

	cfg := &Config{
		APIURL:        *apiURL,
		FluxURL:       fu,
		APIKey:        key,
		Model:         *model,
		BatchModel:    bm,
		FluxModel:     fm,
		Language:       *language,
		Audio:         audioData,
		SmartFormat:   *smartFormat,
		SkipFlux:      *skipFlux,
		SkipBatch:     *skipBatch,
		SkipStreaming: *skipStreaming,
	}

	tests := []struct {
		name string
		fn   func(*Config) error
		skip bool
	}{
		{"Health Check", testHealth, false},
		{"Batch STT", testBatch, cfg.SkipBatch},
		{"Streaming STT", testStreaming, cfg.SkipStreaming},
		{"Flux STT", testFlux, cfg.SkipFlux},
	}

	passed, failed, skipped := 0, 0, 0
	fmt.Println()

	for _, t := range tests {
		if t.skip {
			fmt.Printf("  SKIP  %s\n", t.name)
			skipped++
			continue
		}
		fmt.Printf("  RUN   %s\n", t.name)
		start := time.Now()
		if err := t.fn(cfg); err != nil {
			fmt.Printf("  FAIL  %s (%s): %v\n", t.name, time.Since(start).Round(time.Millisecond), err)
			failed++
		} else {
			fmt.Printf("  PASS  %s (%s)\n", t.name, time.Since(start).Round(time.Millisecond))
			passed++
		}
		fmt.Println()
	}

	fmt.Printf("Results: %d passed, %d failed, %d skipped\n\n", passed, failed, skipped)
	if failed > 0 {
		os.Exit(1)
	}
}
