package ws_connection

type WSConnection struct {
	SocialId     string `dynamodbav:"socialId"`
	Provider     string `dynamodbav:"provider"`
	ConnectionId string `dynamodbav:"connectionId"`
}
