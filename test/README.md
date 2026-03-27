# Deepgram Self-Hosted Validation Tool

A Go CLI tool to smoke-test a Deepgram self-hosted deployment. Validates health, streaming STT (v1), and Flux STT (v2) endpoints with real-time transcript streaming.

## Prerequisites

- Go 1.22+
- Network access to the Deepgram API (via port-forward or direct)
- A Deepgram API key

## Setup

Port-forward the Deepgram API service:

```bash
# Standard API (Nova streaming)
kubectl port-forward svc/deepgram-api-external 8080:8080 -n deepgram

# If Flux runs on a separate service/port
kubectl port-forward svc/deepgram-flux-api-external 8081:8080 -n deepgram
```

## Build & Run

```bash
cd test
go build -o dg-validate .

# Using env var
export DEEPGRAM_API_KEY=<your-api-key>

# Test streaming STT only
./dg-validate -model nova-2-general -skip-batch -skip-flux

# Test Flux STT only
./dg-validate -skip-batch -skip-streaming -flux-model flux-general-en

# Test Flux on a separate port
./dg-validate -skip-batch -skip-streaming -flux-url localhost:8081 -flux-model flux-general-en

# Test both streaming and Flux
./dg-validate -model nova-2-general -skip-batch -flux-model flux-general-en

# Test Flux on a different API than streaming
./dg-validate -model nova-2-general -skip-batch \
  -flux-url localhost:8081 -flux-model flux-general-en

# Test Nova-3 without smart formatting (no NER model needed)
./dg-validate -model nova-3 -smart-format=false -skip-batch -skip-flux

# Test Nova-3 with smart formatting (requires entity-detector NER model)
./dg-validate -model nova-3 -skip-batch -skip-flux
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-url` | `localhost:8080` | Deepgram API host:port |
| `-flux-url` | same as `-url` | Flux API host:port (if different from main API) |
| `-key` | `DEEPGRAM_API_KEY` env | Deepgram API key |
| `-model` | `nova-2-general` | Model name for streaming STT |
| `-batch-model` | same as `-model` | Model for batch STT |
| `-flux-model` | same as `-model` | Model for Flux STT (e.g. `flux-general-en`) |
| `-language` | `en` | Language code for streaming/batch (e.g. en, es, fr) |
| `-audio` | `../benchmarking/audio.8k.wav` | Path to test audio file (WAV) |
| `-smart-format` | `true` | Enable smart formatting (requires NER model for Nova-3) |
| `-skip-batch` | `false` | Skip the Batch STT test |
| `-skip-streaming` | `false` | Skip the Streaming STT test |
| `-skip-flux` | `false` | Skip the Flux STT test |

## Tests

| Test | Endpoint | Protocol | Description |
|------|----------|----------|-------------|
| Health Check | `GET /v1/status/engine` | HTTP | Verifies engine is connected |
| Batch STT | `POST /v1/listen` | HTTP | Sends audio file, verifies transcription (requires batch-capable model) |
| Streaming STT | `ws:///v1/listen` | WebSocket | Streams audio in real-time, shows interim and final results |
| Flux STT | `ws:///v2/listen` | WebSocket | Flux turn-based STT, streams transcript as words arrive |

## Notes

- **Streaming models only**: Most Deepgram models (Nova-2, Nova-3) are streaming-only. Use `-skip-batch` unless you have a batch-capable model loaded.
- **Flux uses v2 endpoint**: Flux models require `/v2/listen` and the model name includes the language (e.g. `flux-general-en`), so the `-language` flag is not used for Flux.
- **Port-forward stability**: A failed request (wrong model name, etc.) can crash `kubectl port-forward`. Restart it before retrying.
- **Transcripts stream live**: Streaming STT shows interim (gray) and final (green) results. Flux shows words appearing as they're transcribed.
- **Smart formatting and Nova-3**: Nova-3 on `release-260319`+ requires the `entity-detector` NER model files for `smart_format=true`. If you don't have the NER models loaded, use `-smart-format=false` to test Nova-3 without formatting. Nova-2 works with smart formatting without the NER model.

## Audio File

By default, the tool uses `../benchmarking/audio.8k.wav` (8kHz, 16-bit PCM mono, ~20 seconds). Use the `-audio` flag to specify a different file.
