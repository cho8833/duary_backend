package shared

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type ServerResponse[T any] struct {
	Message *string `json:"message"`
	Status  int     `json:"status"`
	Data    T       `json:"data"`
}

func ErrorResponse(message string, statusCode int) *ServerResponse[any] {
	return &ServerResponse[any]{
		&message,
		statusCode,
		nil,
	}
}

func ResponseFromError(error error, statusCode int) *ServerResponse[any] {
	errorMessage := error.Error()
	return &ServerResponse[any]{
		&errorMessage,
		statusCode,
		nil,
	}
}

func ResponseWithData(data any) *ServerResponse[any] {
	okString := "OK"
	return &ServerResponse[any]{
		&okString,
		200,
		&data,
	}
}

func SUCCESS() *ServerResponse[any] {
	okString := "OK"

	return &ServerResponse[any]{
		&okString,
		200,
		nil,
	}
}

func LambdaResponseWithData(data any) events.APIGatewayProxyResponse {
	res := ResponseWithData(data)
	b, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
	}
}

func LambdaResponseWithDataAndHeader(data any, header map[string]string) events.APIGatewayProxyResponse {
	res := ResponseWithData(data)
	b, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		Headers:    header,
		StatusCode: 200,
		Body:       string(b),
	}
}

func LambdaAppErrorResponse(err ApplicationError) events.APIGatewayProxyResponse {
	return LambdaErrorResponse(err, err.StatusCode())
}
func LambdaErrorResponse(err error, statusCode int) events.APIGatewayProxyResponse {
	res := ResponseFromError(err, statusCode)
	b, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
	}
}
