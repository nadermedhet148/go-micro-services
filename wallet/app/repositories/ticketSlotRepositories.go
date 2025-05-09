package repositories

import (
	"os"
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	Save(ticket entity.Ticket) (int, error)
	UpdateTicket(ticket entity.Ticket) error
	DeleteTicket(ticket entity.Ticket) error
	GetAllTickets() []entity.Ticket
	GetTicket(id string) entity.Ticket
	GetTicketBySlotId(slotId int) entity.Ticket
}

type TicketDatabase struct {
	connection *gorm.DB
}

func NewTicketRepository() TicketRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}
	if os.Getenv("DB_HOST") != "" {
		db.AutoMigrate(&entity.Ticket{})
	} else {
		db.AutoMigrate(&entity.TicketTesting{})
	}
	return &TicketDatabase{
		connection: db,
	}
}

func (db *TicketDatabase) Save(Ticket entity.Ticket) (int, error) {
	data := &Ticket
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}

func (db *TicketDatabase) UpdateTicket(Ticket entity.Ticket) error {
	data := &Ticket
	data.UpdatedAt = time.Now()
	data.CreatedAt = time.Now()
	db.connection.Save(data)
	return nil
}

func (db *TicketDatabase) DeleteTicket(Ticket entity.Ticket) error {
	db.connection.Delete(&Ticket)
	return nil
}

func (db *TicketDatabase) GetAllTickets() []entity.Ticket {
	var Tickets []entity.Ticket
	db.connection.Preload(clause.Associations).Find(&Tickets)
	return Tickets
}

func (db *TicketDatabase) GetTicket(ref_number string) entity.Ticket {
	var Ticket entity.Ticket
	db.connection.Preload(clause.Associations).Where("ref_number = ?", ref_number).First(&Ticket)
	return Ticket
}

func (db *TicketDatabase) GetTicketBySlotId(slotId int) entity.Ticket {
	var Ticket entity.Ticket
	db.connection.Preload(clause.Associations).Where("slot_id = ?", slotId).First(&Ticket)
	return Ticket
}
