package message

import (
	"context"

	"github.com/dharmab/skyeye/pkg/simpleradio"
	"github.com/dharmab/skyeye/pkg/simpleradio/voice"
)

type Message[T any] struct {
	Context     context.Context
	TraceId     string
	ClientName  string
	Data        T
	Frequencies []voice.Frequency

	GameTimeHour   int
	GameTimeMinute int
	GameTimeSecond int
}

func FromMessage[T any](ctx context.Context, msg Message[T], data T) Message[T] {
	return Message[T]{Context: ctx, TraceId: msg.TraceId,
		Frequencies: msg.Frequencies, ClientName: msg.ClientName, Data: data,
		GameTimeHour: msg.GameTimeHour, GameTimeMinute: msg.GameTimeMinute,
		GameTimeSecond: msg.GameTimeSecond}
}

func FromTransmission[T any](ctx context.Context, transmission simpleradio.Transmission, data T) Message[T] {
	return Message[T]{Context: ctx, TraceId: transmission.TraceID,
		Frequencies: transmission.Frequencies, ClientName: transmission.ClientName, Data: data}
}

type OutgoingMessage struct {
	Message Message[string]
	Model   string
}
