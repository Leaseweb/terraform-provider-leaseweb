package dedicated_server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPciCard(t *testing.T) {
	pciCard := NewPciCard("description")
	assert.Equal(t, "description", pciCard.Description)
	assert.Equal(t, "description", fmt.Sprintf("%v", pciCard))
}
