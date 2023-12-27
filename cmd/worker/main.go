package main

import (
	"codebase/internal/worker/email"
	"github.com/hibiken/asynq"
)

func main() {
	processor := worker_email.NewEmailTaskProcessor(asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	})
	if err := processor.Process(); err != nil {
		panic(err)
	}
}
