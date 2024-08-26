package dedicated_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewContract(t *testing.T) {

	got := NewContract("id", "customerId", "deliveryStatus", "reference", "salesOrgId")
	want := Contract{
		Id:             "id",
		CustomerId:     "customerId",
		DeliveryStatus: "deliveryStatus",
		Reference:      "reference",
		SalesOrgId:     "salesOrgId",
	}
	assert.Equal(t, want, got)
}
