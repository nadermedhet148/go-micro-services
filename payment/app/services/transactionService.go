package services

import (
	"errors"

	"encoding/json"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/kafka"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type TransactionService interface {
	CerateRechargeTransaction(Transaction entity.WalletRechargeRequest) (int, error)
	UpdateTransaction(Transaction entity.TransactionUpdateRequest) error
}

type transactionService struct {
	repositories         repositories.TransactionRepository
	paymentEventProducer kafka.PaymentEventProducer
}

const (
	DEBIT  = "debit"
	CREDIT = "credit"
)

func NewTransactionService(repository repositories.TransactionRepository, paymentEventProducer kafka.PaymentEventProducer) TransactionService {
	return &transactionService{
		repositories:         repository,
		paymentEventProducer: paymentEventProducer,
	}
}

func (service *transactionService) CerateRechargeTransaction(req entity.WalletRechargeRequest) (int, error) {

	Transaction := entity.Transaction{
		WALLET_ID:  req.WALLET_ID,
		AMOUNT:     req.AMOUNT,
		STATUS:     "pending",
		REGION:     req.REGION,
		TYPE:       DEBIT,
		REF_NUMBER: uuid.New().String(),
	}

	id, err := service.repositories.Save(req.REGION, Transaction)
	if err != nil {
		return 0, err
	}
	event := kafka.PaymentEvent{
		WALLET_ID: Transaction.WALLET_ID,
		AMOUNT:    Transaction.AMOUNT,
		STATUS:    "pending",
		TYPE:      Transaction.TYPE,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return 0, err
	}
	service.paymentEventProducer.PushToTopicWithPartition(eventBytes)
	return id, nil
}
func (service *transactionService) UpdateTransaction(req entity.TransactionUpdateRequest) error {
	Transaction := service.repositories.GetByRefNumber(req.REGION, req.REF_NUMBER)
	if Transaction.ID == 0 {
		return errors.New("transaction not found")
	}
	if Transaction.STATUS != "pending" {
		return errors.New("transaction already completed")
	}
	Transaction.STATUS = req.STATUS
	err := service.repositories.Update(req.REGION, Transaction)
	if err != nil {
		return err
	}
	event := kafka.PaymentEvent{
		WALLET_ID: Transaction.WALLET_ID,
		AMOUNT:    Transaction.AMOUNT,
		STATUS:    Transaction.STATUS,
		TYPE:      Transaction.TYPE,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	service.paymentEventProducer.PushToTopicWithPartition(eventBytes)

	return nil

}
