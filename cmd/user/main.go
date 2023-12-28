package main

import (
	postV1 "codebase/api/post/v1"
	userV1 "codebase/api/user/v1"
	"codebase/config"
	postTrpt "codebase/module/post/transport"
	userTrpt "codebase/module/user/transport"
	errorhandling "codebase/pkg/error"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msgf("Recovered from panic: %v", err)
		}
	}()

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Cfg.GrpcPort))
	if err != nil {
		log.Fatal().Msgf("Error listen port detail = %v", err)
	}
	grpcServer := grpc.NewServer()

	// FIXME: Register transport
	userTransport := userTrpt.NewUserTransport()
	userV1.RegisterUserServer(grpcServer, userTransport)

	postTransport := postTrpt.NewPostTransport(asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", config.Cfg.RedisHost, config.Cfg.RedisPort),
		DB:   config.Cfg.RedisDatabase,
	})
	postV1.RegisterPostServer(grpcServer, postTransport)

	go func() {
		log.Fatal().Msgf("Serving gRPC on %s", grpcServer.Serve(listen).Error())
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("%s:%d", "0.0.0.0", config.Cfg.GrpcPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Msgf("Error dial grpc server detail = %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal().Msgf("Error close grpc client detail = %v", err)
		}
	}(conn)

	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
			w.Header().Set("Vary", "Origin")
			return nil
		}),
		runtime.WithErrorHandler(func(ctx context.Context, sm *runtime.ServeMux, marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			// return Internal when Marshal failed
			const fallback = `{"code": 13, "message": "failed to marshal error message"}`

			s := status.Convert(err)
			pb := s.Proto()

			w.Header().Del("Trailer")
			w.Header().Del("Transfer-Encoding")

			contentType := marshaller.ContentType(pb)
			w.Header().Set("Content-Type", contentType)

			if s.Code() == codes.Unauthenticated {
				w.Header().Set("WWW-Authenticate", s.Message())
			}

			buf, marshalErr := marshaller.Marshal(pb)
			if marshalErr != nil {
				grpclog.Infof("Failed to marshal error message %q: %v", s, marshalErr)
				w.WriteHeader(http.StatusInternalServerError)

				if _, err := io.WriteString(w, fallback); err != nil {
					grpclog.Infof("Failed to write response: %v", err)
				}
				return
			}
			if s.Code() >= 4000 && s.Code() < 5000 {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(errorhandling.HTTPStatusFromCode(s.Code()))
			}
			if _, err := w.Write(buf); err != nil {
				grpclog.Infof("Failed to write response: %v", err)
			}
		}),
	)

	// FIXME: Register handler grpc gateway
	err = userV1.RegisterUserHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal().Msgf("Error register handler grpc gateway detail = %v", err)
	}
	err = postV1.RegisterPostHandler(context.Background(), gatewayMux, conn)
	if err != nil {
		log.Fatal().Msgf("Error register handler grpc gateway detail = %v", err)
	}

	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Accept-Language"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(gatewayMux)

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Cfg.HttpPort),
		Handler: withCors,
	}

	log.Info().Msgf("Serving gRPC-Gateway on %s", fmt.Sprintf("0.0.0.0:%d", config.Cfg.HttpPort))
	err = gwServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msgf("Error listen and serve http server detail = %v", err)
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
