package main

import (
	"fmt"
	"os"
)

func loadAudio(path string) ([]byte, error) {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		fmt.Printf("  INFO  Using audio file: %s\n", path)
		return data, nil
	}

	defaults := []string{
		"../benchmarking/audio.8k.wav",
		"benchmarking/audio.8k.wav",
	}
	for _, p := range defaults {
		if data, err := os.ReadFile(p); err == nil {
			fmt.Printf("  INFO  Using audio file: %s\n", p)
			return data, nil
		}
	}

	return nil, fmt.Errorf("no audio file found; use -audio flag to specify path (e.g. -audio ../benchmarking/audio.8k.wav)")
}
