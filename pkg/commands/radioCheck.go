package commands

import (
	"fmt"
	"strings"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
)

// "Anapa tower, Uzi 2-1, radio check"

type RadioCheckParser struct {
}

type RadioCheck struct {
	Message       *message.Message[string]
	globalContext *GlobalCommandContext
}

func (m *RadioCheck) String() string {
	return "RadioCheckCommand"
}

func (m *RadioCheck) Execute() (string, error) {
	fmt.Printf("Radio check %s", m.Message.ClientName)

	intro := ""
	if m.globalContext.rand.IntN(3) == 0 {
		if m.Message.GameTimeHour <= 10 {
			intro = "morning"
		} else if m.Message.GameTimeHour <= 16 {
			intro = "afternoon"
		} else {
			intro = "evening"
		}
		intro = fmt.Sprintf("good %s, %s.", intro, m.Message.ClientName)
	} else {
		intro = fmt.Sprintf("%s,", m.Message.ClientName)
	}

	bodyText := ""
	variation := m.globalContext.rand.IntN(5)
	if variation == 0 {
		bodyText = "got you loud and clear"
	} else if variation == 1 {
		bodyText = "loud and clear"
	} else if variation == 2 {
		bodyText = "read you five by five"
	} else if variation == 3 {
		bodyText = "lima charlie"
	} else if variation == 4 {
		bodyText = "reading you loud and clear"
	}

	return fmt.Sprintf("%s %s", intro, bodyText), nil
}

func (p *RadioCheckParser) Parse(globalContext *GlobalCommandContext, message *message.Message[string]) PlayerCommand {
	if strings.Contains(message.Data, "radio check") {
		return &RadioCheck{globalContext: globalContext, Message: message}
	}
	return nil
}
