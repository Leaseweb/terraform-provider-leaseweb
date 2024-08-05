package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type Balance string

func (b Balance) String() string {
	return string(b)
}

const (
	BalanceRoundRobin Balance = "roundrobin"
	BalanceLeastConn  Balance = "leastconn"
	BalanceSource     Balance = "source"
)

var balances = []Balance{BalanceRoundRobin, BalanceLeastConn, BalanceSource}

func NewBalance(value string) (Balance, error) {
	return enum_utils.FindEnumForString(value, balances, BalanceRoundRobin)
}
