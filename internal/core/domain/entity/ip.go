package entity

import (
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type Ip struct {
	Ip            string
	PrefixLength  string
	Version       int
	NullRouted    bool
	MainIp        bool
	NetworkType   enum.NetworkType
	ReverseLookup *string
	Ddos          *Ddos
}

type OptionalIpValues struct {
	Ddos          *Ddos
	ReverseLookup *string
}

func NewIp(
	ip string,
	prefixLength string,
	version int,
	nullRouted bool,
	mainIp bool,
	networkType enum.NetworkType,
	options OptionalIpValues,
) Ip {
	ipObject := Ip{
		Ip:           ip,
		PrefixLength: prefixLength,
		Version:      version,
		NullRouted:   nullRouted,
		MainIp:       mainIp,
		NetworkType:  networkType,
	}

	ipObject.Ddos = options.Ddos
	ipObject.ReverseLookup = options.ReverseLookup

	return ipObject
}
