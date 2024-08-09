package public_cloud

type Price struct {
	HourlyPrice  string
	MonthlyPrice string
}

func NewPrice(hourlyPrice string, monthlyPrice string) Price {
	return Price{
		HourlyPrice:  hourlyPrice,
		MonthlyPrice: monthlyPrice,
	}
}
