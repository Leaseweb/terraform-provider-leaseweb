package domain

type Prices struct {
	Currency       string
	CurrencySymbol string
	Compute        Price
	Storage        Storage
}

func NewPrices(
	currency string,
	currencySymbol string,
	compute Price,
	storage Storage,
) Prices {
	return Prices{
		Currency:       currency,
		CurrencySymbol: currencySymbol,
		Compute:        compute,
		Storage:        storage,
	}
}
