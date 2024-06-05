package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetBoolValue(t *testing.T) {
	type args struct {
		hasValue bool
		value    bool
	}

	tests := []struct {
		name string
		args args
		want basetypes.BoolValue
	}{
		{name: "Should return true", args: args{hasValue: true, value: true}, want: basetypes.NewBoolValue(true)},
		{name: "Should return false", args: args{hasValue: true, value: false}, want: basetypes.NewBoolValue(false)},
		{name: "Should return null", args: args{hasValue: false, value: true}, want: basetypes.NewBoolNull()},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetBoolValue(tt.args.hasValue, tt.args.value), tt.name)
	}
}

func TestGetDateTime(t *testing.T) {
	zeroTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "0001-01-01 00:00:00 +0000 UTC ")
	timestamp, _ := time.Parse("2006-01-02 15:04:05", "2023-12-14 17:09:47")

	type args struct {
		value time.Time
	}

	tests := []struct {
		name string
		args args
		want basetypes.StringValue
	}{
		{name: "Should return a time", args: args{value: timestamp}, want: basetypes.NewStringValue("2023-12-14 17:09:47 +0000 UTC")},
		{name: "Should return null", args: args{value: zeroTime}, want: basetypes.NewStringNull()},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetDateTime(tt.args.value), tt.name)
	}
}

func TestGetFloatValue(t *testing.T) {
	type args struct {
		hasValue bool
		value    float32
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "Should return a float", args: args{hasValue: true, value: 1.2300000190734863}, want: 1.2300000190734863},
		{name: "Should return empty", args: args{hasValue: false, value: 1.22}, want: 0},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetFloatValue(tt.args.hasValue, tt.args.value).ValueFloat64(), tt.name)
	}
}

func TestGetIntValue(t *testing.T) {
	type args struct {
		hasValue bool
		value    int32
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "Should return an int", args: args{hasValue: true, value: 1}, want: 1},
		{name: "Should return null", args: args{hasValue: false, value: 2}, want: 0},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, GetIntValue(tt.args.hasValue, tt.args.value).ValueInt64(), tt.name)
	}
}

func TestGetStringValue(t *testing.T) {
	type args struct {
		hasValue bool
		value    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Should return a string", args: args{hasValue: true, value: "tralala"}, want: "tralala"},
		{name: "Should return null", args: args{hasValue: false, value: "tralala"}, want: ""},
	}
	for _, tt := range tests {

		assert.Equal(t, tt.want, GetStringValue(tt.args.hasValue, tt.args.value).ValueString(), tt.name)
	}
}
