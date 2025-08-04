package test

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/ws_connection"
	"github.com/cho8833/duary_lambda/ws"
	"testing"
)

func Test_send(t *testing.T) {
	dbClient, err := model.GetDynamoDBClient()

	if err != nil {
		t.Fatalf(err.Error())
	}

	wsRepo := ws_connection.NewRepository(dbClient)
	wsSvc, err := ws.NewService(&wsRepo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = wsSvc.Send("11093972176936539879", "google", ws.CoupleConnected, "test")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
