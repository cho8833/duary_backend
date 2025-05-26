package scheduler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
)

func GetSchedulerClient() (*scheduler.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := scheduler.NewFromConfig(cfg)

	return client, nil
}
