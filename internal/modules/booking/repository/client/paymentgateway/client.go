package paymentgateway

import (
	"booking-event/internal/infra/paymentgateway"
	"context"
)

type Client struct {
	paymentGateway paymentgateway.PaymentGateway
}

func NewClient(paymentGateway paymentgateway.PaymentGateway) *Client {
	return &Client{paymentGateway: paymentGateway}
}

func (c *Client) CreatePayment(ctx context.Context, payment *paymentgateway.Payment) error {
	return c.paymentGateway.CreatePayment(ctx, payment)
}
