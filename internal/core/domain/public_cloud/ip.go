package public_cloud

type Ip struct {
	Ip string
}

func NewIp(ip string) Ip {
	return Ip{
		Ip: ip,
	}
}
