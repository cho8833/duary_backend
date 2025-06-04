package member

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	"log"
)

type Service interface {
	UpdateMember(request *UpdateMemberReq) (*Member, shared.ApplicationError)
	UpdateMemberTransaction(request *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction) (*Member, shared.ApplicationError)
	GetMember(socialId string, provider string) (*Member, shared.ApplicationError)
	FindById(socialId string, provider string) (*Member, shared.ApplicationError)
	SaveMember(request *SaveMemberReq) (*Member, shared.ApplicationError)
	UpdateFcmToken(socialId string, provider string, fcmToken *string) (*Member, shared.ApplicationError)
}

type ServiceImpl struct {
	repo Repository
}

func NewMemberService(repo Repository) *ServiceImpl {
	return &ServiceImpl{repo: repo}
}

func (svc *ServiceImpl) FindById(socialId string, provider string) (*Member, shared.ApplicationError) {
	member, err := svc.repo.FindBySocialIdAndProvider(socialId, provider)
	if temp := new(types.ResourceNotFoundException); errors.As(err, &temp) {
		return nil, shared.UserNotFound{}
	}

	if err != nil {
		log.Printf(err.Error())
		return nil, shared.DBReadError{}
	}
	return member, nil
}

func (svc *ServiceImpl) SaveMember(request *SaveMemberReq) (*Member, shared.ApplicationError) {
	newMember := FromSaveMemberReq(request)

	result, err := svc.repo.SaveMember(newMember)
	if err != nil {
		log.Printf(err.Error())
		return nil, shared.DBSaveError{}
	}
	return result, nil
}

func (svc *ServiceImpl) UpdateMember(request *UpdateMemberReq) (*Member, shared.ApplicationError) {
	
	updatedMember, err := svc.repo.UpdateNonNil(request)
	if err != nil {
		log.Printf("failed to update member. req: %+v, error: %s", request, err.Error())
		return nil, shared.DBUpdateError{}
	}
	return updatedMember, nil
}

func (svc *ServiceImpl) UpdateMemberTransaction(request *UpdateMemberReq, transaction *model.DynamoDBWriteTransaction) (*Member, shared.ApplicationError) {
	updatedMember, err := svc.repo.UpdateNonNilTransaction(request, transaction)
	if err != nil {
		log.Printf("failed to get UpdateNonNilTransaction. error:%s", err.Error())
		return nil, shared.DBUpdateError{}
	}
	return updatedMember, nil
}

func (svc *ServiceImpl) GetMember(socialId string, provider string) (*Member, shared.ApplicationError) {
	member, err := svc.repo.FindBySocialIdAndProvider(socialId, provider)
	if err != nil {
		log.Printf("failed to get member with socialId: %d, provider: %s \nerror: %s", socialId, provider, err.Error())
		return nil, shared.DBReadError{}
	}
	return member, nil
}
func (svc *ServiceImpl) UpdateFcmToken(socialId string, provider string, fcmToken *string) (*Member, shared.ApplicationError) {
	member, err := svc.repo.UpdateFcmToken(socialId, provider, fcmToken)
	if err != nil {
		log.Printf("failed to update member. req: %+v, error: %s", member, err.Error())
		return nil, shared.DBUpdateError{}
	}
	return member, nil
}
