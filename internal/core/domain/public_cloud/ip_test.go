package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIp(t *testing.T) {
	ip := NewIp("127.0.0.1")

	assert.Equal(t, "127.0.0.1", ip.Ip)
}
