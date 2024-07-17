package modify_plan

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain"
)

func TestTypeValidator_HashTypeChanged(t *testing.T) {
	type fields struct {
		stateInstanceId   types.String
		stateInstanceType types.String
		planInstanceType  types.String
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"always returns false on resource creation",
			fields{
				stateInstanceId: basetypes.NewStringUnknown(),
			},
			false,
		},
		{
			"always returns false on resource import",
			fields{
				stateInstanceId:  basetypes.NewStringValue("123"),
				planInstanceType: basetypes.NewStringUnknown(),
			},
			false,
		},
		{
			"always returns false when type doesn't change",
			fields{
				stateInstanceId:   basetypes.NewStringValue("123"),
				stateInstanceType: basetypes.NewStringValue("lsw.m3.large"),
				planInstanceType:  basetypes.NewStringValue("lsw.m3.large"),
			},
			false,
		},
		{
			"returns true when type changes",
			fields{
				stateInstanceId:   basetypes.NewStringValue("123"),
				stateInstanceType: basetypes.NewStringValue("lsw.m3.large"),
				planInstanceType:  basetypes.NewStringValue("lsw.m4.large"),
			},
			true,
		},
	}
	for _, tt := range tests {
		v := NewTypeValidator(tt.fields.stateInstanceId, tt.fields.stateInstanceType, tt.fields.planInstanceType)
		t.Run(tt.name, func(t *testing.T) {
			got := v.HashTypeChanged()
			assert.Equal(t, tt.want, got, fmt.Sprintf("%v", tt.name))
		})
	}
}

func TestTypeValidator_IsTypeValid(t *testing.T) {
	type fields struct {
		stateInstanceId   types.String
		stateInstanceType types.String
		planInstanceType  types.String
	}
	type args struct {
		allowedInstanceTypes domain.InstanceTypes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "is valid",
			fields: fields{
				stateInstanceId:   basetypes.NewStringValue("123"),
				stateInstanceType: basetypes.NewStringValue("oldValue"),
				planInstanceType:  basetypes.NewStringValue("newValue"),
			},
			args: args{
				allowedInstanceTypes: domain.InstanceTypes{{Name: "newValue"}},
			},
			want: true,
		},
		{
			name: "is not valid",
			fields: fields{
				stateInstanceId:   basetypes.NewStringValue("123"),
				stateInstanceType: basetypes.NewStringValue("oldValue"),
				planInstanceType:  basetypes.NewStringValue("newValue"),
			},
			args: args{
				allowedInstanceTypes: domain.InstanceTypes{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := TypeValidator{
				stateInstanceId:   tt.fields.stateInstanceId,
				stateInstanceType: tt.fields.stateInstanceType,
				planInstanceType:  tt.fields.planInstanceType,
			}
			assert.Equalf(
				t,
				tt.want,
				v.IsTypeValid(tt.args.allowedInstanceTypes),
				"IsTypeValid(%v)",
				tt.args.allowedInstanceTypes,
			)
		})
	}
}

func TestTypeValidator_IsBeingCreated(t *testing.T) {
	type fields struct {
		stateInstanceId   types.String
		stateInstanceType types.String
		planInstanceType  types.String
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "is being created",
			fields: fields{},
			want:   true,
		},
		{
			name:   "is not created",
			fields: fields{stateInstanceId: basetypes.NewStringValue("123")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := TypeValidator{
				stateInstanceId:   tt.fields.stateInstanceId,
				stateInstanceType: tt.fields.stateInstanceType,
				planInstanceType:  tt.fields.planInstanceType,
			}
			assert.Equalf(t, tt.want, v.IsBeingCreated(), "IsBeingCreated()")
		})
	}
}
