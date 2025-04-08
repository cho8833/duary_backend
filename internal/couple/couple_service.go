package couple

import (
	"github.com/cho8833/duary_lambda/internal/util"
	"log"
	"math/rand"
	"time"
)

type Service interface {
	CreateCouple(req *CreateCoupleReq, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError)
	UpdateCouple(couple *Couple, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError)
	FindByCoupleCode(coupleCode *string) (*Couple, util.ApplicationError)
}

type ServiceImpl struct {
	repository Repository
}

func NewCoupleService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (svc *ServiceImpl) CreateCouple(req *CreateCoupleReq, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError) {
	couple := &Couple{
		RelationDate: &req.RelationDate,
		Code:         svc.generateCoupleCode(),
	}

	if transaction != nil {
		input, err := svc.repository.GetSaveCoupleTransaction(couple)
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

func (svc *ServiceImpl) UpdateCouple(couple *Couple, transaction *util.DynamoDBWriteTransaction) (*Couple, util.ApplicationError) {
	if transaction != nil {
		writeTransaction, err := svc.repository.GetSaveCoupleTransaction(couple)
		if err != nil {
			return nil, util.DBUpdateError{}
		}
		transaction.AddTransaction(writeTransaction)
		return couple, nil
	} else {
		updatedCouple, err := svc.repository.SaveCouple(couple)
		if err != nil {
			return nil, util.DBUpdateError{}
		}
		return updatedCouple, nil
	}

}

func (svc *ServiceImpl) FindByCoupleCode(coupleCode *string) (*Couple, util.ApplicationError) {
	couples, err := svc.repository.FindByCoupleCode(coupleCode)
	if err != nil {
		return nil, util.CoupleNotFound{}
	}
	if len(couples) != 1 {
		log.Printf("Couple with coupleCode: %s invalid. Found %d", *coupleCode, len(couples))
		return nil, util.InternalServerError{}
	}
	return &couples[0], nil
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
