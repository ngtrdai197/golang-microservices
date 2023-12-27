package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type HelloDeliveryPayload struct {
	Msg string `json:"msg"`
}

func (distributor *RedisTaskDistributor) DistributeTaskHelloDelivery(ctx context.Context, p *HelloDeliveryPayload, opts ...asynq.Option) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to markshal task payload: %w", err)
	}
	task := asynq.NewTask(DeliverySayHello, payload, opts...)
	if err != nil {
		log.Fatal().Msgf("Error create task msg delivery detail = %v", err)
		return err
	}

	// Process the task immediately.
	info, err := distributor.client.EnqueueContext(ctx, task)
	log.Info().Msgf("info = %+v", info)
	if err != nil {
		return fmt.Errorf("Error add task msg delivery into queue detail = %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Msg("Successfully enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) HandleHelloDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p HelloDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Info().Msgf("Sending msg=%s", p.Msg)
	return nil
}
