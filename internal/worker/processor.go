package worker

import (
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				QueuePriorityCritical: 6,
				QueuePriorityDefault:  3,
				QueuePriorityLow:      1,
			},
		},
	)
	return &RedisTaskProcessor{
		server: server,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(DeliveryEmailQueue, processor.HandleEmailDeliveryTask)
	mux.HandleFunc(DeliverySayHello, processor.HandleHelloDeliveryTask)

	return processor.server.Run(mux)
}
