package resource

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncedMap_Get(t *testing.T) {
	t.Run("existing value is returned", func(t *testing.T) {
		syncedMap := NewSyncedMap[string, string]()

		syncedMap.Set("foo", "bar")
		got, ok := syncedMap.Get("foo")

		assert.True(t, ok)
		assert.Equal(t, "bar", got)
	})

	t.Run("non existent value returns false", func(t *testing.T) {
		syncedMap := NewSyncedMap[string, string]()

		_, ok := syncedMap.Get("foo")

		assert.False(t, ok)
	})
}

func ExampleNewSyncedMap() {
	syncedMap := NewSyncedMap[string, string]()
	syncedMap.Set("foo", "bar")

	fmt.Println(syncedMap.Get("foo"))
	// Output: bar true
}
