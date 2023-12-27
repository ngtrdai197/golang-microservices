package worker_email

import (
	"codebase/internal/worker"
	"context"
	"github.com/hibiken/asynq"
)

const DeliveryEmail = "task:email"

type Processor interface {
	HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error
	Process() error
}

type processor struct {
	s *asynq.Server
}

func NewEmailTaskProcessor(redisOpt asynq.RedisClientOpt) Processor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				worker.QueuePriorityCritical: 6,
				worker.QueuePriorityDefault:  3,
				worker.QueuePriorityLow:      1,
			},
		},
	)
	return &processor{
		s: server,
	}
}

func (p *processor) Process() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(DeliveryEmail, p.HandleEmailDeliveryTask)

	return p.s.Run(mux)
}
