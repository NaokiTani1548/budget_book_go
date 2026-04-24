package valueobject

import "errors"

type Money struct {
	amount float64
}

func NewMoney(amount float64) (Money, error) {
	if amount <= 0 {
		return Money{}, errors.New("金額は0以上です")
	}
	return Money{amount: amount}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) isGreaterThan(other Money) bool {
	return m.amount > other.amount
}