package atcclienttesthelpers

import (
	"context"
	"sync"

	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/types"
	"github.com/stretchr/testify/mock"
)

type MockSRSClient struct {
	mock.Mock
}

func (m *MockSRSClient) Run(ctx context.Context, wg *sync.WaitGroup) error {
	args := m.Called(ctx, wg)
	return args.Error(0)
}

func (m *MockSRSClient) Send(msg types.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockSRSClient) Receive() <-chan simpleradio.Transmission {
	args := m.Called()
	return args.Get(0).(chan simpleradio.Transmission)
}

func (m *MockSRSClient) Transmit(transmission simpleradio.Transmission) {
	m.Called(transmission)
}

func (m *MockSRSClient) Frequencies() []simpleradio.RadioFrequency {
	args := m.Called()
	return args.Get(0).([]simpleradio.RadioFrequency)
}

func (m *MockRecognizer) Recognize(ctx context.Context, pcm []float32, enableTranscriptionLogging bool) (string, error) {
	args := m.Called(ctx, pcm, enableTranscriptionLogging)
	return args.String(0), args.Error(1)
}

func (m *MockSRSClient) ClientsOnFrequency() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockSRSClient) HumansOnFrequency() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockSRSClient) BotsOnFrequency() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockSRSClient) IsOnFrequency(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

type MockRecognizer struct {
	mock.Mock
}

func IsAudioEqual(a []float32, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
