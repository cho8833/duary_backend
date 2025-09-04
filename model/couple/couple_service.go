package couple

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	uuid2 "github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type Service interface {
	Create(req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction, id ...string) (*Couple, shared.ApplicationError)
	Update(couple *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError)
	FindByCoupleCode(coupleCode *string) (*Couple, shared.ApplicationError)
	FindById(coupleId string) (*Couple, shared.ApplicationError)
	Delete(id string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError
	CreateWithId(id string, req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError)
	GenerateUID() *string
}

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (svc *ServiceImpl) Create(req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction, id ...string) (*Couple, shared.ApplicationError) {
	var coupleId *string
	if len(id) == 0 {
		coupleId = svc.GenerateUID()
	} else {
		coupleId = &id[0]
	}
	couple := &Couple{
		Id:                 coupleId,
		RelationDate:       req.RelationDate,
		Code:               svc.generateCoupleCode(),
		Members:            req.Members,
		ConnectedMemberIds: req.ConnectedMemberIds,
	}

	if transaction != nil {
		input, err := svc.repository.SaveTransaction(couple)
		if err != nil {
			log.Printf("failed to get saveCoupleTransaction. error:%s", err.Error())
			return nil, shared.DBSaveError{}
		}
		transaction.AddTransaction(input)
		return couple, nil
	} else {
		couple, err := svc.repository.Save(couple)
		if err != nil {
			log.Printf("failed to save couple\nreq: %+v\nerror:%s", req, err.Error())
			return nil, shared.DBSaveError{}
		}
		return couple, nil
	}
}

func (svc *ServiceImpl) Update(req *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError) {
	if req.RelationDate != nil {
		now := time.Now().UTC()
		relationDate := *req.RelationDate
		if relationDate.After(now) {
			return nil, shared.ValidateError{Message: "처음 만난 날을 다시 설정해주세요"}
		}
		relationDate = time.Date(relationDate.Year(), relationDate.Month(), relationDate.Day(), 0, 0, 0, 0, relationDate.Location())
		req.RelationDate = &relationDate
	}
	if transaction != nil {
		updatedCouple, err := svc.repository.UpdateTransaction(req, transaction)
		if err != nil {
			log.Printf("failed to update couple. req: %+v, error: %s", req, err.Error())
			return nil, shared.DBUpdateError{}
		}
		return updatedCouple, nil
	} else {
		updatedCouple, err := svc.repository.Update(req)
		if err != nil {
			log.Printf("failed to update couple. req: %+v, error: %s", req, err.Error())
			return nil, shared.DBUpdateError{}
		}
		return updatedCouple, nil
	}
}

func (svc *ServiceImpl) Delete(id string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError {
	deleteTransaction, err := svc.repository.DeleteTransaction(id)
	if err != nil {
		return shared.DBDeleteError{}
	}
	transaction.AddTransaction(deleteTransaction)
	return nil
}

func (svc *ServiceImpl) FindById(coupleId string) (*Couple, shared.ApplicationError) {
	couple, err := svc.repository.FindById(&coupleId)
	if temp := new(types.ResourceNotFoundException); errors.As(err, &temp) {
		return nil, shared.CoupleNotFound{}
	}

	if err != nil {
		return nil, shared.DBReadError{}
	}
	return couple, nil
}

func (svc *ServiceImpl) FindByCoupleCode(coupleCode *string) (*Couple, shared.ApplicationError) {
	couples, err := svc.repository.FindByCoupleCode(coupleCode)
	if err != nil {
		return nil, shared.DBReadError{}
	}
	if len(couples) > 1 {
		log.Printf("Couple with coupleCode: %s invalid. Found %d", *coupleCode, len(couples))
		return nil, shared.InternalServerError{}
	} else if len(couples) == 0 {
		return nil, shared.CoupleNotFound{}
	} else {
		return &couples[0], nil
	}

}

func (svc *ServiceImpl) CreateWithId(id string, req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError) {
	return svc.Create(req, transaction, id)
}

func (svc *ServiceImpl) GenerateUID() *string {
	uuid := uuid2.New().String()[:32]
	return &uuid
}

func (svc *ServiceImpl) generateCoupleCode() *string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) // 현재 시간을 시드로 설정
	b := make([]byte, 9)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	result := string(b)
	return &result
}
