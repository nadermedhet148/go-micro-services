package services

import (
	"errors"
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/rabbitmq"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RechargeWalletCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "recharge_wallet_total",
			Help: "The total number of processed transactions",
		},
	)
	RechargeWalletDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "recharge_wallet_duration_seconds",
			Help:    "Duration of transaction processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(RechargeWalletCounter)
	prometheus.MustRegister(RechargeWalletDuration)
}

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
	start := time.Now()
	defer func() {
		RechargeWalletCounter.Inc()
		duration := time.Since(start).Seconds()
		RechargeWalletDuration.Observe(duration)

	}()

	wallet := service.repositories.Get(req.Region, req.WALLET_ID)
	if wallet.ID == 0 {
		return errors.New("wallet not found")
	}

	req.Region = wallet.Region
	err := service.paymentProducer.PublishPayment(req)

	return err
}
