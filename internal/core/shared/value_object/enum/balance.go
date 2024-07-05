package enum

type Balance string

const (
	BalanceRoundRobin Balance = "ROUNDROBIN"
	BalanceLeastConn  Balance = "LEASTCONN"
	BalanceSource     Balance = "SOURCE"
)
