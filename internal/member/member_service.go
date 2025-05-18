package member

import (
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
)

type Service interface {
	UpdateMember(request *UpdateMemberReq, transaction *util.DynamoDBWriteTransaction) (*Member, util.ApplicationError)
	GetMember(socialId string, provider string) (*Member, util.ApplicationError)
}

type ServiceImpl struct {
	repo Repository
}

func NewMemberService(repo Repository) *ServiceImpl {
	return &ServiceImpl{repo: repo}
}

func (svc *ServiceImpl) UpdateMember(request *UpdateMemberReq, transaction *util.DynamoDBWriteTransaction) (*Member, util.ApplicationError) {
	if transaction != nil {
		updatedMember, err := svc.repo.UpdateMemberTransaction(request, transaction)
		if err != nil {
			log.Printf("failed to get UpdateMemberTransaction. error:%s", err.Error())
			return nil, util.DBUpdateError{}
		}
		return updatedMember, nil
	}
	updatedMember, err := svc.repo.UpdateMember(request)
	if err != nil {
		log.Printf("failed to update member. req: %+v, error: %s", request, err.Error())
		return nil, util.DBUpdateError{}
	}
	return updatedMember, nil
}

func (svc *ServiceImpl) GetMember(socialId string, provider string) (*Member, util.ApplicationError) {
	member, err := svc.repo.FindBySocialIdAndProvider(socialId, provider)
	if err != nil {
		log.Printf("failed to get member with socialId: %d, provider: %s \nerror: %s", socialId, provider, err.Error())
		return nil, util.DBReadError{}
	}
	return member, nil
}
