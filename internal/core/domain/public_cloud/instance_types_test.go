package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceTypes_ToArray(t *testing.T) {
	instanceTypes := InstanceTypes{InstanceType{Name: "tralala"}}
	got := instanceTypes.ToArray()
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}

func TestInstanceTypes_ContainsName(t *testing.T) {
	t.Run("return true if name exists", func(t *testing.T) {
		instanceTypes := InstanceTypes{InstanceType{Name: "tralala"}}
		assert.True(t, instanceTypes.ContainsName("tralala"))
	})

	t.Run("return false if name does not exist", func(t *testing.T) {
		instanceTypes := InstanceTypes{InstanceType{Name: "piet"}}
		assert.False(t, instanceTypes.ContainsName("tralala"))
	})
}

func TestInstanceTypes_GetByName(t *testing.T) {
	t.Run(
		"return instanceType if instanceType with name is found",
		func(t *testing.T) {
			want := InstanceType{Name: "tralala"}
			instanceTypes := InstanceTypes{want}
			got, err := instanceTypes.GetByName("tralala")

			assert.NoError(t, err)
			assert.Equal(t, want, *got)
		},
	)

	t.Run(
		"return error if instanceType with name is not found",
		func(t *testing.T) {
			instanceTypes := InstanceTypes{InstanceType{Name: "piet"}}
			_, err := instanceTypes.GetByName("tralala")

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)
}
