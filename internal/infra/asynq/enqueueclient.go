package asynq

import (
	"context"

	"github.com/hibiken/asynq"
)

type AsyncTaskEnqueueClient interface {
	Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) error
}

type enqueueClient struct {
	client *asynq.Client
}

func NewEnqueueClient(config Config) AsyncTaskEnqueueClient {
	return &enqueueClient{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr}),
	}
}

func (c *enqueueClient) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) error {
	_, err := c.client.EnqueueContext(ctx, task, opts...)
	return err
}
