package loader

// StackStatus represents the status of a Docker Compose stack.
type StackStatus string

const (
	// StackStatusRunning indicates the stack is running.
	StackStatusRunning StackStatus = "running"
	// StackStatusStopped indicates the stack is stopped but containers exist.
	StackStatusStopped StackStatus = "stopped"
	// StackStatusDown indicates the stack has no containers.
	StackStatusDown StackStatus = "down"
)

// Stack represents a Docker Compose stack.
type Stack struct {
	Name   string
	Dir    string
	Status StackStatus
}

// Action represents a stack operation.
type Action string

// Action constants for stack operations.
const (
	ActionStart   Action = "start"
	ActionStop    Action = "stop"
	ActionRestart Action = "restart"
	ActionReload  Action = "reload"
	ActionDown    Action = "down"
	ActionList    Action = "list"
)

// IsValid checks if the action is a valid operation.
func (a Action) IsValid() bool {
	switch a {
	case ActionStart, ActionStop, ActionRestart, ActionReload, ActionDown, ActionList:
		return true
	default:
		return false
	}
}
