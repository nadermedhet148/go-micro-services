package repositories

import (
	"testing"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TicketRepositoryTestSuite struct {
	suite.Suite
	ctx *gin.Context
	db  *gorm.DB
}

func (suite *TicketRepositoryTestSuite) SetupTest() {
	suite.db, _ = config.ConnectDB()
}

func (suite *TicketRepositoryTestSuite) TestA_BuildNewTicketRepository() {
	repoTest := NewTicketRepository()
	var dummyImpl *TicketRepository
	assert.NotNil(suite.T(), repoTest)
	assert.Implements(suite.T(), dummyImpl, repoTest)
}

func (suite *TicketRepositoryTestSuite) TestB_CreateTicket() {
	repoTest := NewTicketRepository()
	dummyTicket := entity.Ticket{
		ID: 1,
	}
	_, err := repoTest.SaveTicket(dummyTicket)
	assert.Nil(suite.T(), err)
}

func (suite *TicketRepositoryTestSuite) TestC_UpdateTicket() {
	repoTest := NewTicketRepository()
	dummyTicket := entity.Ticket{
		ID:         1,
		REF_NUMBER: "123",
		SLOT_ID:    1,
	}
	TicketDummy := repoTest.UpdateTicket(dummyTicket)
	assert.Nil(suite.T(), TicketDummy)
}

func (suite *TicketRepositoryTestSuite) TestE_GetAllTickets() {
	repoTest := NewTicketRepository()
	TicketDummy := repoTest.GetAllTickets()
	assert.NotNil(suite.T(), TicketDummy)
}

func (suite *TicketRepositoryTestSuite) TestF_GetTicket() {
	repoTest := NewTicketRepository()
	TicketDummy := repoTest.GetTicket("1")
	assert.NotNil(suite.T(), TicketDummy)
}

func (suite *TicketRepositoryTestSuite) TestH_RemoveTicket() {
	repoTest := NewTicketRepository()
	dummyTicket := entity.Ticket{
		ID: 1,
	}
	TicketDummy := repoTest.DeleteTicket(dummyTicket)
	assert.Nil(suite.T(), TicketDummy)
}

func TestTicketRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TicketRepositoryTestSuite))
}
