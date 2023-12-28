package workeremail

import (
	"codebase/constant"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type DeliveryPayload struct {
	Msg string `json:"msg"`
}

func (distributor *taskDistributor) DistributeTaskEmailDelivery(ctx context.Context, p *DeliveryPayload, opts ...asynq.Option) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to markshal task payload: %w", err)
	}
	task := asynq.NewTask(constant.DeliveryEmail, payload, opts...)
	if err != nil {
		log.Fatal().Msgf("Error create task msg delivery detail = %v", err)
		return err
	}

	// Process the task immediately.
	info, err := distributor.client.EnqueueContext(ctx, task)
	log.Info().Msgf("info = %+v", info)
	if err != nil {
		return fmt.Errorf("error add task msg delivery into queue detail = %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("task_id", info.ID).Str("queue", info.Queue).Msg("Successfully enqueued task")
	return nil
}

func (p *processor) HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var payload DeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Info().Str("type", t.Type()).Str("task_id", t.ResultWriter().TaskID()).Msgf("Sending msg=%s", payload.Msg)
	return nil
}
