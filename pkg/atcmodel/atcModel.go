package atcmodel

import (
	"context"
	"fmt"

	"github.com/ErikGoldman/DCSAtcOverhaul/pkg/message"
	"github.com/dharmab/skyeye/pkg/sim"
	"github.com/dharmab/skyeye/pkg/simpleradio/types"
	"github.com/dharmab/skyeye/pkg/spatial"
	"github.com/martinlindhe/unit"
	"github.com/rs/zerolog/log"
)

const IS_AIRBORN_AGL = 10 * unit.Meter

type AtcModel struct {
	Map          AtcMap
	Squads       map[types.Radio][]*AtcSquadron
	PlaneToSquad map[uint64]*AtcSquadron

	AllPlaneData map[uint64]*sim.Updated
	CallsignToId map[string]*uint64
}

// needed to avoid circular dependencies with parsed commands
type AtcCommand interface {
	Execute(atc *AtcModel, messageOut chan message.OutgoingMessage) error
}

func (a *AtcModel) Start(ctx context.Context, simStarted chan sim.Started, simUpdated chan sim.Updated, simFaded chan sim.Faded, commands chan AtcCommand,
	messageOut chan message.OutgoingMessage) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("atc loop exiting")
			return

		case <-simStarted:
			// new mission started
			log.Info().Msg("atc model notified of mission change")
			a.reset()

		case updated := <-simUpdated:
			if squad, ok := a.PlaneToSquad[updated.Labels.ID]; ok {
				if plane, okPlane := squad.PlaneStates[updated.Labels.ID]; okPlane {
					plane.UpdateFromTrack(updated.Frame)
				} else {
					log.Warn().Msgf("could not find plane %d (%s) in squad", updated.Labels.ID, updated.Labels.Name)
				}
			}
			a.AllPlaneData[updated.Labels.ID] = &updated
			if _, ok := a.CallsignToId[updated.Labels.Name]; !ok {
				log.Info().Msgf("added callsign mapping %s -> %d", updated.Labels.Name, updated.Labels.ID)
				a.CallsignToId[updated.Labels.Name] = &updated.Labels.ID
			}

		case removed := <-simFaded:
			if squad, ok := a.PlaneToSquad[removed.ID]; ok {
				log.Info().Msgf("removing plane %d from squad due to disconnection", removed.ID)
				delete(squad.PlaneStates, removed.ID)
				delete(a.PlaneToSquad, removed.ID)
			}
			if planeData, ok := a.AllPlaneData[removed.ID]; ok {
				log.Info().Msgf("removing plane %d from records due to disconnection", removed.ID)
				delete(a.CallsignToId, planeData.Labels.Name)
				delete(a.AllPlaneData, removed.ID)
			}

		case cmd := <-commands:
			log.Info().Msgf("atc executing command %s", cmd)
			cmd.Execute(a, messageOut)
		}
	}
}

type SquadSearchResult struct {
	PlaneId          uint64
	PlaneType        string
	Distance         int
	IsAlreadyInSquad bool
	AGL              unit.Length
}

// if someone says "squad of two F16s" -- which two F16s?
func (a *AtcModel) FindCandidatesForSquad(ctx context.Context, leaderId uint64,
	squadmatePlanes []string, maxDistance unit.Length) ([]SquadSearchResult, error) {

	var leaderData *sim.Updated
	leaderData, foundLeader := a.AllPlaneData[leaderId]
	if !foundLeader {
		return nil, fmt.Errorf("leader with id %d does not exist", leaderId)
	}

	planeTypes := make(map[string]struct{})
	for _, str := range squadmatePlanes {
		planeTypes[str] = struct{}{}
	}

	var candidatePlanes = []SquadSearchResult{}
	for id, planeData := range a.AllPlaneData {
		if id == leaderId {
			continue
		}

		if _, ok := planeTypes[planeData.Labels.ACMIName]; !ok {
			continue
		}

		distance := spatial.Distance(leaderData.Frame.Point, planeData.Frame.Point)
		if distance > maxDistance {
			continue
		}

		_, isAlreadyInSquad := a.PlaneToSquad[id]

		candidatePlanes = append(candidatePlanes, SquadSearchResult{
			PlaneId:          planeData.Labels.ID,
			PlaneType:        planeData.Labels.ACMIName,
			Distance:         int(distance),
			IsAlreadyInSquad: isAlreadyInSquad,
			AGL:              *planeData.Frame.AGL,
		})
	}

	return candidatePlanes, nil
}

func (a *AtcModel) AddSquadron(planes []uint64) {

}

func (a *AtcModel) reset() {
	log.Info().Msg("resetting atc model")
	for r := range a.Squads {
		delete(a.Squads, r)
	}
}
