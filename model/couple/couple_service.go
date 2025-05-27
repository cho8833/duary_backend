package couple

import (
	"github.com/cho8833/duary_lambda/model"
	"github.com/cho8833/duary_lambda/shared"
	uuid2 "github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type Service interface {
	CreateCouple(req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction, id ...string) (*Couple, shared.ApplicationError)
	UpdateCouple(couple *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError)
	FindByCoupleCode(coupleCode *string) (*Couple, shared.ApplicationError)
	FindById(coupleId string) (*Couple, shared.ApplicationError)
	DeleteCouple(id string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError
	CreateCoupleWithId(id string, req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError)
	GenerateUID() *string
}

type ServiceImpl struct {
	repository Repository
}

func NewCoupleService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (svc *ServiceImpl) CreateCouple(req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction, id ...string) (*Couple, shared.ApplicationError) {
	var coupleId *string
	if len(id) == 0 {
		coupleId = svc.GenerateUID()
	} else {
		coupleId = &id[0]
	}
	couple := &Couple{
		Id:           coupleId,
		RelationDate: &req.RelationDate,
		Code:         svc.generateCoupleCode(),
		Members:      req.Members,
	}

	if transaction != nil {
		input, err := svc.repository.SaveCoupleTransaction(couple)
		if err != nil {
			log.Printf("failed to get saveCoupleTransaction. error:%s", err.Error())
			return nil, shared.DBSaveError{}
		}
		transaction.AddTransaction(input)
		return couple, nil
	} else {
		couple, err := svc.repository.SaveCouple(couple)
		if err != nil {
			log.Printf("failed to save couple\nreq: %+v\nerror:%s", req, err.Error())
			return nil, shared.DBSaveError{}
		}
		return couple, nil
	}
}

func (svc *ServiceImpl) UpdateCouple(req *UpdateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError) {
	now := time.Now().UTC()
	relationDate := req.RelationDate.UTC()
	if relationDate.After(now) {
		return nil, shared.NewCustomApplicationError("처음 만난 날을 다시 설정해주세요")
	}
	if transaction != nil {
		updatedCouple, err := svc.repository.UpdateCoupleTransaction(req, transaction)
		if err != nil {
			log.Printf("failed to update couple. req: %+v, error: %s", req, err.Error())
			return nil, shared.DBUpdateError{}
		}
		return updatedCouple, nil
	} else {
		updatedCouple, err := svc.repository.UpdateCouple(req)
		if err != nil {
			log.Printf("failed to update couple. req: %+v, error: %s", req, err.Error())
			return nil, shared.DBUpdateError{}
		}
		return updatedCouple, nil
	}
}

func (svc *ServiceImpl) DeleteCouple(id string, transaction *model.DynamoDBWriteTransaction) shared.ApplicationError {
	deleteTransaction, err := svc.repository.DeleteCoupleTransaction(id)
	if err != nil {
		return shared.DBDeleteError{}
	}
	transaction.AddTransaction(deleteTransaction)
	return nil
}

func (svc *ServiceImpl) FindById(coupleId string) (*Couple, shared.ApplicationError) {
	couple, err := svc.repository.FindById(&coupleId)
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

func (svc *ServiceImpl) CreateCoupleWithId(id string, req *CreateCoupleReq, transaction *model.DynamoDBWriteTransaction) (*Couple, shared.ApplicationError) {
	return svc.CreateCouple(req, transaction, id)
}

func (svc *ServiceImpl) GenerateUID() *string {
	uuid := uuid2.New().String()
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
