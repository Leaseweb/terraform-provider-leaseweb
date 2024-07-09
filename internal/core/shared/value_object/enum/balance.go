package enum

type Balance string

func (b Balance) String() string {
	return string(b)
}

type Balances []Balance

const (
	BalanceRoundRobin Balance = "ROUNDROBIN"
	BalanceLeastConn  Balance = "LEASTCONN"
	BalanceSource     Balance = "SOURCE"
)

var BalanceValues = Balances{BalanceRoundRobin, BalanceLeastConn, BalanceSource}
