package ws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func GetApiGwClient() (*apigatewaymanagementapi.Client, error) {
	//wsUrl := os.Getenv("WEBSOCKET_URL")
	wsUrl := "https://xgbo7rm9jk.execute-api.ap-northeast-2.amazonaws.com/production/"
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = &wsUrl
	})

	return client, nil
}
