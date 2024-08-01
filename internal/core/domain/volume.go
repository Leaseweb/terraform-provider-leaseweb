package domain

type Volume struct {
	Size float64
	Unit string
}

func NewVolume(size float64, unit string) Volume {
	return Volume{
		Size: size,
		Unit: unit,
	}
}
