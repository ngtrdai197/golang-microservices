package transport

import (
	"codebase/api/post/v1"
	workeremail "codebase/internal/worker/email"
	"context"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PostTransport interface {
	v1.PostServer
}

type postTransport struct {
	v1.UnimplementedPostServer
	redisOpts asynq.RedisClientOpt
}

func NewPostTransport(redisOpts asynq.RedisClientOpt) PostTransport {
	return &postTransport{
		redisOpts: redisOpts,
	}
}

func (p *postTransport) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	task := workeremail.NewEmailDeliveryDistributor(p.redisOpts)
	err := task.DistributeTaskEmailDelivery(ctx, &workeremail.DeliveryPayload{
		Msg: "Create new post with name: " + req.Name + " and content: " + req.Content,
	})
	if err != nil {
		log.Err(err).Msg("Error create task new post delivery")
		return nil, err
	}
	return &v1.CreatePostResponse{
		Name:    req.Name,
		Content: req.Content,
	}, nil
}
