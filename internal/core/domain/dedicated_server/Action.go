package dedicated_server

type Action struct {
	Type            string
	LastTriggeredAt *string
}

type OptionalActionValues struct {
	LastTriggeredAt *string
}

func NewAction(actionType string, optional OptionalActionValues) Action {
	action := Action{Type: actionType}

	if optional.LastTriggeredAt != nil {
		action.LastTriggeredAt = optional.LastTriggeredAt
	}

	return action
}
