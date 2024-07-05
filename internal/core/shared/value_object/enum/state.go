package enum

type State string

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
