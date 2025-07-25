package test

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/model/event"
	"log"
	"testing"
	"time"
)

func Test_GetBetweenStartAndEndDate(t *testing.T) {
	client, _ := model.GetDynamoDBClient()

	repo := event.NewRepository(client)
	svc := event.NewService(repo)

	now := timePtr(time.Now())
	rangeStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	rangeEnd := rangeStart.Add(24 * time.Hour)
	events, err := svc.GetBetweenStartAndEndDate("dafe98ab-9eac-4b66-ba47-a389783d1a19", rangeStart, rangeEnd)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(events)

}
