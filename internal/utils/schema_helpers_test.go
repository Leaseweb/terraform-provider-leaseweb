package utils

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func TestAction_firstAction(t *testing.T) {
	t.Run(
		"return `imported` when unsupported actions contain `createAction`",
		func(t *testing.T) {
			got := CreateAction.firstAction([]Action{CreateAction})
			want := "imported"

			assert.Equal(t, want, got)
		},
	)

	t.Run(
		"return `created` when unsupported actions do not contain `createAction`",
		func(t *testing.T) {
			got := CreateAction.firstAction([]Action{CreateAction})
			want := "imported"

			assert.Equal(t, want, got)
		},
	)

}

func TestAction_string(t *testing.T) {
	t.Run(
		"expected string is returned for the default user case",
		func(t *testing.T) {
			tests := []struct {
				name string
				a    Action
				want string
			}{
				{
					name: "create action",
					a:    CreateAction,
					want: "Once created, this resource cannot be created",
				},
				{
					name: "read action",
					a:    ReadAction,
					want: "Once created, this resource cannot be read",
				},
				{
					name: "update action",
					a:    UpdateAction,
					want: "Once created, this resource cannot be updated",
				},
				{
					name: "delete action",
					a:    DeleteAction,
					want: "Once created, this resource cannot be deleted",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					got := tt.a.string(nil)
					assert.Equal(t, tt.want, got)
				})
			}
		},
	)

	t.Run(
		"expected string is returned when unsupported actions contain created and the current action is created",
		func(t *testing.T) {
			response := resource.SchemaResponse{}
			AddUnsupportedActionsNotation(&response, []Action{CreateAction})
			assert.Equal(
				t,
				"This resource cannot be created, only imported",
				CreateAction.string([]Action{CreateAction}),
			)
		},
	)

}

func TestAddUnsupportedActionsNotation(t *testing.T) {
	response := resource.SchemaResponse{}
	AddUnsupportedActionsNotation(&response, []Action{UpdateAction})

	assert.Equal(
		t,
		"**Note:**\n- Once created, this resource cannot be updated.",
		response.Schema.GetMarkdownDescription(),
	)
}

func ExampleAddUnsupportedActionsNotation() {
	response := resource.SchemaResponse{}
	AddUnsupportedActionsNotation(&response, []Action{UpdateAction})

	fmt.Println(response.Schema.GetMarkdownDescription())
	// Output:
	// **Note:**
	// - Once created, this resource cannot be updated.
}
