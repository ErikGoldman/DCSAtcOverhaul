package atcmodel

import (
	"context"
	"fmt"

	"github.com/dharmab/skyeye/pkg/sim"
	"github.com/dharmab/skyeye/pkg/spatial"
	"github.com/martinlindhe/unit"
	"github.com/rs/zerolog/log"
)

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

func (a *AtcModel) DoesExistingSquadMatchTypes(leader uint64, planeTypes *[]string) bool {
	squad, hasSquad := a.PlaneToSquad[leader]
	if !hasSquad {
		return false
	}

	planeTypeCount := make(map[string]int)
	for _, planeType := range *planeTypes {
		planeTypeCount[planeType]++
	}

	for planeId := range squad.PlaneStates {
		if planeId == leader {
			continue
		}

		planeInfo, exists := a.AllPlaneData[planeId]
		if !exists {
			log.Error().Msgf("Could not find %d when looking for existing squad for %d", planeId, leader)
			continue
		}

		if planeCount, exists := planeTypeCount[planeInfo.Labels.ACMIName]; exists {
			planeTypeCount[planeInfo.Labels.ACMIName] = planeCount - 1
		}
	}

	for _, planeCount := range planeTypeCount {
		if planeCount != 0 {
			return false
		}
	}
	return true
}
