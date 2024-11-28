package deepgramspeaker

import (
	"context"
	"fmt"
	"strings"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	speak "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
	websocketv1 "github.com/deepgram/deepgram-go-sdk/pkg/client/speak/v1/websocket"
	"github.com/rs/zerolog/log"
)

type TextToSpeech interface {
	GenerateSpeech(model string, text string, out chan []byte) error
	Disconnect() error
}

// Implement your own callback
type MyCallback struct {
	out chan []byte
}

func (c MyCallback) Metadata(md *msginterfaces.MetadataResponse) error {
	fmt.Printf("\n[Metadata] Received\n")
	fmt.Printf("Metadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	return nil
}

func (c MyCallback) Binary(byMsg []byte) error {
	c.out <- byMsg
	return nil
}

func (c MyCallback) Flush(fl *msginterfaces.FlushedResponse) error {
	fmt.Printf("\n[Flushed] Received\n")
	c.out <- nil
	return nil
}

func (c MyCallback) Warning(wr *msginterfaces.WarningResponse) error {
	fmt.Printf("\n[Warning] Received\n")
	fmt.Printf("Warning.Code: %s\n", wr.WarnCode)
	fmt.Printf("Warning.Description: %s\n\n", wr.WarnMsg)
	return nil
}

func (c MyCallback) Error(er *msginterfaces.ErrorResponse) error {
	fmt.Printf("\n[Error] Received\n")
	fmt.Printf("Error.Code: %s\n", er.ErrCode)
	fmt.Printf("Error.Description: %s\n\n", er.ErrMsg)
	return nil
}

func (c MyCallback) Close(cr *msginterfaces.CloseResponse) error {
	fmt.Printf("\n[Close] Received\n")
	c.out <- nil
	return nil
}

func (c MyCallback) Open(or *msginterfaces.OpenResponse) error {
	fmt.Printf("\n[Open] Received\n")
	return nil
}

func (c MyCallback) Clear(cr *msginterfaces.ClearedResponse) error {
	fmt.Printf("\n[Clear] Received\n")
	return nil
}

func (c MyCallback) UnhandledEvent([]byte) error {
	fmt.Printf("\n[UnhandledEvent] Received\n")
	return nil
}

type DeepgramSpeakSynthesizer struct {
	dgClient *websocketv1.WSCallback
	apiKey   string
}

func NewSpeechSynthesizer(apiKey string) *DeepgramSpeakSynthesizer {
	//speak.InitWithDefault()
	return &DeepgramSpeakSynthesizer{
		dgClient: nil,
		apiKey:   apiKey,
	}
}

func (d *DeepgramSpeakSynthesizer) GenerateSpeech(model string, text string, out chan []byte) error {
	ttsOptions := &interfaces.WSSpeakOptions{
		Model:      model,
		Encoding:   "linear16",
		SampleRate: 24000,
	}

	ctx := context.Background()
	callback := MyCallback{
		out: out,
	}

	dgClient, err := speak.NewWSUsingCallback(ctx, d.apiKey, &interfaces.ClientOptions{}, ttsOptions, callback)
	d.dgClient = dgClient
	if err != nil {
		fmt.Println("ERROR creating TTS connection:", err)
		return err
	}

	bConnected := dgClient.Connect()
	if !bConnected {
		fmt.Println("Client.Connect failed")
		return err
	}

	err = dgClient.SpeakWithText(text)
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return err
	}

	err = dgClient.Flush()
	if err != nil {
		fmt.Printf("Error sending text input: %v\n", err)
		return err
	}

	return nil
}

func (d *DeepgramSpeakSynthesizer) Disconnect() error {
	log.Info().Msg("disconnecting from deepgram TTS")
	d.dgClient.Stop()
	return nil
}
