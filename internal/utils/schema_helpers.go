package utils

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	BetaDescription string = "**Warning:** This functionality is in BETA. Documentation might be incorrect or incomplete. Functionality might change with the final release."
)

type Action int

const (
	CreateAction Action = iota
	ReadAction
	UpdateAction
	DeleteAction
)

func (a Action) string(unsupportedActions []Action) string {
	var secondAction string
	var format = "Once %s, this resource cannot be %s"

	switch a {
	case CreateAction:
		secondAction = "created"
	case ReadAction:
		secondAction = "read"
	case UpdateAction:
		secondAction = "updated"
	case DeleteAction:
		secondAction = "deleted"
	default:
		log.Fatal(fmt.Printf("do not know how to handle action: %q", a))
	}

	return fmt.Sprintf(
		format,
		a.firstAction(unsupportedActions),
		secondAction,
	)
}

func (a Action) firstAction(unsupportedActions []Action) string {
	if slices.Contains(unsupportedActions, CreateAction) {
		return "imported"
	}

	return "created"
}

// AddUnsupportedActionsNotation lets the end user know which actions aren't supported in the markdown description.
func AddUnsupportedActionsNotation(
	response *resource.SchemaResponse,
	unsupportedActions []Action,
) {
	response.Schema.MarkdownDescription += "**Note:**"
	for _, action := range unsupportedActions {
		response.Schema.MarkdownDescription += fmt.Sprintf(
			"\n- %s.",
			action.string(unsupportedActions),
		)
	}
}
