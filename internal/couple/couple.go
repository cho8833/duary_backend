package couple

import (
	"github.com/cho8833/duary_lambda/internal/member"
	"time"
)

type Couple struct {
	Id           *string          `json:"id" dynamodbav:"id"`
	RelationDate *time.Time       `json:"relationDate" dynamodbav:"relationDate"`
	Members      []*member.Member `json:"members" dynamodbav:"members"`
	Code         *string          `json:"code" dynamodbav:"code"`
}
