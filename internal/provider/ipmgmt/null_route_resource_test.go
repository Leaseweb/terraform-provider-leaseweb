package ipmgmt

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/stretchr/testify/assert"
)

func Test_adaptNullRouteToResourceModel(t *testing.T) {
	nulledAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:47",
	)
	automatedUnnullingAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:48",
	)
	unnulledAt, _ := time.Parse(
		"2006-01-02 15:04:05",
		"2023-12-14 17:09:49",
	)
	unnulledBy := "unnulled@example.com"
	ticketID := "ticketId"
	comment := "comment"

	sdkAssignedContract := ipmgmt.AssignedContract{
		Id: "assignedContractId",
	}

	nullroutedIP := ipmgmt.NullRoutedIP{
		Id:                   "id",
		Ip:                   "1.2.3.4",
		NulledAt:             nulledAt,
		NulledBy:             "john@example.com",
		NullLevel:            3,
		AutomatedUnnullingAt: *ipmgmt.NewNullableTime(&automatedUnnullingAt),
		UnnulledAt:           *ipmgmt.NewNullableTime(&unnulledAt),
		UnnulledBy:           *ipmgmt.NewNullableString(&unnulledBy),
		TicketId:             *ipmgmt.NewNullableString(&ticketID),
		Comment:              *ipmgmt.NewNullableString(&comment),
		EquipmentId:          "equipmentId",
		AssignedContract:     *ipmgmt.NewNullableAssignedContract(&sdkAssignedContract),
	}

	diags := diag.Diagnostics{}
	got := adaptNullRouteToResourceModel(nullroutedIP, &diags, context.TODO())

	assert.False(t, diags.HasError())

	assignedContract := assignedContractResourceModel{}
	got.AssignedContract.As(context.TODO(), &assignedContract, basetypes.ObjectAsOptions{})
	assert.Equal(t, "assignedContractId", assignedContract.ID.ValueString())

	assert.Equal(t, "id", got.ID.ValueString())
	assert.Equal(t, "1.2.3.4", got.IP.ValueString())
	assert.Equal(
		t,
		"2023-12-14 17:09:49 +0000 UTC",
		got.UnnulledAt.ValueString(),
	)
	assert.Equal(t, "unnulled@example.com", got.UnnulledBy.ValueString())
	assert.Equal(t, int32(3), got.NullLevel.ValueInt32())
	assert.Equal(
		t,
		"2023-12-14 17:09:48 +0000 UTC",
		got.AutomaticUnnullingAt.ValueString(),
	)
	assert.Equal(
		t,
		"2023-12-14 17:09:49 +0000 UTC",
		got.UnnulledAt.ValueString(),
	)
	assert.Equal(t, "unnulled@example.com", got.UnnulledBy.ValueString())
	assert.Equal(t, "ticketId", got.TicketID.ValueString())
	assert.Equal(t, "comment", got.Comment.ValueString())
	assert.Equal(t, "equipmentId", got.EquipmentID.ValueString())
}
