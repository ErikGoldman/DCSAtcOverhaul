package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/dharmab/skyeye/pkg/coalitions"
	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	atcclient "github.com/ErikGoldman/DCSAtcOverhaul/pkg/client"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/deepgramRecognizer"
	deepgramspeaker "github.com/ErikGoldman/DCSAtcOverhaul/pkg/deepgramSpeaker"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/dharmab/skyeye/pkg/telemetry"
)

func main() {
	telemetryAddress := flag.String("telemetryAddress", "", "The address of the Tacview server")
	if telemetryAddress == nil {
		panic("No telemetry address defined")
	}

	config := types.ClientConfiguration{
		Address:                   "192.168.86.40:5002",
		ClientName:                "test",
		ExternalAWACSModePassword: "test",
		GUID:                      uuid.New().String(),
		Coalition:                 coalitions.Blue,
		ConnectionTimeout:         10 * time.Second,
		AllowRecording:            true,
		Mute:                      false,
		Radios: []types.Radio{
			{
				Frequency:        305000000.0,
				IsEncrypted:      false,
				ShouldRetransmit: true,
				Modulation:       types.ModulationAM,
			},
		},
	}

	// Read and parse the config.json file
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config.json")
	}

	var configData struct {
		Deepgram struct {
			APIKey string `json:"api_key"`
		} `json:"deepgram"`
	}

	err = json.Unmarshal(configFile, &configData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config.json")
	}

	recognizer := deepgramRecognizer.NewAtcDeepgramRecognizer(configData.Deepgram.APIKey)

	var telemetryClient telemetry.Client
	log.Info().Str("address", *telemetryAddress).Msg("constructing telemetry client")
	telemetryClient = telemetry.NewTelemetryClient(
		*telemetryAddress,
		"hostname",
		"",            // password
		500,           //timeout
		2*time.Second, // refresh rate in seconds
	)

	a := &atcclient.AtcApplication{
		Recognizer:                 recognizer,
		EnableTranscriptionLogging: true,
		TranscribedMessages:        make(chan message.Message[string]),
		CommandProcessor:           atcclient.LoadCommandProcessor(),
		SpeechSynthesizer:          deepgramspeaker.NewSpeechSynthesizer(configData.Deepgram.APIKey),
		TelemetryClient:            telemetryClient,
	}

	log.Info().Msgf("config: %v", config)

	srsClient, err := simpleradio.NewClient(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create SRS client")
		return
	}

	log.Info().Msgf("running")

	go a.Start(srsClient)
	var wg sync.WaitGroup

	srsClient.Run(
		context.Background(),
		&wg,
	)

	log.Info().Msgf("done")
}
