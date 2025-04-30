package couple

import (
	"github.com/cho8833/duary_lambda/internal/util"
	uuid2 "github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type Service interface {
	CreateCouple(req *CreateCoupleReq, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError)
	UpdateCouple(couple *UpdateCoupleReq, transaction *util.DynamoDBWriteTransaction) util.ApplicationError
	FindByCoupleCode(coupleCode *string) (*Couple, util.ApplicationError)
	FindById(coupleId string) (*Couple, util.ApplicationError)
	DeleteCouple(id string, transaction *util.DynamoDBWriteTransaction) util.ApplicationError
}

type ServiceImpl struct {
	repository Repository
}

func NewCoupleService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (svc *ServiceImpl) CreateCouple(req *CreateCoupleReq, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError) {
	couple := &Couple{
		Id:           svc.generateUID(),
		RelationDate: &req.RelationDate,
		Code:         svc.generateCoupleCode(),
		Members:      req.Members,
	}

	if transaction != nil {
		input, err := svc.repository.SaveCoupleTransaction(couple)
		if err != nil {
			log.Printf("failed to get saveCoupleTransaction. error:%s", err.Error())
			return nil, util.DBSaveError{}
		}
		transaction.AddTransaction(input)
		return couple, nil
	} else {
		couple, err := svc.repository.SaveCouple(couple)
		if err != nil {
			log.Printf("failed to save couple\nreq: %+v\nerror:%s", req, err.Error())
			return nil, util.DBSaveError{}
		}
		return couple, nil
	}
}

func (svc *ServiceImpl) UpdateCouple(couple *UpdateCoupleReq, transaction *util.DynamoDBWriteTransaction) util.ApplicationError {
	writeTransaction, err := svc.repository.UpdateCoupleTransaction(couple)
	if err != nil {
		return util.DBUpdateError{}
	}
	transaction.AddTransaction(writeTransaction)
	return nil
}

func (svc *ServiceImpl) DeleteCouple(id string, transaction *util.DynamoDBWriteTransaction) util.ApplicationError {
	deleteTransaction, err := svc.repository.DeleteCoupleTransaction(id)
	if err != nil {
		return util.DBDeleteError{}
	}
	transaction.AddTransaction(deleteTransaction)
	return nil
}

func (svc *ServiceImpl) FindById(coupleId string) (*Couple, util.ApplicationError) {
	couple, err := svc.repository.FindById(&coupleId)
	if err != nil {
		return nil, util.CoupleNotFound{}
	}
	return couple, nil
}

func (svc *ServiceImpl) FindByCoupleCode(coupleCode *string) (*Couple, util.ApplicationError) {
	couples, err := svc.repository.FindByCoupleCode(coupleCode)
	if err != nil {
		return nil, util.DBReadError{}
	}
	if len(couples) > 1 {
		log.Printf("Couple with coupleCode: %s invalid. Found %d", *coupleCode, len(couples))
		return nil, util.InternalServerError{}
	} else if len(couples) == 0 {
		return nil, util.CoupleNotFound{}
	} else {
		return &couples[0], nil
	}

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

func (svc *ServiceImpl) generateUID() *string {
	uuid := uuid2.New().String()
	return &uuid
}
