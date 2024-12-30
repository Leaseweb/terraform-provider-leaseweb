package ipmgmt

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/stretchr/testify/assert"
)

func Test_adaptIPToIPResourceModel(t *testing.T) {
	diags := diag.Diagnostics{}

	reverseLookup := "reverseLookup"
	nullLevel := int32(3)
	sdkIP := ipmgmt.Ip{
		Ip:               "1.2.3.4",
		Version:          4,
		Type:             "NORMAL_IP",
		PrefixLength:     0,
		Primary:          false,
		ReverseLookup:    *ipmgmt.NewNullableString(&reverseLookup),
		NullRouted:       true,
		NullLevel:        *ipmgmt.NewNullableInt32(&nullLevel),
		UnnullingAllowed: false,
		EquipmentId:      "equipmentId",
		AssignedContract: *ipmgmt.NewNullableAssignedContract(
			&ipmgmt.AssignedContract{
				Id: "contractId",
			},
		),
		Subnet: ipmgmt.Subnet{
			Id:           "subnetId",
			NetworkIp:    "5.6.7.8",
			PrefixLength: 24,
			Gateway:      "2.2.2.2",
		},
	}
	got := adaptIPToIPResourceModel(sdkIP, context.TODO(), &diags)

	assert.False(t, diags.HasError())

	assert.Equal(t, "1.2.3.4", got.IP.ValueString())
	assert.Equal(t, int32(4), got.Version.ValueInt32())
	assert.Equal(t, "NORMAL_IP", got.Type.ValueString())
	assert.False(t, got.Primary.ValueBool())
	assert.Equal(t, reverseLookup, got.ReverseLookup.ValueString())
	assert.True(t, got.NullRouted.ValueBool())
	assert.Equal(t, nullLevel, got.NullLevel.ValueInt32())
	assert.False(t, got.UnnullingAllowed.ValueBool())
	assert.Equal(t, "equipmentId", got.EquipmentID.ValueString())

	assignedContract := assignedContractResourceModel{}
	got.AssignedContract.As(
		context.TODO(),
		&assignedContract,
		basetypes.ObjectAsOptions{},
	)
	assert.Equal(t, "contractId", assignedContract.ID.ValueString())

	subnet := subnetResourceModel{}
	got.Subnet.As(context.TODO(), &subnet, basetypes.ObjectAsOptions{})
	assert.Equal(t, "subnetId", subnet.ID.ValueString())
	assert.Equal(t, "5.6.7.8", subnet.NetworkIP.ValueString())
	assert.Equal(t, int32(24), subnet.PrefixLength.ValueInt32())
	assert.Equal(t, "2.2.2.2", subnet.Gateway.ValueString())

}
