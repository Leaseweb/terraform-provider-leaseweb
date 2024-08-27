package dedicated_server

type Contract struct {
	Id             string
	CustomerId     string
	DeliveryStatus string
	Reference      string
	SalesOrgId     string
}

func NewContract(id, customerId, deliveryStatus, reference, salesOrgId string) Contract {
	return Contract{
		Id:             id,
		CustomerId:     customerId,
		DeliveryStatus: deliveryStatus,
		Reference:      reference,
		SalesOrgId:     salesOrgId,
	}
}
