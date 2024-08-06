package model

type Storage struct {
	Local   Price `tfsdk:"local"`
	Central Price `tfsdk:"central"`
}
