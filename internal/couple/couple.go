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

type CreateCoupleReq struct {
	RelationDate time.Time        `json:"relationDate"`
	Members      []*member.Member `json:"members"`
}

type UpdateCoupleReq struct {
	Id           *string          `json:"id" dynamodbav:"id,omitempty"`
	Members      []*member.Member `json:"members" dynamodbav:"members,omitempty"`
	RelationDate *time.Time       `json:"relationDate" dynamodbav:"relationDate,omitempty"`
	CoupleCode   *string          `json:"coupleCode" dynamodbav:"coupleCode,omitempty"`
}
