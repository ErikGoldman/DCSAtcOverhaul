package atcclienttest

import (
	"testing"

	atcclient "github.com/ErikGoldman/DCSAtcOverhaul/pkg/client"
	atcclienttesthelpers "github.com/ErikGoldman/DCSAtcOverhaul/pkg/client/testhelpers"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/voice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSRSLoop_SimpleSingleFrequency(t *testing.T) {
	mockClient := &atcclienttesthelpers.MockSRSClient{}
	mockRecognizer := &atcclienttesthelpers.MockRecognizer{}

	app := &atcclient.AtcApplication{
		Recognizer:                 mockRecognizer,
		EnableTranscriptionLogging: false,
	}

	// EXPECTATIONS
	transmissionChan := make(chan simpleradio.Transmission)
	mockClient.On("Receive").Return(transmissionChan)

	mockRecognizer.On("Recognize",
		mock.Anything, mock.Anything, mock.Anything).Return("recognized text", nil)

	// FUNCTION
	go app.Start(mockClient)
	transmissionChan <- simpleradio.Transmission{
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
	app.Stop()

	// ASSERTIONS
	mockClient.AssertExpectations(t)

	mockRecognizer.AssertNumberOfCalls(t, "Recognize", 1)
	mockRecognizer.AssertCalled(t, "Recognize",
		mock.Anything, // context
		mock.MatchedBy(func(pcm []float32) bool {
			expectedPCM := []float32{1, 2, 3, 4, 10}
			return atcclienttesthelpers.IsAudioEqual(expectedPCM, pcm)
		}),
		mock.MatchedBy(func(transcriptLoggingEnabled bool) bool {
			return !transcriptLoggingEnabled
		}),
	)

	var msg message.Message[string]
	select {
	case msg = <-app.TranscribedMessages:
	default:
		assert.Fail(t, "Expected one message on TranscribedMessages, but got none")
		return
	}

	assert.EqualValues(t, msg.Data, "recognized text")
	assert.Len(t, msg.Frequencies, 1)
	assert.EqualValues(t, msg.Frequencies[0].Frequency, 123.4)
	assert.EqualValues(t, msg.Frequencies[0].Modulation, 2)
	assert.EqualValues(t, msg.Frequencies[0].Encryption, 1)
	assert.EqualValues(t, msg.ClientName, "MyClientName")
	assert.EqualValues(t, msg.TraceId, "MyTraceId")

	select {
	case msgTwo := <-app.TranscribedMessages:
		assert.Fail(t, "Expected one message on TranscribedMessages, but got %s", msgTwo.Data)
		return
	default:
		// No additional messages, which is expected
	}
}

func TestSRSLoop_SimpleMultiFrequency(t *testing.T) {
	mockClient := &atcclienttesthelpers.MockSRSClient{}
	mockRecognizer := &atcclienttesthelpers.MockRecognizer{}

	app := &AtcApplication{
		Recognizer:                 mockRecognizer,
		EnableTranscriptionLogging: false,
	}

	// EXPECTATIONS
	transmissionChan := make(chan simpleradio.Transmission)
	mockClient.On("Receive").Return(transmissionChan)

	mockRecognizer.On("Recognize",
		mock.Anything, mock.Anything, mock.Anything).Return("recognized text", nil)

	// FUNCTION
	go app.Start(mockClient)
	transmissionChan <- simpleradio.Transmission{
		Frequencies: []voice.Frequency{
			voice.Frequency{
				Frequency:  245.8,
				Modulation: 8,
				Encryption: 0,
			},
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
	app.Stop()

	// ASSERTIONS
	mockClient.AssertExpectations(t)

	mockRecognizer.AssertNumberOfCalls(t, "Recognize", 1)
	mockRecognizer.AssertCalled(t, "Recognize",
		mock.Anything, // context
		mock.Anything, // context
		mock.Anything, // context
	)

	var msg message.Message[string]
	select {
	case msg = <-app.TranscribedMessages:
	default:
		assert.Fail(t, "Expected one message on TranscribedMessages, but got none")
		return
	}

	assert.Len(t, msg.Frequencies, 2)
	assert.EqualValues(t, msg.Frequencies[0].Frequency, 245.8)
	assert.EqualValues(t, msg.Frequencies[1].Frequency, 123.4)

	select {
	case msgTwo := <-app.TranscribedMessages:
		assert.Fail(t, "Expected one message on TranscribedMessages, but got %s", msgTwo.Data)
		return
	default:
		// No additional messages, which is expected
	}
}
