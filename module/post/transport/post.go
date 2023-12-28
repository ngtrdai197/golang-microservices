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
	distributor := workeremail.NewEmailDeliveryDistributor(p.redisOpts)
	err := distributor.DistributeTaskEmailDelivery(ctx, &workeremail.DeliveryPayload{
		Msg: "Create new post with name: " + req.Name + " and content: " + req.Content,
	})
	/**
	 * Note: Close redis client after finish,
	 * reduce connection to redis server and avoid memory leak in long-running process (like server)
	 * when create new redis client for each request
	 */
	defer func(distributor workeremail.TaskDistributor) {
		err := distributor.Close()
		if err != nil {
			log.Fatal().Msgf("Error close redis client detail = %v", err)
		}
	}(distributor)
	if err != nil {
		log.Err(err).Msg("Error create distributor new post delivery")
		return nil, err
	}
	return &v1.CreatePostResponse{
		Name:    req.Name,
		Content: req.Content,
	}, nil
}
