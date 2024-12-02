package atcmodel

import (
	"context"
	"testing"

	"github.com/dharmab/skyeye/pkg/bearings"
	"github.com/dharmab/skyeye/pkg/coalitions"
	"github.com/dharmab/skyeye/pkg/sim"
	"github.com/dharmab/skyeye/pkg/spatial"
	"github.com/dharmab/skyeye/pkg/trackfiles"
	"github.com/martinlindhe/unit"
	"github.com/paulmach/orb"
)

func TestFindCandidatesForSquad(t *testing.T) {
	tests := []struct {
		name          string
		setupModel    func() *AtcModel
		leaderId      uint64
		planeTypes    []string
		maxDistance   unit.Length
		wantResults   []SquadSearchResult
		wantErr       bool
		errorContains string
	}{
		{
			name: "ignores leader as candidate",
			setupModel: func() *AtcModel {
				model := &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}

				// Leader F-16
				model.AllPlaneData[1] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       1,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: orb.Point{45, 45},
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				return model
			},
			leaderId:    1,
			planeTypes:  []string{"F-16C_50"},
			maxDistance: 5000 * unit.Meter,
			wantResults: []SquadSearchResult{},
			wantErr:     false,
		},

		{
			name: "finds nearby aircraft of correct type",
			setupModel: func() *AtcModel {
				model := &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}

				// Leader F-16
				model.AllPlaneData[1] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       1,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: orb.Point{45, 45},
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-16
				model.AllPlaneData[2] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       2,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*1000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Far F-16
				model.AllPlaneData[3] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       3,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*8000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-15 (wrong type)
				model.AllPlaneData[4] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       4,
						ACMIName: "F-15C",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*1000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				return model
			},
			leaderId:    1,
			planeTypes:  []string{"F-16C_50"},
			maxDistance: 5000 * unit.Meter,
			wantResults: []SquadSearchResult{
				{
					PlaneId:          2,
					PlaneType:        "F-16C_50",
					Distance:         1000,
					IsAlreadyInSquad: false,
					AGL:              1000 * unit.Meter,
				},
			},
			wantErr: false,
		},
		{
			name: "handles aircraft already in squad",
			setupModel: func() *AtcModel {
				model := &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}

				// Leader
				model.AllPlaneData[1] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       1,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: orb.Point{45, 45},
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-16 already in squad
				model.AllPlaneData[2] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       2,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*1000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Mark plane 2 as already in a squad
				model.PlaneToSquad[2] = &AtcSquadron{}

				return model
			},
			leaderId:    1,
			planeTypes:  []string{"F-16C_50"},
			maxDistance: 5000 * unit.Meter,
			wantResults: []SquadSearchResult{
				{
					PlaneId:          2,
					PlaneType:        "F-16C_50",
					Distance:         1000,
					IsAlreadyInSquad: true,
					AGL:              1000 * unit.Meter,
				},
			},
			wantErr: false,
		},
		{
			name: "ignore if different coalition",
			setupModel: func() *AtcModel {
				model := &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}

				model.AllPlaneData[1] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:        1,
						ACMIName:  "F-16C_50",
						Coalition: coalitions.Blue,
					},
					Frame: trackfiles.Frame{
						Point: orb.Point{45, 45},
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-16 already in squad
				model.AllPlaneData[2] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:        2,
						ACMIName:  "F-16C_50",
						Coalition: coalitions.Red,
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*1000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				return model
			},
			leaderId:      1,
			planeTypes:    []string{"F-16C_50"},
			maxDistance:   5000 * unit.Meter,
			wantResults:   []SquadSearchResult{},
			wantErr:       false,
			errorContains: "",
		},
		{
			name: "leader not found",
			setupModel: func() *AtcModel {
				return &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}
			},
			leaderId:      1,
			planeTypes:    []string{"F-16C_50"},
			maxDistance:   5000 * unit.Meter,
			wantResults:   nil,
			wantErr:       true,
			errorContains: "does not exist",
		},
		{
			name: "multiple plane types",
			setupModel: func() *AtcModel {
				model := &AtcModel{
					AllPlaneData: make(map[uint64]*sim.Updated),
					PlaneToSquad: make(map[uint64]*AtcSquadron),
				}

				// Leader
				model.AllPlaneData[1] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       1,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: orb.Point{45, 45},
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-16
				model.AllPlaneData[2] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       2,
						ACMIName: "F-16C_50",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*1000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				// Nearby F-15
				model.AllPlaneData[3] = &sim.Updated{
					Labels: trackfiles.Labels{
						ID:       3,
						ACMIName: "F-15C",
					},
					Frame: trackfiles.Frame{
						Point: spatial.PointAtBearingAndDistance(orb.Point{45, 45}, bearings.NewTrueBearing(10), unit.Meter*2000),
						AGL:   ptrlength(1000 * unit.Meter),
					},
				}

				return model
			},
			leaderId:    1,
			planeTypes:  []string{"F-16C_50", "F-15C"},
			maxDistance: 5000 * unit.Meter,
			wantResults: []SquadSearchResult{
				{
					PlaneId:          2,
					PlaneType:        "F-16C_50",
					Distance:         1000,
					IsAlreadyInSquad: false,
					AGL:              1000 * unit.Meter,
				},
				{
					PlaneId:          3,
					PlaneType:        "F-15C",
					Distance:         2000,
					IsAlreadyInSquad: false,
					AGL:              1000 * unit.Meter,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			ctx := context.Background()

			results, err := model.FindCandidatesForSquad(ctx, tt.leaderId, tt.planeTypes, tt.maxDistance)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q but got nil", tt.errorContains)
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q but got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(results) != len(tt.wantResults) {
				t.Errorf("got %d results, want %d", len(results), len(tt.wantResults))
				return
			}

			for i, result := range results {
				want := tt.wantResults[i]
				if result.PlaneId != want.PlaneId ||
					result.PlaneType != want.PlaneType ||
					result.Distance != want.Distance ||
					result.IsAlreadyInSquad != want.IsAlreadyInSquad ||
					result.AGL != want.AGL {
					t.Errorf("result %d mismatch:\ngot: %+v\nwant: %+v", i, result, want)
				}
			}
		})
	}
}

// Helper functions
func ptrlength(l unit.Length) *unit.Length {
	return &l
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}
