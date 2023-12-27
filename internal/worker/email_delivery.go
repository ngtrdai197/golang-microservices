package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type EmailDeliveryPayload struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (distributor *RedisTaskDistributor) DistributeTaskEmailDelivery(ctx context.Context, p *EmailDeliveryPayload, opts ...asynq.Option) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to markshal task payload: %w", err)
	}
	task := asynq.NewTask(DeliveryEmailQueue, payload, opts...)
	if err != nil {
		log.Fatal().Msgf("Error create task email delivery detail = %v", err)
	}

	// Process the task immediately.
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("Error add task email delivery into queue detail = %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Msg("Successfully enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Info().Msgf("Sending Email to user: name=%s, phone=%s", p.Name, p.Phone)
	return nil
}
