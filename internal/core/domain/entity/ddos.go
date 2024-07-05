package entity

type Ddos struct {
	DetectionProfile string
	ProtectionType   string
}

func NewDdos(detectionProfile string, protectionType string) Ddos {
	return Ddos{
		DetectionProfile: detectionProfile,
		ProtectionType:   protectionType,
	}
}
