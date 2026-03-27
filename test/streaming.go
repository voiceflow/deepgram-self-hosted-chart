package main

import "fmt"

func testStreaming(cfg *Config) error {
	url := fmt.Sprintf("ws://%s/v1/listen?model=%s&language=%s&smart_format=%t&punctuate=true&interim_results=true",
		cfg.APIURL, cfg.Model, cfg.Language, cfg.SmartFormat)

	res, err := runWebSocketTest(cfg, url)
	if err != nil {
		return err
	}

	if len(res.AllTranscripts) > 0 {
		fmt.Println("         Full transcript:")
		for i, t := range res.AllTranscripts {
			fmt.Printf("           [%d] %s\n", i+1, t)
		}
	}
	fmt.Printf("         %d messages, %d final\n", res.Messages, res.Finals)
	return nil
}
