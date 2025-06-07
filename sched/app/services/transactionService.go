package services

import (
	"errors"
	"time"

	"encoding/json"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	opsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "transaction_ops_processed_total",
			Help: "The total number of processed transactions",
		},
	)
	transactionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "transaction_duration_seconds",
			Help:    "Duration of transaction processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(opsProcessed)
	prometheus.MustRegister(transactionDuration)
}

type TransactionService interface {
	UpdateTransaction(Transaction entity.TransactionUpdateRequest) error
}

type transactionService struct {
	repositories repositories.TransactionRepository
}

const (
	DEBIT  = "debit"
	CREDIT = "credit"
)

func NewTransactionService(repository repositories.TransactionRepository, paymentEventProducer kafka.PaymentEventProducer) TransactionService {
	return &transactionService{
		repositories: repository,
	}
}

func (service *transactionService) UpdateTransaction(req entity.TransactionUpdateRequest) error {
	start := time.Now()
	defer func() {
		updateProcessed.Inc()
		duration := time.Since(start).Seconds()
		transactionUpdateDuration.Observe(duration)
	}()

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
