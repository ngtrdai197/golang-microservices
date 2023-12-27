package main

import (
	"codebase/internal/worker/email"
	"context"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	distribute := worker_email.NewEmailDeliveryDistributor(asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	})

	defer func(distribute worker_email.TaskDistributor) {
		err := distribute.Close()
		if err != nil {
			log.Fatal().Msgf("Error close redis client detail = %v", err)
		}
	}(distribute)

	err := distribute.DistributeTaskEmailDelivery(context.Background(), &worker_email.DeliveryPayload{
		Msg: "dainguyen.iammm+1@gmail.com",
	}, asynq.ProcessIn(5*time.Second)) // delay process task duration 5s

	if err != nil {
		log.Fatal().Msgf("Error distribute task email delivery detail = %v", err)
	}
}
