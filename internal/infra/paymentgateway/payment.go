package paymentgateway

import "context"

type Payment struct {
	ID       string
	Amount   float64
	Currency string
	Status   string
	IntentID string
}

type Refund struct {
	IntentID string
	Reason   string
}

type PaymentGateway interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	CreateRefund(ctx context.Context, refund *Refund) error
}
