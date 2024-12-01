package atcclient

import (
	"context"
	"encoding/binary"
	"errors"
	"time"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/atcmodel"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/commands"
	deepgramspeaker "github.com/ErikGoldman/DCSAtcOverhaul/pkg/deepgramSpeaker"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/dharmab/skyeye/pkg/recognizer"
	"github.com/dharmab/skyeye/pkg/sim"
	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/telemetry"
	"github.com/rs/zerolog/log"
)

type AtcApplication struct {
	Recognizer                 recognizer.Recognizer
	SpeechSynthesizer          deepgramspeaker.TextToSpeech
	CommandProcessor           commands.CommandProcessorInterface
	EnableTranscriptionLogging bool
	TelemetryClient            telemetry.Client
	AtcModel                   atcmodel.AtcModel

	incomingPlayerCommands chan<- atcmodel.AtcCommand

	simStarted chan<- sim.Started
	simUpdated chan<- sim.Updated
	simFaded   chan<- sim.Faded

	TranscribedMessages chan message.Message[string]
	OutgoingMessages    chan message.OutgoingMessage

	stopCtx      context.Context
	stopCancelFn context.CancelFunc
}

func (a *AtcApplication) srsLoop(radioClient simpleradio.Client) {
	for {
		select {
		case <-a.stopCtx.Done():
			log.Info().Msg("srsLoop stopping")
			return

		case transmission := <-radioClient.Receive():
			if len(transmission.Frequencies) == 0 {
				log.Error().Msg("got transmission without a frequency")
				continue
			}

			frequencies := transmission.Frequencies
			log.Info().Msgf("received transmission from frequency %s", &frequencies[0])

			a.recognizeTransmission(context.Background(), nil, transmission, a.TranscribedMessages)
		}
	}
}

func (a *AtcApplication) processTranscriptLoop(radioClient simpleradio.Client) {
	for {
		select {
		case <-a.stopCtx.Done():
			log.Info().Msg("transcript loop stopping")
			return

		case msg := <-a.TranscribedMessages:
			log.Info().Msg("processing transcription")
			cmd, err := a.CommandProcessor.ProcessText(context.Background(), &msg)
			if err == nil {
				log.Info().Msgf("sending command to ATC %s", cmd)
				a.incomingPlayerCommands <- cmd.ParsedCommand
			} else {
				log.Info().Msgf("command parsing failed")
			}
		}
	}
}

func convertLinear16ToFloat32(linear16 []byte) []float32 {
	// Calculate the number of samples
	numSamples := len(linear16) / 2
	float32Data := make([]float32, numSamples)

	for i := 0; i < numSamples; i++ {
		// Read int16 from byte array
		sample := int16(binary.LittleEndian.Uint16(linear16[i*2 : (i+1)*2]))
		// Normalize to float32
		float32Data[i] = float32(sample) / 32768.0 // 32768.0 is the max value for int16
	}

	return float32Data
}

// do these in serial for now. in theory we can do it in parallel across multiple frequencies
// but for now block and have just one ATC thread globally (I think it's fine)
func (a *AtcApplication) processOutgoingAudioLoop(radioClient simpleradio.Client) {
	for {
		select {
		case <-a.stopCtx.Done():
			log.Info().Msg("outgoing audio loop stopping")
			return

		case msg := <-a.OutgoingMessages:
			log.Info().Msg("processing outgoing message")

			audioChannel := make(chan []byte, 5)

			// use text to speech API and pipe results into the radio client
			err := a.SpeechSynthesizer.GenerateSpeech(msg.Model, msg.Message.Data, audioChannel)
			log.Info().Msg("done generating speech")
			if err != nil {
				log.Error().Err(err).Msg("error generating speech")
				continue
			}

			func() {
				for {
					select {
					case audioBytes := <-audioChannel:
						if audioBytes == nil {
							log.Info().Msg("got end of TTS stream")
							return
						}

						log.Info().Msg("sending voice transmission")
						transmission := simpleradio.Transmission{
							TraceID:     msg.Message.TraceId,
							ClientName:  msg.Message.ClientName,
							Frequencies: msg.Message.Frequencies,
							Audio:       convertLinear16ToFloat32(audioBytes),
						}
						radioClient.Transmit(transmission)
					}
				}
			}()

			a.SpeechSynthesizer.Disconnect()
		}
	}
}

func (a *AtcApplication) recognizeTransmission(processCtx context.Context, requestCtx context.Context,
	transmission simpleradio.Transmission, out chan<- message.Message[string]) {

	recogizerCtx, cancel := context.WithTimeout(processCtx, 30*time.Second)
	defer func() {
		if recogizerCtx.Err() != nil && errors.Is(recogizerCtx.Err(), context.DeadlineExceeded) {
			log.Error().Msg("timeout processing speech")
		}
	}()
	defer cancel()

	log.Info().Msg("recognizing audio")
	//start := time.Now()
	text, err := a.Recognizer.Recognize(recogizerCtx, transmission.Audio, a.EnableTranscriptionLogging)
	if err != nil {
		log.Error().Err(err).Msg("error recognizing audio sample")
		return
	}

	if a.EnableTranscriptionLogging {
		log.Info().Msgf("recognized text: %s", text)
	}
	out <- message.FromTransmission(requestCtx, transmission, text)

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

func LoadCommandProcessor() *commands.CommandProcessor {
	rand := commands.RealGenerator{}
	cp := commands.NewCommandProcessor(&rand)

	cp.RegisterParser(&commands.RadioCheckParser{})

	return cp
}

func (a *AtcApplication) Start(srsClient simpleradio.Client) {
	a.stopCtx, a.stopCancelFn = context.WithCancel(context.Background())
	a.TranscribedMessages = make(chan message.Message[string], 5)
	a.OutgoingMessages = make(chan message.OutgoingMessage, 5)

	go a.TelemetryClient.Run(a.stopCtx, nil)
	go func() {
		log.Info().Msg("streaming telemetry data")
		a.TelemetryClient.Stream(a.stopCtx, nil, a.simStarted, a.simUpdated, a.simFaded)
	}()

	go a.processOutgoingAudioLoop(srsClient)
	go a.processTranscriptLoop(srsClient)
	a.srsLoop(srsClient)
}

func (a *AtcApplication) Stop() {
	a.stopCancelFn()
}
