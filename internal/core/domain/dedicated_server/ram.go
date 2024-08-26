package dedicated_server

type Ram struct {
	Size int32
	Unit string
}

func NewRam(size int32, unit string) Ram {
	return Ram{
		Size: size,
		Unit: unit,
	}
}
