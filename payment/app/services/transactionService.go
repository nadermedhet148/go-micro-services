package services

import (
	"errors"
	"time"

	"encoding/json"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/kafka"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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
	updateProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "transaction_update_processed_total",
			Help: "The total number of processed transaction updates",
		})
	transactionUpdateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "transaction_update_duration_seconds",
			Help:    "Duration of transaction update processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(opsProcessed)
	prometheus.MustRegister(transactionDuration)
	prometheus.MustRegister(updateProcessed)
	prometheus.MustRegister(transactionUpdateDuration)
}

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

	start := time.Now()
	defer func() {
		opsProcessed.Inc()
		duration := time.Since(start).Seconds()
		transactionDuration.Observe(duration)

	}()

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
