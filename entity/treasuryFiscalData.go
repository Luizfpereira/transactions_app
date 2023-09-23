package entity

import (
	"github.com/shopspring/decimal"
)

type TreasuryFiscalData struct {
	CoutryCurrencyDesc string          `json:"country_currency_desc,omitempty"`
	ExchageRate        decimal.Decimal `json:"exchange_rate,omitempty"`
	RecordDate         string          `json:"record_date,omitempty"`
}

type TreasuryApi struct {
	Data []TreasuryFiscalData `json:"data,omitempty"`
}
