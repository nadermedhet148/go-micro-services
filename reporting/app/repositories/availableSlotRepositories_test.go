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

type AvailableSlotRepositoryTestSuite struct {
	suite.Suite
	ctx *gin.Context
	db  *gorm.DB
}

func (suite *AvailableSlotRepositoryTestSuite) SetupTest() {
	suite.db, _ = config.ConnectDB()
}

func (suite *AvailableSlotRepositoryTestSuite) TestA_BuildNewAvailableSlotRepository() {
	repoTest := NewAvailableSlotRepository()
	var dummyImpl *AvailableSlotRepository
	assert.NotNil(suite.T(), repoTest)
	assert.Implements(suite.T(), dummyImpl, repoTest)
}

func (suite *AvailableSlotRepositoryTestSuite) TestB_CreateAvailableSlot() {
	repoTest := NewAvailableSlotRepository()
	dummyAvailableSlot := entity.AvailableSlot{
		ID:       1,
		LOCATION: "Jakarta",
		EVENT_ID: 1,
	}
	_, err := repoTest.SaveAvailableSlot(dummyAvailableSlot)
	assert.Nil(suite.T(), err)
}

func (suite *AvailableSlotRepositoryTestSuite) TestC_UpdateAvailableSlot() {
	repoTest := NewAvailableSlotRepository()
	dummyAvailableSlot := entity.AvailableSlot{
		ID:       1,
		LOCATION: "Jakarta",
		EVENT_ID: 1,
	}
	AvailableSlotDummy := repoTest.UpdateAvailableSlot(dummyAvailableSlot)
	assert.Nil(suite.T(), AvailableSlotDummy)
}

func (suite *AvailableSlotRepositoryTestSuite) TestE_GetAllAvailableSlots() {
	repoTest := NewAvailableSlotRepository()
	AvailableSlotDummy := repoTest.GetAllAvailableSlots()
	assert.NotNil(suite.T(), AvailableSlotDummy)
}

func (suite *AvailableSlotRepositoryTestSuite) TestF_GetAvailableSlot() {
	repoTest := NewAvailableSlotRepository()
	AvailableSlotDummy := repoTest.GetAvailableSlot("1")
	assert.NotNil(suite.T(), AvailableSlotDummy)
}

func (suite *AvailableSlotRepositoryTestSuite) TestH_RemoveAvailableSlot() {
	repoTest := NewAvailableSlotRepository()
	dummyAvailableSlot := entity.AvailableSlot{
		ID: 1,
	}
	AvailableSlotDummy := repoTest.DeleteAvailableSlot(dummyAvailableSlot)
	assert.Nil(suite.T(), AvailableSlotDummy)
}

func TestAvailableSlotRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableSlotRepositoryTestSuite))
}
