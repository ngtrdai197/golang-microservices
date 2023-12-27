package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskEmailDelivery(
		ctx context.Context, payload *EmailDeliveryPayload, opts ...asynq.Option,
	) error
	DistributeTaskHelloDelivery(ctx context.Context, p *HelloDeliveryPayload, opts ...asynq.Option) error
	Close() error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(opts asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributor{
		client: asynq.NewClient(opts),
	}
}

func (distributor *RedisTaskDistributor) Close() error {
	return distributor.client.Close()
}
