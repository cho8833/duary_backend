package common

import (
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/cho8833/duary_lambda/internal/couple"
	"github.com/cho8833/duary_lambda/internal/member"
	"time"
)

type StartDuaryReq struct {
	Name         *string    `json:"name"`
	Birthday     *time.Time `json:"birthday"`
	RelationDate *time.Time `json:"relationDate"`
	MyCharacter  *string    `json:"myCharacter"`
	Provider     string
	SocialId     int64
}

type StartDuaryRes struct {
	Member *member.Member          `json:"member"`
	Couple *couple.Couple          `json:"couple"`
	Token  *jwtutil.ApplicationJWT `json:"token"`
}

type ConnectCoupleReq struct {
	CoupleCode *string `json:"coupleCode"`
}

type CheckConnectedRes struct {
	Token  *jwtutil.ApplicationJWT `json:"token"`
	Couple *couple.Couple          `json:"couple"`
}
