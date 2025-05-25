package services

import (
	"errors"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/rabbitmq"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
)

type WalletService interface {
	CerateWallet(wallet entity.Wallet) (int, error)
	RechargeWallet(req entity.WalletRechargeRequest) error
}

type walletService struct {
	repositories    repositories.WalletRepository
	paymentProducer *rabbitmq.PaymentProducer
}

func NewWalletService(repository repositories.WalletRepository, paymentProducer *rabbitmq.PaymentProducer) WalletService {
	return &walletService{
		repositories:    repository,
		paymentProducer: paymentProducer,
	}
}

func (service *walletService) CerateWallet(wallet entity.Wallet) (int, error) {
	id, err := service.repositories.Save(wallet.Region, wallet)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (service *walletService) RechargeWallet(req entity.WalletRechargeRequest) error {
	wallet := service.repositories.Get(req.Region, req.WALLET_ID)
	if wallet.ID == 0 {
		return errors.New("wallet not found")
	}

	err := service.paymentProducer.PublishPayment(req)

	return err
}
