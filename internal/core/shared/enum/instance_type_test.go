package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceType_String(t *testing.T) {
	got := InstanceTypeC3Large.String()

	assert.Equal(t, "lsw.c3.large", got)
}

func TestNewInstanceType(t *testing.T) {
	want := InstanceTypeC3Large
	got, err := NewInstanceType("lsw.c3.large")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestInstanceType_Values(t *testing.T) {
	want := []string{
		"lsw.m3.large",
		"lsw.m3.xlarge",
		"lsw.m3.2xlarge",
		"lsw.m4.large",
		"lsw.m4.xlarge",
		"lsw.m4.2xlarge",
		"lsw.m4.4xlarge",
		"lsw.m5.large",
		"lsw.m5.xlarge",
		"lsw.m5.2xlarge",
		"lsw.m5.4xlarge",
		"lsw.m5a.large",
		"lsw.m5a.xlarge",
		"lsw.m5a.2xlarge",
		"lsw.m5a.4xlarge",
		"lsw.m5a.8xlarge",
		"lsw.m5a.12xlarge",
		"lsw.m6a.large",
		"lsw.m6a.xlarge",
		"lsw.m6a.2xlarge",
		"lsw.m6a.4xlarge",
		"lsw.m6a.8xlarge",
		"lsw.m6a.12xlarge",
		"lsw.m6a.16xlarge",
		"lsw.m6a.24xlarge",
		"lsw.c3.large",
		"lsw.c3.xlarge",
		"lsw.c3.2xlarge",
		"lsw.c3.4xlarge",
		"lsw.c4.large",
		"lsw.c4.xlarge",
		"lsw.c4.2xlarge",
		"lsw.c4.4xlarge",
		"lsw.c5.large",
		"lsw.c5.xlarge",
		"lsw.c5.2xlarge",
		"lsw.c5.4xlarge",
		"lsw.c5a.large",
		"lsw.c5a.xlarge",
		"lsw.c5a.2xlarge",
		"lsw.c5a.4xlarge",
		"lsw.c5a.9xlarge",
		"lsw.c5a.12xlarge",
		"lsw.c6a.large",
		"lsw.c6a.xlarge",
		"lsw.c6a.2xlarge",
		"lsw.c6a.4xlarge",
		"lsw.c6a.8xlarge",
		"lsw.c6a.12xlarge",
		"lsw.c6a.16xlarge",
		"lsw.c6a.24xlarge",
		"lsw.r3.large",
		"lsw.r3.xlarge",
		"lsw.r3.2xlarge",
		"lsw.r4.large",
		"lsw.r4.xlarge",
		"lsw.r4.2xlarge",
		"lsw.r5.large",
		"lsw.r5.xlarge",
		"lsw.r5.2xlarge",
		"lsw.r5a.large",
		"lsw.r5a.xlarge",
		"lsw.r5a.2xlarge",
		"lsw.r5a.4xlarge",
		"lsw.r5a.8xlarge",
		"lsw.r5a.12xlarge",
		"lsw.r6a.large",
		"lsw.r6a.xlarge",
		"lsw.r6a.2xlarge",
		"lsw.r6a.4xlarge",
		"lsw.r6a.8xlarge",
		"lsw.r6a.12xlarge",
		"lsw.r6a.16xlarge",
		"lsw.r6a.24xlarge",
	}

	got := InstanceTypeC3Large.Values()

	assert.EqualValues(t, want, got)
}
