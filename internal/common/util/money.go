package util

import "github.com/Rhymond/go-money"

func ConvertToMoney(amount int, currency string) *money.Money {
	return money.New(int64(amount), currency)
}
