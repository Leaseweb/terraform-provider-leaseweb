package public_cloud

type StorageSize struct {
	Size float64
	Unit string
}

func NewStorageSize(size float64, unit string) StorageSize {
	return StorageSize{
		Size: size,
		Unit: unit,
	}
}
