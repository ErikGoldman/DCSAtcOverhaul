package atcmodel

import (
	"github.com/dharmab/skyeye/pkg/trackfiles"
)

type PlaneState interface {
	UpdateFromTrack(update trackfiles.Frame)
	TransitionToState()
}

type AtcSquadron struct {
	PlaneStates map[uint64]PlaneState
}
