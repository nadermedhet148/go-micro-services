package services

import (
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/rabbitmq"
	"github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	transactionsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "transaction_ops_processed_total",
			Help: "The total number of processed transactions",
		},
	)
	opsDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "transaction_duration_seconds",
			Help:    "Duration of transaction processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(transactionsProcessed)
	prometheus.MustRegister(opsDuration)
}

type TransactionService interface {
	RunExpiredTransactionCleanup() error
}

type transactionService struct {
	repositories    repositories.TransactionRepository
	paymentProducer rabbitmq.PaymentProducer
}

func NewTransactionService(repository repositories.TransactionRepository, paymentProducer rabbitmq.PaymentProducer) TransactionService {
	return &transactionService{
		repositories:    repository,
		paymentProducer: paymentProducer,
	}
}

func (service *transactionService) RunExpiredTransactionCleanup() error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		opsDuration.Observe(duration)

	}()

	for bucket := 0; bucket < 10; bucket++ {
		pendingTrxs := service.repositories.GePendingTrxs(bucket)

		for _, trx := range pendingTrxs {
			transactionsProcessed.Inc()
			service.paymentProducer.PublishPayment(entity.TransactionUpdateRequest{
				REF_NUMBER: trx.REF_NUMBER,
				STATUS:     "expired",
				REGION:     trx.REGION,
			})
		}
	}
	return nil
}
