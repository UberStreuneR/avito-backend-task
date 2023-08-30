package tests

import (
	"avito-task/initializers"
	"avito-task/services"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SegmentServiceTestSuite struct {
	suite.Suite
	u services.UserService
	s services.SegmentService
}

func (suite *SegmentServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *SegmentServiceTestSuite) TearDownTest() {
	suite.s.DeleteAll()
}

func (suite *SegmentServiceTestSuite) TestAddSegment() {
	segment, err := suite.s.AddOne("AVITO_TEST_SEGMENT")
	suite.Equal(err, nil)
	segmentInDB, err := suite.s.GetOne("AVITO_TEST_SEGMENT")
	suite.Equal(err, nil)
	suite.Equal(segment.Name, segmentInDB.Name)
}

func (suite *SegmentServiceTestSuite) TestGetAllSegments() {
	suite.s.AddOne("AVITO_ANOTHER_SEGMENT")
	suite.s.AddOne("AVITO_VOICE_MESSAGES")
	segments, err := suite.s.GetAll()
	suite.Equal(err, nil)
	suite.Equal(len(segments), 2)
	suite.Equal(segments[0].Name, "AVITO_ANOTHER_SEGMENT")
	suite.Equal(segments[1].Name, "AVITO_VOICE_MESSAGES")
}

func (suite *SegmentServiceTestSuite) TestUpdateSegment() {
	suite.s.AddOne("AVITO_ANOTHER_SEGMENT")
	suite.s.UpdateOne("AVITO_ANOTHER_SEGMENT", "AVITO_DISCOUNT_50")
	_, err := suite.s.GetOne("AVITO_ANOTHER_SEGMENT")
	suite.Equal(err, gorm.ErrRecordNotFound)
	updated, err := suite.s.GetOne("AVITO_DISCOUNT_50")
	suite.Equal(err, nil)
	suite.Equal(updated.Name, "AVITO_DISCOUNT_50")
}

func (suite *SegmentServiceTestSuite) TestDeleteSegment() {
	suite.s.AddOne("AVITO_DISCOUNT_50")
	suite.s.DeleteOne("AVITO_DISCOUNT_50")
	_, err := suite.s.GetOne("AVITO_DISCOUNT_50")
	suite.Equal(err, gorm.ErrRecordNotFound)
}

func (suite *SegmentServiceTestSuite) TestAddSegmentToUser() {
	user, _ := suite.u.AddOne(2000)
	suite.s.AddOne("TEST_SEGMENT")
	suite.s.AddOne("TEST_SEGMENT2")
	suite.s.AddSegmentsToUser(user.ID, []string{"TEST_SEGMENT", "TEST_SEGMENT2"})
	segments, err := suite.s.GetSegmentsForUser(user.ID)
	suite.Equal(err, nil)
	suite.Equal(len(segments), 2)
	suite.Equal(segments[0].Name, "TEST_SEGMENT")
	suite.Equal(segments[1].Name, "TEST_SEGMENT2")
}

func (suite *SegmentServiceTestSuite) TestRemoveSegmentFromUser() {
	suite.s.AddOne("TEST_SEGMENT")
	suite.s.AddOne("TEST_SEGMENT2")
	suite.s.AddSegmentsToUser(2000, []string{"TEST_SEGMENT", "TEST_SEGMENT2"})
	suite.s.RemoveUserFromSegment(2000, "TEST_SEGMENT")
	segments, _ := suite.s.GetSegmentsForUser(2000)
	suite.Equal(len(segments), 1)
	suite.Equal(segments[0].Name, "TEST_SEGMENT2")
	segment, err := suite.s.GetOne("TEST_SEGMENT")
	suite.Equal(err, nil)
	suite.Equal(segment.Name, "TEST_SEGMENT")
	suite.Equal(len(segment.Users), 0)
}

func (suite *SegmentServiceTestSuite) TestRemoveSegmentsFromUser() {
	suite.u.AddOne(2001)
	suite.s.AddOne("TEST_SEGMENT")
	suite.s.AddOne("TEST_SEGMENT2")
	suite.s.AddOne("TEST_SEGMENT3")
	suite.s.AddSegmentsToUser(2001, []string{"TEST_SEGMENT", "TEST_SEGMENT2", "TEST_SEGMENT3"})
	err := suite.s.RemoveSegmentsFromUser(2001, []string{"TEST_SEGMENT", "TEST_SEGMENT3"})
	suite.Equal(err, nil)
	segments, _ := suite.s.GetSegmentsForUser(2001)
	suite.Equal(len(segments), 1)
	suite.Equal(segments[0].Name, "TEST_SEGMENT2")
}

func (suite *SegmentServiceTestSuite) TestAddNonExistentSegments() {
	user, _ := suite.u.AddOne(3000)
	err := suite.s.AddSegmentsToUser(user.ID, []string{"FAKE_SEGMENT", "FAKE_SEGMENT2"})
	suite.Equal(err, nil)
	segments, _ := suite.s.GetSegmentsForUser(user.ID)
	suite.Equal(len(segments), 0)
	suite.s.AddOne("TEST_SEGMENT")
	err = suite.s.AddSegmentsToUser(user.ID, []string{"FAKE_SEGMENT", "TEST_SEGMENT"})
	segments, _ = suite.s.GetSegmentsForUser(user.ID)
	suite.Equal(len(segments), 1)
}

func (suite *SegmentServiceTestSuite) TestAddToNonExistentUser() {
	suite.s.AddOne("TEST_SEGMENT2")
	suite.s.AddOne("REAL_SEGMENT")
	err := suite.s.AddSegmentsToUser(1200, []string{"TEST_SEGMENT2", "REAL_SEGMENT"})
	suite.Equal(err, gorm.ErrRecordNotFound)
	segments, _ := suite.s.GetAll()
	suite.Equal(len(segments), 2)
	suite.Equal(len(segments[0].Users), 0)
	suite.Equal(len(segments[1].Users), 0)
}

func TestSegmentService(t *testing.T) {
	config, err := initializers.LoadConfig("../..")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	u := services.CreateUserService(db)
	s := services.CreateSegmentService(db)
	suite.Run(t, &SegmentServiceTestSuite{u: u, s: s})
}
