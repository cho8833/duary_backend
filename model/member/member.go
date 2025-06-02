package member

import (
	"github.com/cho8833/duary_lambda/model/event"
	"time"
)

const DefaultAlarmOffset = event.Min15

type Member struct {
	Name        *string           `json:"name" dynamodbav:"name"`
	Birthday    *time.Time        `json:"birthday" dynamodbav:"birthday"`
	Gender      *string           `json:"gender" dynamodbav:"gender"` // man, woman, other
	FcmToken    *string           `json:"fcmToken" dynamodbav:"fcmToken"`
	AccessToken *string           `json:"accessToken" dynamodbav:"accessToken"`
	Provider    string            `json:"provider" dynamodbav:"provider"`
	SocialId    string            `json:"socialId" dynamodbav:"socialId"`
	Email       *string           `json:"email" dynamodbav:"email"`
	CoupleId    *string           `json:"coupleId" dynamodbav:"coupleId"`
	Character   *string           `json:"character" dynamodbav:"character"`
	LoverAlarm  event.AlarmOffset `json:"loverAlarm" dynamodbav:"loverAlarm"`
	MyAlarm     event.AlarmOffset `json:"myAlarm" dynamodbav:"myAlarm"`
}

func (m *Member) ApplyFrom(req UpdateMemberReq) {
	if req.Name != nil {
		m.Name = req.Name
	}
	if req.Birthday != nil {
		m.Birthday = req.Birthday
	}
	if req.FcmToken != nil {
		m.FcmToken = req.FcmToken
	}
	if req.AccessToken != nil {
		m.AccessToken = req.AccessToken
	}
	if req.CoupleId != nil {
		m.CoupleId = req.CoupleId
	}
	if req.Character != nil {
		m.Character = req.Character
	}
	if req.LoverAlarm != nil {
		m.LoverAlarm = *req.LoverAlarm
	}
	if req.MyAlarm != nil {
		m.MyAlarm = *req.MyAlarm
	}
}

func (m *Member) GetId() string {
	return m.SocialId + "-" + m.Provider
}

func FromSaveMemberReq(req *SaveMemberReq) *Member {
	return &Member{
		FcmToken:    req.FcmToken,
		AccessToken: req.AccessToken,
		Character:   req.Character,
		Name:        req.Name,
		CoupleId:    req.CoupleId,
		Birthday:    req.Birthday,
		SocialId:    req.SocialId,
		Provider:    req.Provider,
		LoverAlarm:  req.LoverAlarm,
		MyAlarm:     req.MyAlarm,
	}
}

type UpdateMemberReq struct {
	Name        *string            `dynamodbav:"name"`
	Birthday    *time.Time         `dynamodbav:"birthday"`
	FcmToken    *string            `dynamodbav:"fcmToken"`
	AccessToken *string            `dynamodbav:"accessToken"`
	Provider    string             `dynamodbav:"provider"`
	SocialId    string             `dynamodbav:"socialId"`
	CoupleId    *string            `dynamodbav:"coupleId"`
	Character   *string            `dynamodbav:"character"`
	LoverAlarm  *event.AlarmOffset `dynamodbav:"loverAlarm"`
	MyAlarm     *event.AlarmOffset `dynamodbav:"myAlarm"`
}

type SaveMemberReq struct {
	Name        *string           `json:"name"`
	Birthday    *time.Time        `json:"birthday"`
	FcmToken    *string           `json:"fcmToken"`
	AccessToken *string           `json:"accessToken"`
	Provider    string            `json:"provider"`
	SocialId    string            `json:"socialId"`
	CoupleId    *string           `json:"coupleId"`
	Character   *string           `json:"character"`
	LoverAlarm  event.AlarmOffset `json:"loverAlarm"`
	MyAlarm     event.AlarmOffset `json:"myAlarm"`
}
