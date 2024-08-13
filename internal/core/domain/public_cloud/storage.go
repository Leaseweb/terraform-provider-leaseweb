package public_cloud

type Storage struct {
	Local   Price
	Central Price
}

func NewStorage(local Price, central Price) Storage {
	return Storage{Local: local, Central: central}
}
