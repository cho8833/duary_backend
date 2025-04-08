package common

import (
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"time"
)

type InitDuaryInfoReq struct {
	Name           *string    `json:"name"`
	Birthday       *time.Time `json:"birthday"`
	RelationDate   *time.Time `json:"relationDate"`
	OtherCharacter *string    `json:"otherCharacter"`
	MyCharacter    *string    `json:"myCharacter"`
	Provider       string
	SocialId       int64
}

type InitDuaryInfoRes struct {
	Member *member.Member `json:"member"`
	Couple *couple.Couple `json:"couple"`
}

type ConnectCoupleReq struct {
	CoupleCode *string `json:"coupleCode"`
}
