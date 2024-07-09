package enum

type NetworkType string

func (n NetworkType) String() string {
	return string(n)
}

type NetworkTypes []NetworkType

const (
	NetworkTypeInternal NetworkType = "INTERNAL"
	NetworkTypePublic   NetworkType = "PUBLIC"
)

var NetworkTypeValues = NetworkTypes{NetworkTypeInternal, NetworkTypePublic}
