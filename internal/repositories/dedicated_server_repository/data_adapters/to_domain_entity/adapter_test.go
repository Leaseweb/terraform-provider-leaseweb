package to_domain_entity

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/stretchr/testify/assert"
)

func Test_adaptDedicatedServerId(t *testing.T) {
	t.Run("ID is set", func(t *testing.T) {
		// TODO: we need to check if ID is required attribute for dedicated server!
		dedicatedServer := dedicatedServer.NewServer()
		dedicatedServer.SetId("123456")
		got, err := AdaptDedicatedServer(*dedicatedServer)

		assert.NoError(t, err)
		assert.Equal(t, "123456", got.Id)
	})
}
