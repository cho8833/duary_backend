package couple

import (
	"github.com/cho8833/duary_backend/model/member"
	"time"
)

type Couple struct {
	Id                 *string         `json:"id" dynamodbav:"id"`
	RelationDate       *time.Time      `json:"relationDate" dynamodbav:"relationDate"`
	Members            []member.Member `json:"members" dynamodbav:"members"`
	Code               *string         `json:"code" dynamodbav:"code"`
	ConnectedMemberIds []string        `json:"-" dynamodbav:"connectedMemberIds"`
}

func (c *Couple) ApplyFrom(req UpdateCoupleReq) {
	if req.RelationDate != nil {
		c.RelationDate = req.RelationDate
	}
	if req.Members != nil {
		c.Members = req.Members
	}
	if req.Code != nil {
		c.Code = req.Code
	}
	if req.ConnectedMemberIds != nil {
		c.ConnectedMemberIds = req.ConnectedMemberIds
	}
}

type CreateCoupleReq struct {
	RelationDate       *time.Time      `json:"relationDate"`
	Members            []member.Member `json:"members"`
	ConnectedMemberIds []string        `json:"-"`
}

type UpdateCoupleReq struct {
	Id                 *string         `json:"id" dynamodbav:"id,omitempty"`
	Members            []member.Member `json:"members" dynamodbav:"members,omitempty"`
	RelationDate       *time.Time      `json:"relationDate" dynamodbav:"relationDate,omitempty"`
	Code               *string         `json:"code" dynamodbav:"code,omitempty"`
	ConnectedMemberIds []string        `json:"-" dynamodbav:"connectedMemberIds,omitempty"`
}
