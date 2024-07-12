package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerificationEmail(
		ctx context.Context,
		payload *PayloadSendVerificationEmail,
		opts ...asynq.Option,
	) error
}

type RedisDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(opt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(opt)
	return &RedisDistributor{client: client}
}
