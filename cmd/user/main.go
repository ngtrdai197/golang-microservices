package main

import (
	"codebase/internal/worker"
	"context"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {
	distribute := worker.NewRedisTaskDistributor(asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	})

	defer func(distribute worker.TaskDistributor) {
		err := distribute.Close()
		if err != nil {
			log.Fatal().Msgf("Error close redis client detail = %v", err)
		}
	}(distribute)

	err := distribute.DistributeTaskHelloDelivery(context.Background(), &worker.HelloDeliveryPayload{
		Msg: "Hello World",
	})

	if err != nil {
		log.Fatal().Msgf("Error distribute task hello delivery detail = %v", err)
	}
}
