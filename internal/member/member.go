package member

import "time"

type Member struct {
	Name        *string    `json:"name" dynamodbav:"name"`
	Birthday    *time.Time `json:"birthday" dynamodbav:"birthday"`
	Gender      *string    `json:"gender" dynamodbav:"gender"` // man, woman, other
	FcmToken    *string    `json:"fcmToken" dynamodbav:"fcmToken"`
	AccessToken *string    `json:"accessToken" dynamodbav:"accessToken"`
	Provider    string     `json:"provider" dynamodbav:"provider"`
	SocialId    int64      `json:"socialId" dynamodbav:"socialId"`
	Email       *string    `json:"email" dynamodbav:"email"`
	CoupleId    *string    `json:"coupleId" dynamodbav:"coupleId"`
	Character   *string    `json:"character" dynamodbav:"character"`

	// SOLO, COUPLE
	Status *string `json:"status" dynamodbav:"status"`
}

type UpdateMemberReq struct {
	Name        *string    `dynamodbav:"name"`
	Birthday    *time.Time `dynamodbav:"birthday"`
	Gender      *string    `dynamodbav:"gender"` // man, woman, other
	FcmToken    *string    `dynamodbav:"fcmToken"`
	AccessToken *string    `dynamodbav:"accessToken"`
	Provider    string     `dynamodbav:"provider"`
	SocialId    int64      `dynamodbav:"socialId"`
	Email       *string    `dynamodbav:"email"`
	CoupleId    *string    `dynamodbav:"coupleId"`
	Character   *string    `dynamodbav:"character"`
}
