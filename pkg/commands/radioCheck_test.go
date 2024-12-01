package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/dharmab/skyeye/pkg/simpleradio/voice"
)

func TestRadioCheck(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name         string
		input        string
		clientName   string
		shouldMatch  bool
		expectedText string
		randValues   []int
	}{
		{
			name:         "basic radio check",
			input:        "radio check",
			clientName:   "test1",
			expectedText: "test1, got you loud and clear",
			shouldMatch:  true,
			randValues:   []int{1, 0},
		},
		{
			name:         "radio check with prefix",
			input:        "tower radio check",
			clientName:   "test2",
			expectedText: "good morning, test2. read you five by five",
			shouldMatch:  true,
			randValues:   []int{0, 2},
		},
		{
			name:        "unrelated message",
			input:       "request taxi",
			clientName:  "test3",
			shouldMatch: false,
		},
	}

	rand := &MockGenerator{}

	cp := NewCommandProcessor(rand)
	cp.RegisterParser(&RadioCheckParser{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.ResetTo(tt.randValues)

			msg := message.Message[string]{
				Context:    context.Background(),
				TraceId:    "trace1",
				ClientName: tt.clientName,
				Frequencies: []voice.Frequency{
					{
						Frequency:  123,
						Modulation: 1,
						Encryption: 0,
					},
				},
				Data: tt.input,
			}
			result, err := cp.ProcessText(context.Background(), &msg)

			outChan := make(chan message.OutgoingMessage, 10)

			if tt.shouldMatch {
				assert.Nil(err, "Expected match but got error: %v", err)

				err := result.ParsedCommand.Execute(nil, outChan)
				assert.Nil(err, "Failed to execute command: %v", err)

				select {
				case response := <-outChan:
					assert.NotNil(response, "Expected non-empty response")
					assert.Equal(tt.expectedText, response.Message.Data, "Expected reply text to match")
				default:
					assert.Fail("no response generated")
				}
			} else {
				assert.NotNil(err, "Expected error for non-matching input")
			}
		})
	}
}
