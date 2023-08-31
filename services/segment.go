package services

import (
	"avito-task/entity"
	"fmt"

	"gorm.io/gorm"
)

var Segments SegmentService

type SegmentService struct {
	DB *gorm.DB
}

func CreateSegmentService(db *gorm.DB) SegmentService {
	return SegmentService{db}
}

func (s SegmentService) GetAll() ([]*entity.Segment, error) {
	var segments []*entity.Segment
	results := s.DB.Preload("Users").Find(&segments)
	if results.Error != nil {
		return segments, results.Error
	}
	return segments, nil
}

func (s SegmentService) GetOne(name string) (*entity.Segment, error) {
	var segment *entity.Segment
	result := s.DB.Preload("Users").First(&segment, "name = ?", name)
	if result.Error != nil {
		return segment, result.Error
	}
	return segment, nil
}

func (s SegmentService) GetSegmentsForUser(id uint) ([]*entity.Segment, error) {
	var user *entity.User
	result := s.DB.Preload("Segments").Find(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Segments, nil
}

func (s SegmentService) AddOne(name string) (*entity.Segment, error) {
	segment := &entity.Segment{Name: name}
	result := s.DB.Create(segment)
	if result.Error != nil {
		return nil, result.Error
	}
	return segment, nil
}

func (s SegmentService) AddOneWithPercent(name string, percent int) (*entity.Segment, error) {
	segment := &entity.Segment{Name: name}
	var users []*entity.User
	result := s.DB.Create(segment)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	usersToAdd := pickPercentOfUsers(percent, users)
	for _, user := range usersToAdd {
		s.AddSegmentsToUser(user.ID, []string{segment.Name})
	}
	segment.Users = usersToAdd
	return segment, nil
}

func (s SegmentService) UpdateOne(name, newName string) (*entity.Segment, error) {
	segment, err := s.GetOne(name)
	if err != nil {
		return nil, err
	}
	segment.Name = newName
	result := s.DB.Model(&segment).Where("name = ?", name).Update("name", newName)
	if result.Error != nil {
		return nil, result.Error
	}
	return segment, nil
}

func (s SegmentService) DeleteOne(name string) error {
	result := s.DB.Delete(&entity.Segment{}, "name = ?", name)
	return result.Error
}

func (s SegmentService) DeleteAll() error {
	result := s.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.Segment{})
	return result.Error
}

func (s SegmentService) AddSegmentsToUser(id uint, strSegments []string) error {
	var segments []*entity.Segment
	result := s.DB.Model(&entity.Segment{}).Where("name IN (?)", strSegments).Find(&segments)
	if result.Error != nil {
		return result.Error
	}
	var user *entity.User
	result = s.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	user.Segments = append(user.Segments, segments...)
	result = s.DB.Save(user)
	s.DB.Save(user.Segments)
	sl := CreateSegmentLogService(s.DB)
	if result.Error == nil {
		for _, seg := range segments {
			sl.AddOne(id, seg.Name, "added")
		}
	}
	return result.Error
}

func (s SegmentService) RemoveSegmentsFromUser(id uint, strSegments []string) error {
	var user *entity.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	segmentHash := make(map[string]bool)
	for _, seg := range strSegments {
		segmentHash[seg] = true
	}
	segments, err := s.GetSegmentsForUser(id)
	if err != nil {
		return err
	}
	sl := CreateSegmentLogService(s.DB)
	for _, seg := range segments {
		if segmentHash[seg.Name] {
			s.DB.Model(user).Association("Segments").Delete(seg)
			sl.AddOne(id, seg.Name, "deleted")
		}
	}
	return nil
}

func (s SegmentService) RemoveUserFromSegment(id uint, segmentName string) error {
	segment, err := s.GetOne(segmentName)
	if err != nil {
		return err
	}
	var user *entity.User
	result := s.DB.Preload("Segments").Find(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return result.Error
	}
	if len(user.Segments) == 0 {
		return nil
	}
	err = s.DB.Model(user).Association("Segments").Delete(segment)
	sl := CreateSegmentLogService(s.DB)
	if err == nil {
		sl.AddOne(id, segmentName, "deleted")
	}
	return err
}
