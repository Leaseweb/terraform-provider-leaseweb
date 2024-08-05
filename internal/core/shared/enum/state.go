package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type State string

func (s State) String() string {
	return string(s)
}

const (
	StateCreating   State = "CREATING"
	StateDestroyed  State = "DESTROYED"
	StateDestroying State = "DESTROYING"
	StateFailed     State = "FAILED"
	StateRunning    State = "RUNNING"
	StateStarting   State = "STARTING"
	StateStopped    State = "STOPPED"
	StateStopping   State = "STOPPING"
	StateUnknown    State = "UNKNOWN"
)

var states = []State{
	StateCreating,
	StateDestroyed,
	StateDestroying,
	StateFailed,
	StateRunning,
	StateStarting,
	StateStopped,
	StateStopping,
	StateUnknown,
}

func NewState(s string) (State, error) {
	return enum_utils.FindEnumForString(s, states, StateUnknown)
}
