package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/atcmodel"
	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/rs/zerolog/log"
)

type PlayerCommand interface {
	Execute(atc *atcmodel.AtcModel, messageOut chan message.OutgoingMessage) error
}

// PlayerCommand represents a parsed command
type PlayerCommandMessage struct {
	Message       *message.Message[string]
	ParsedCommand PlayerCommand
}

type GlobalCommandContext struct {
	rand Random
}

// PlayerCommandParser defines the interface for command parsers
type PlayerCommandParser interface {
	Parse(globalContext *GlobalCommandContext, message *message.Message[string]) PlayerCommand
}

type CommandProcessorInterface interface {
	ProcessText(ctx context.Context, message *message.Message[string]) (PlayerCommandMessage, error)
}

// CommandProcessor handles the registration and processing of command parsers
type CommandProcessor struct {
	parsers       []PlayerCommandParser
	globalContext *GlobalCommandContext
}

func NewCommandProcessor(rand Random) *CommandProcessor {
	return &CommandProcessor{
		parsers: make([]PlayerCommandParser, 0),
		globalContext: &GlobalCommandContext{
			rand: rand,
		},
	}
}

// RegisterParser adds a new parser to the processor
func (cp *CommandProcessor) RegisterParser(parser PlayerCommandParser) {
	cp.parsers = append(cp.parsers, parser)
}

// ProcessCommand attempts to parse the input string using registered parsers
func (cp *CommandProcessor) ProcessText(ctx context.Context, message *message.Message[string]) (PlayerCommandMessage, error) {
	message.Data = strings.ToLower(message.Data)

	for _, parser := range cp.parsers {
		if cmd := parser.Parse(cp.globalContext, message); cmd != nil {
			log.Info().Msgf("Matched to command %s", cmd)
			return PlayerCommandMessage{
				Message:       message,
				ParsedCommand: cmd,
			}, nil
		} else {
			log.Info().Msgf("Could not match to command")
		}
	}
	return PlayerCommandMessage{}, fmt.Errorf("no parser found for command: %s", message.Data)
}
