package usecases

import (
	"net/http/httptest"
	"testing"
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	repositories "github.com/coroo/go-starter/app/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// dummy data
var dummyAvailableSlot = []entity.AvailableSlot{
	{
		ID:        1,
		LOCATION:  "Jakarta",
		EVENT_ID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, {
		ID:        2,
		LOCATION:  "Jakarta",
		EVENT_ID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

type repoMockAvailableSlot struct {
	mock.Mock
}

func (r *repoMockAvailableSlot) SaveAvailableSlot(AvailableSlot entity.AvailableSlot) (int, error) {
	return 0, nil
}

func (r *repoMockAvailableSlot) UpdateAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	args := r.Called(AvailableSlot)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (r *repoMockAvailableSlot) DeleteAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	args := r.Called(AvailableSlot)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (r *repoMockAvailableSlot) GetAllAvailableSlots() []entity.AvailableSlot {
	return dummyAvailableSlot
}

func (r *repoMockAvailableSlot) GetAvailableSlot(id string) []entity.AvailableSlot {
	return dummyAvailableSlot
}

func (r *repoMockAvailableSlot) GetAvailableSlotByUuid(uuid string) entity.AvailableSlot {
	return dummyAvailableSlot[0]
}

func (r *repoMockAvailableSlot) GetActiveAvailableSlotByCode(code string) entity.AvailableSlot {
	return dummyAvailableSlot[0]
}

func (r *repoMockAvailableSlot) GetAvailableSlotByCode(code string) entity.AvailableSlot {
	return dummyAvailableSlot[0]
}

func (r *repoMockAvailableSlot) CloseDB() {
}

type AvailableSlotUsecaseTestSuite struct {
	suite.Suite
	repositoryTest repositories.AvailableSlotRepository
}

func (suite *AvailableSlotUsecaseTestSuite) SetupTest() {
	suite.repositoryTest = new(repoMockAvailableSlot)
}

func (suite *AvailableSlotUsecaseTestSuite) TestBuildAvailableSlotService() {
	resultTest := NewAvailableSlotService(suite.repositoryTest)
	var dummyImpl *AvailableSlotService
	assert.NotNil(suite.T(), resultTest)
	assert.Implements(suite.T(), dummyImpl, resultTest)
	// assert.NotNil(suite.T(), resultTest.(*AvailableSlotService).repositories)
}

func (suite *AvailableSlotUsecaseTestSuite) TestSaveAvailableSlotUsecase() {
	suite.repositoryTest.(*repoMockAvailableSlot).On("SaveAvailableSlot", dummyAvailableSlot[0]).Return(nil)
	useCaseTest := NewAvailableSlotService(suite.repositoryTest)
	// dummyAvailableSlot[0].Password = "Change Password"
	data, _ := useCaseTest.SaveAvailableSlot(dummyAvailableSlot[0])
	assert.NotNil(suite.T(), data)
}

func (suite *AvailableSlotUsecaseTestSuite) TestUpdateAvailableSlotUsecase() {
	suite.repositoryTest.(*repoMockAvailableSlot).On("UpdateAvailableSlot", dummyAvailableSlot[0]).Return(nil)
	useCaseTest := NewAvailableSlotService(suite.repositoryTest)
	err := useCaseTest.UpdateAvailableSlot(dummyAvailableSlot[0])
	assert.Nil(suite.T(), err)
}

func (suite *AvailableSlotUsecaseTestSuite) TestDeleteAvailableSlotUsecase() {
	suite.repositoryTest.(*repoMockAvailableSlot).On("DeleteAvailableSlot", dummyAvailableSlot[0]).Return(nil)
	useCaseTest := NewAvailableSlotService(suite.repositoryTest)
	err := useCaseTest.DeleteAvailableSlot(dummyAvailableSlot[0])
	assert.Nil(suite.T(), err)
}

func (suite *AvailableSlotUsecaseTestSuite) TestGetAllAvailableSlots() {
	suite.repositoryTest.(*repoMockAvailableSlot).On("GetAllAvailableSlots", dummyAvailableSlot).Return(dummyAvailableSlot)
	useCaseTest := NewAvailableSlotService(suite.repositoryTest)
	dummyAvailableSlot := useCaseTest.GetAllAvailableSlots()
	assert.Equal(suite.T(), dummyAvailableSlot, dummyAvailableSlot)
}

func (suite *AvailableSlotUsecaseTestSuite) TestGetAvailableSlot() {
	suite.repositoryTest.(*repoMockAvailableSlot).On("GetAvailableSlot", dummyAvailableSlot[0].ID).Return(dummyAvailableSlot[0], nil)
	useCaseTest := NewAvailableSlotService(suite.repositoryTest)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	dummyAvailableSlot := useCaseTest.GetAvailableSlot(c.Param("id"))
	assert.NotNil(suite.T(), dummyAvailableSlot[0])
	assert.Equal(suite.T(), dummyAvailableSlot[0], dummyAvailableSlot[0])
}

func TestAvailableSlotUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AvailableSlotUsecaseTestSuite))
}
