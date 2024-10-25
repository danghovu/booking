package paymentgateway

import (
	"context"
	"fmt"
)

type noopPaymentGateway struct{}

func NewNoopPaymentGateway() PaymentGateway {
	return &noopPaymentGateway{}
}

func (pg *noopPaymentGateway) CreatePayment(ctx context.Context, payment *Payment) error {
	fmt.Println("create payment", payment)
	return nil
}

func (pg *noopPaymentGateway) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	return nil, nil
}

func (pg *noopPaymentGateway) CreateRefund(ctx context.Context, refund *Refund) error {
	return nil
}
