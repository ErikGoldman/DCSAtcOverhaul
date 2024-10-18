package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dharmab/skyeye/pkg/coalitions"
	"github.com/dharmab/skyeye/pkg/recognizer"
	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/deepgramRecognizer"
)

type atcApplication struct {
	recognizer                 recognizer.Recognizer
	enableTranscriptionLogging bool
	transcribedMessages        chan Message[string]
}

func (a *atcApplication) srsLoop(radioClient simpleradio.Client) {
	for {
		select {
		case transmission := <-radioClient.Receive():
			log.Info().Msg("received transmission")
			a.recognizeSample(context.Background(), nil, transmission.Audio, a.transcribedMessages)
		}
	}
}

type Message[T any] struct {
	Context context.Context
	Data    T
}

func AsMessage[T any](ctx context.Context, data T) Message[T] {
	return Message[T]{Context: ctx, Data: data}
}

func (a *atcApplication) recognizeSample(processCtx context.Context, requestCtx context.Context, audio simpleradio.Audio, out chan<- Message[string]) {
	recogizerCtx, cancel := context.WithTimeout(processCtx, 30*time.Second)
	defer func() {
		if recogizerCtx.Err() != nil && errors.Is(recogizerCtx.Err(), context.DeadlineExceeded) {
			log.Error().Msg("timeout processing speech")
		}
	}()
	defer cancel()

	log.Info().Msg("recognizing audio sample")
	//start := time.Now()
	text, err := a.recognizer.Recognize(recogizerCtx, audio, a.enableTranscriptionLogging)
	if err != nil {
		log.Error().Err(err).Msg("error recognizing audio sample")
		return
	}

	log.Info().Msgf("recognized text: %s", text)

	/*
		logger := log.With().Stringer("clockTime", time.Since(start)).Logger()

		requestCtx = traces.WithRecognizedAt(requestCtx, time.Now())
		requestCtx = traces.WithRequestText(requestCtx, text)
		if a.enableTranscriptionLogging {
			logger = logger.With().Str("text", text).Logger()
		}
		logger.Info().Msg("recognized audio")
		out <- AsMessage(requestCtx, text)
	*/
}

/*

{"level":"debug","name":"Goldylox","unitID":0,"frequencies":[],"time":"2024-10-17T16:58:42-04:00","message":"synced with SRS client \"Goldylox\""}
{"level":"debug","name":"player","unitID":16783872,"frequencies":[],"time":"2024-10-17T17:04:19-04:00","message":"synced with SRS client \"player\""}
{"level":"debug","name":"player","unitID":16783872,"frequencies":["305"],"time":"2024-10-17T17:05:02-04:00","message":"synced with SRS client \"player\""}
{"level":"debug","name":"player","unitID":16783872,"frequencies":["249.5","305"],"time":"2024-10-17T17:07:10-04:00","message":"synced with SRS client \"player\""}
{"level":"info","duration":"2.04s","time":"2024-10-17T17:07:17-04:00","message":"received transmission"}
{"level":"info","clientName":"player","len":32640,"time":"2024-10-17T17:07:17-04:00","message":"publishing received audio to receiving channel"}

*/

func main() {
	config := types.ClientConfiguration{
		Address:                   "192.168.4.70:5002",
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

	a := &atcApplication{
		recognizer:                 recognizer,
		enableTranscriptionLogging: true,
		transcribedMessages:        make(chan Message[string]),
	}

	log.Info().Msgf("config: %v", config)

	srsClient, err := simpleradio.NewClient(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create SRS client")
		return
	}

	go a.srsLoop(srsClient)
	var wg sync.WaitGroup

	log.Info().Msgf("running")

	atcRecognizer, ok := recognizer.(*deepgramRecognizer.AtcDeepgramRecognizer)
	if !ok {
		log.Error().Msg("Failed to cast recognizer to AtcDeepgramRecognizer")
		return
	}

	text, err := atcRecognizer.Debug_ReadFromWavFile("voice test.wav")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read WAV file")
		return
	}

	log.Info().Msgf("recognized text: %s", text)

	srsClient.Run(
		context.Background(),
		&wg,
	)

	log.Info().Msgf("done")
}
