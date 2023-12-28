package workeremail

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskEmailDelivery(ctx context.Context, p *DeliveryPayload, opts ...asynq.Option) error
	Close() error
}

type taskDistributor struct {
	client *asynq.Client
}

func NewEmailDeliveryDistributor(opts asynq.RedisClientOpt) TaskDistributor {
	return &taskDistributor{
		client: asynq.NewClient(opts)}
}

func (distributor *taskDistributor) Close() error {
	return distributor.client.Close()
}
