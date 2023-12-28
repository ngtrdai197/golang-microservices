package main

import (
	"codebase/config"
	"codebase/internal/worker/email"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msgf("Recovered from panic: %v", err)
		}
	}()
	processor := workeremail.NewEmailTaskProcessor(asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", config.Cfg.RedisHost, config.Cfg.RedisPort),
		DB:   config.Cfg.RedisDatabase,
	})
	if err := processor.Process(); err != nil {
		panic(err)
	}
}

func init() {
	config.LoadConfig()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		if level == zerolog.ErrorLevel {
			e.Str("stack", fmt.Sprintf("%+v", errors.WithStack(errors.New(msg))))
		}
	}))
}
