package relay

import (
	"sync/atomic"

	"github.com/blocknative/dreamboat/pkg/structs"
)

type AtomicSharedState struct {
	duties     atomic.Value
	validators atomic.Value
	genesis    atomic.Value
}

func (as *AtomicSharedState) Beacon() *structs.BeaconState {
	state := &structs.BeaconState{}

	if val := as.duties.Load(); val != nil {
		state.DutiesState = val.(structs.DutiesState)
	}
	if val := as.validators.Load(); val != nil {
		state.ValidatorsState = val.(structs.ValidatorsState)
	}
	if val := as.genesis.Load(); val != nil {
		state.GenesisInfo = val.(structs.GenesisInfo)
	}

	return state
}
