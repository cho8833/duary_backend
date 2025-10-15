package ws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"os"
)

func GetApiGwClient(stage string) (*apigatewaymanagementapi.Client, error) {
	var wsUrl string
	if stage == "dev" {
		wsUrl = os.Getenv("DEV_WEBSOCKET_URL")
	} else {
		wsUrl = os.Getenv("PROD_WEBSOCKET_URL")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = &wsUrl
	})

	return client, nil
}
