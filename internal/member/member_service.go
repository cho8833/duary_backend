package member

import (
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

type Service interface {
	UpdateMember(request *UpdateMemberReq, transaction *util.DynamoDBWriteTransaction) (*Member, util.ApplicationError)
	GetMember(socialId int64, provider string) (*Member, util.ApplicationError)
}

type ServiceImpl struct {
	repo Repository
}

func NewMemberService(repo Repository) *ServiceImpl {
	return &ServiceImpl{repo: repo}
}

func (svc *ServiceImpl) UpdateMember(request *UpdateMemberReq, transaction *util.DynamoDBWriteTransaction) (*Member, util.ApplicationError) {
	if transaction != nil {
		targetMember, err := svc.repo.FindBySocialIdAndProvider(request.SocialId, request.Provider)
		if err != nil {
			return nil, util.DBReadError{}
		}
		update, err := svc.repo.GetUpdateMemberTransaction(request)
		if err != nil {
			log.Printf("failed to get UpdateMemberTransaction. error:%s", err.Error())
			return nil, util.DBUpdateError{}
		}
		transaction.AddTransaction(update)
		svc.updateField(targetMember, request)
		return targetMember, nil
	}
	member, err := svc.repo.UpdateMember(request)
	if err != nil {
		log.Printf("failed to update member. req: %+v, error: %s", request, err.Error())
		return nil, util.DBUpdateError{}
	}
	return member, nil
}

func (svc *ServiceImpl) GetMember(socialId int64, provider string) (*Member, util.ApplicationError) {
	member, err := svc.repo.FindBySocialIdAndProvider(socialId, provider)
	if err != nil {
		log.Printf("failed to get member with socialId: %d, provider: %s \nerror: %s", socialId, provider, err.Error())
		return nil, util.DBReadError{}
	}
	return member, nil
}

func (svc *ServiceImpl) updateField(target *Member, req *UpdateMemberReq) {
	if req.Character != nil {
		target.Character = req.Character
	}
	if req.Name != nil {
		target.Name = req.Name
	}
	if req.Birthday != nil {
		target.Birthday = req.Birthday
	}
	if req.Gender != nil {
		target.Gender = req.Gender
	}
	if req.FcmToken != nil {
		target.FcmToken = req.FcmToken
	}
	if req.Email != nil {
		target.Email = req.Email
	}
	if req.AccessToken != nil {
		target.AccessToken = req.AccessToken
	}
	if req.CoupleId != nil {
		target.CoupleId = req.CoupleId
	}
}
