package usecases

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/externalServices"
	"github.com/coroo/go-starter/app/rabbitmq"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/coroo/go-starter/config"

	"github.com/go-redsync/redsync/v4"
)

type TicketService interface {
	SaveTicket(Ticket entity.Ticket) (int, error)
	SaveTicketWithLock(Ticket entity.Ticket) (int, error)
	SaveTicketWithDLock(Ticket entity.Ticket) (int, error)
	GetAllTickets() []entity.Ticket
	TicketCreateWithHttpTransaction(Ticket entity.Ticket) (int, error)
	TicketCreateWithRabbitMq(Ticket entity.Ticket) (int, error)
}

type ticketService struct {
	repositories repositories.TicketRepository
	mutexes      sync.Map
}

func NewTicketService(repository repositories.TicketRepository) TicketService {
	return &ticketService{
		repositories: repository,
		mutexes:      sync.Map{},
	}
}

func (usecases *ticketService) GetAllTickets() []entity.Ticket {
	return usecases.repositories.GetAllTickets()
}

func (usecases *ticketService) SaveTicket(ticket entity.Ticket) (int, error) {
	var createdTicket = usecases.repositories.GetTicketBySlotId(ticket.SLOT_ID)
	if createdTicket != (entity.Ticket{}) {
		return 0, errors.New("slot already taken")
	}
	return usecases.repositories.Save(ticket)
}

func (usecases *ticketService) SaveTicketWithLock(ticket entity.Ticket) (int, error) {

	key := ticket.REF_NUMBER
	mutexInterface, _ := usecases.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex := mutexInterface.(*sync.Mutex)

	mutex.Lock()
	id, err := usecases.SaveTicket(ticket)
	defer mutex.Unlock()
	return id, err
}

func (usecases *ticketService) SaveTicketWithDLock(ticket entity.Ticket) (int, error) {

	rs := redsync.New(config.ConnectRedis())

	mutexname := "t-" + ticket.REF_NUMBER + "-mutex"
	mutex := rs.NewMutex(mutexname)

	if err := mutex.Lock(); err != nil {
		return 0, errors.New("lock failed")
	}
	id, err := usecases.SaveTicket(ticket)
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return 0, errors.New("unlock failed")
	}
	return id, err
}

func (usecases *ticketService) TicketCreateWithHttpTransaction(ticket entity.Ticket) (int, error) {

	paymentService := externalServices.PaymentService{}

	ticket.Status = "payment_pending"
	id, err := usecases.SaveTicket(ticket)
	ticket.ID = id
	if err == nil {
		response, err := paymentService.SendPayment(externalServices.PaymentRequest{
			REF_NUMBER: ticket.REF_NUMBER,
		})
		if response.Status != "success" || err != nil {
			ticket.Status = "failed"
			usecases.repositories.UpdateTicket(ticket)
			return 0, errors.New("payment failed")
		}
	}
	ticket.Status = "success"
	err = usecases.repositories.UpdateTicket(ticket)
	return id, err
}

func (usecases *ticketService) TicketCreateWithRabbitMq(ticket entity.Ticket) (int, error) {

	ticket.Status = "payment_pending"
	id, err := usecases.SaveTicket(ticket)
	ticket.ID = id
	if err == nil {
		producer, err := rabbitmq.NewPaymentProducer()
		if err != nil {
			return 0, errors.New("payment failed")
		}
		defer producer.Close()
		ticketJSON, err := json.Marshal(ticket)
		if err != nil {
			return 0, errors.New("failed to marshal ticket to JSON")
		}
		err = producer.PublishPayment(string(ticketJSON))
		if err != nil {
			return 0, errors.New("payment failed")
		}
	}
	return id, err
}
