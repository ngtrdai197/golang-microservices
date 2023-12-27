package main

import (
	"codebase/internal/worker"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

func main() {
	processor := worker.NewRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	})
	log.Info().Msg("Starting asynq server...")
	if err := processor.Start(); err != nil {
		panic(err)
	}
}
