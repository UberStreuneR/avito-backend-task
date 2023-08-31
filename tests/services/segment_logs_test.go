package tests

import (
	"avito-task/initializers"
	"avito-task/services"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SegmentLogServiceTestSuite struct {
	suite.Suite
	u  services.UserService
	s  services.SegmentService
	sl services.SegmentLogService
}

func (suite *SegmentLogServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *SegmentLogServiceTestSuite) TearDownTest() {
	suite.s.DeleteAll()
	suite.sl.DeleteAll()
}

func (suite *SegmentLogServiceTestSuite) TestAddOne() {
	res, err := suite.sl.AddOne(1, "TEST_SEGMENT", "delete")
	suite.Equal(err, nil)
	suite.Equal(res.UserID, uint(1))
	suite.Equal(res.SegmentName, "TEST_SEGMENT")
	suite.Equal(res.Operation, "delete")
}

func (suite *SegmentLogServiceTestSuite) TestGetAll() {
	suite.sl.AddOne(1, "TEST_SEGMENT", "delete")
	suite.sl.AddOne(2, "TEST_SEGMENT", "delete")
	res, err := suite.sl.GetAll()
	suite.Equal(err, nil)
	suite.Equal(err, nil)
	suite.Equal(len(res), 2)
}

func (suite *SegmentLogServiceTestSuite) TestAutoLog() {
	u, err := suite.u.AddOne(1500)
	suite.Equal(err, nil)
	s, _ := suite.s.AddOne("TEST_AUTO_SEGMENT")
	s2, _ := suite.s.AddOne("TEST_AUTO_SEGMENT2")
	suite.s.AddSegmentsToUser(u.ID, []string{s.Name, s2.Name})
	logs, err := suite.sl.GetPeriod(u.ID, "2023-01", "2023-12")
	suite.Equal(err, nil)
	suite.Equal(len(logs), 2)
	suite.Equal(logs[0].UserID, u.ID)
	suite.Equal(logs[0].SegmentName, s.Name)
	suite.Equal(logs[0].Operation, "added")
	suite.Equal(logs[1].UserID, u.ID)
	suite.Equal(logs[1].SegmentName, s2.Name)
	suite.Equal(logs[1].Operation, "added")

	suite.s.RemoveSegmentsFromUser(u.ID, []string{"TEST_AUTO_SEGMENT", "TEST_AUTO_SEGMENT2"})
	logs, err = suite.sl.GetPeriod(u.ID, "2023-01", "2023-12")
	suite.Equal(err, nil)
	suite.Equal(len(logs), 4)
	suite.Equal(logs[2].UserID, u.ID)
	suite.Equal(logs[2].SegmentName, s.Name)
	suite.Equal(logs[2].Operation, "deleted")
	suite.Equal(logs[3].UserID, u.ID)
	suite.Equal(logs[3].SegmentName, s2.Name)
	suite.Equal(logs[3].Operation, "deleted")
}
func TestSegmentLogService(t *testing.T) {
	config, err := initializers.LoadConfig("../..")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	u := services.CreateUserService(db)
	s := services.CreateSegmentService(db)
	sl := services.CreateSegmentLogService(db)
	suite.Run(t, &SegmentLogServiceTestSuite{u: u, s: s, sl: sl})
}
