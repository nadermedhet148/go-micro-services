package services

import (
	"errors"

	entity "github.com/coroo/go-starter/app/entity"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type TransactionService interface {
	CerateRechargeTransaction(Transaction entity.WalletRechargeRequest) (int, error)
	updateTransaction(Transaction entity.TransactionUpdateRequest) error
}

type transactionService struct {
	repositories repositories.TransactionRepository
}

func NewTransactionService(repository repositories.TransactionRepository) TransactionService {
	return &transactionService{
		repositories: repository,
	}
}

func (service *transactionService) CerateRechargeTransaction(req entity.WalletRechargeRequest) (int, error) {

	Transaction := entity.Transaction{
		WALLET_ID:  req.WALLET_ID,
		AMOUNT:     req.AMOUNT,
		STATUS:     "pending",
		REF_NUMBER: uuid.New().String(),
	}

	id, err := service.repositories.Save(Transaction)
	if err != nil {
		return 0, err
	}
	// TBD: send to notification
	return id, nil
}
func (service *transactionService) updateTransaction(req entity.TransactionUpdateRequest) error {
	Transaction := service.repositories.GetByRefNumber(req.REF_NUMBER)
	if Transaction.ID == 0 {
		return errors.New("transaction not found")
	}
	if Transaction.STATUS != "pending" {
		return errors.New("transaction already completed")
	}
	Transaction.STATUS = req.STATUS
	err := service.repositories.Update(Transaction)
	if err != nil {
		return err
	}
	return nil

}
