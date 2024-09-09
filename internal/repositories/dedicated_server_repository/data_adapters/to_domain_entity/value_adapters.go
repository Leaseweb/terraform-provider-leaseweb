package to_domain_entity

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
)

func AdaptNullableStringToValue(nullableString dedicatedServer.NullableString) *string {
	return nullableString.Get()
}
