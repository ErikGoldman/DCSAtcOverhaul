package atcclienttest

import (
	"testing"

	atcclient "github.com/ErikGoldman/DCSAtcOverhaul/pkg/client"
	atcclienttesthelpers "github.com/ErikGoldman/DCSAtcOverhaul/pkg/client/testhelpers"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/commands"
	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/voice"
	"github.com/stretchr/testify/mock"
)

type MockTextToSpeech struct {
	mock.Mock
}

func (m *MockTextToSpeech) GenerateSpeech(model string, text string, out chan []byte) error {
	args := m.Called(model, text, out)
	return args.Error(0)
}

func (m *MockTextToSpeech) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func TestClient_EndToEnd(t *testing.T) {
	mockClient := &atcclienttesthelpers.MockSRSClient{}
	mockRecognizer := &atcclienttesthelpers.MockRecognizer{}
	mockTextToSpeech := &MockTextToSpeech{}

	commandProcessor := commands.NewCommandProcessor(&commands.RealGenerator{})
	commandProcessor.RegisterParser(&commands.RadioCheckParser{})

	app := &atcclient.AtcApplication{
		Recognizer:                 mockRecognizer,
		SpeechSynthesizer:          mockTextToSpeech,
		CommandProcessor:           commandProcessor,
		EnableTranscriptionLogging: true,
	}

	// EXPECTATIONS
	recieveChannel := make(chan simpleradio.Transmission)
	mockClient.On("Receive").Return(recieveChannel)

	blockerChannel := make(chan bool)

	mockRecognizer.On("Recognize",
		mock.Anything, mock.Anything, mock.Anything).Return("kutaisi alpha one one radio check", nil)
	mockTextToSpeech.On("GenerateSpeech", mock.Anything, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			byteChan := args.Get(2).(chan []byte)
			byteChan <- []byte{1, 2, 3, 4}
			byteChan <- []byte{5, 6, 7, 8}
			byteChan <- nil
		},
	).Return(nil)
	mockTextToSpeech.On("Disconnect").Return(nil)

	numTransmitCalls := 0
	mockClient.On("Transmit", mock.Anything).Run(func(args mock.Arguments) {
		numTransmitCalls++
		if numTransmitCalls == 2 {
			app.Stop()
			blockerChannel <- true
		}
	})

	// FUNCTION
	go app.Start(mockClient)
	recieveChannel <- simpleradio.Transmission{
		Frequencies: []voice.Frequency{
			voice.Frequency{
				Frequency:  123.4,
				Modulation: 2,
				Encryption: 1,
			},
		},
		TraceID:    "MyTraceId",
		ClientName: "MyClientName",
		Audio:      []float32{1, 2, 3, 4, 10},
	}

	select {
	case _ = <-blockerChannel:
		break
	}

	// ASSERTIONS
	mockClient.AssertExpectations(t)

	mockClient.AssertCalled(t, "Transmit", mock.MatchedBy(func(tr simpleradio.Transmission) bool {
		if tr.ClientName != "MyClientName" || tr.TraceID != "MyTraceId" {
			return false
		}

		if len(tr.Frequencies) != 1 {
			return false
		}

		if len(tr.Audio) != 2 {
			return false
		}
		return true
	}))

	mockTextToSpeech.AssertNumberOfCalls(t, "Disconnect", 1)
}
