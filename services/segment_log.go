package services

import (
	"avito-task/entity"
	"fmt"

	"gorm.io/gorm"
)

var SegmentLogs SegmentLogService

type SegmentLogService struct {
	DB *gorm.DB
}

func CreateSegmentLogService(db *gorm.DB) SegmentLogService {
	return SegmentLogService{db}
}

func (sl *SegmentLogService) AddOne(user_id uint, segment_name, operation string) (*entity.SegmentLog, error) {
	elem := &entity.SegmentLog{UserID: user_id, SegmentName: segment_name, Operation: operation}
	result := sl.DB.Create(elem)
	if result.Error != nil {
		return nil, result.Error
	}
	return elem, nil
}

func (sl *SegmentLogService) GetAll() ([]*entity.SegmentLog, error) {
	var data []*entity.SegmentLog
	result := sl.DB.Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (sl *SegmentLogService) DeleteAll() error {
	result := sl.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.SegmentLog{})
	return result.Error
}

func (sl *SegmentLogService) GetPeriod(user_id uint, date1, date2 string) ([]*entity.SegmentLog, error) {
	var res []*entity.SegmentLog
	d1, err := GetTimeDate(date1)
	if err != nil {
		return res, err
	}
	d2, err := GetTimeDate(date2)
	if err != nil {
		return res, err
	}
	result := sl.DB.Model(&entity.SegmentLog{}).Where("user_id = ? AND created_at BETWEEN (?) AND (?)", user_id, d1, d2).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}

func (sl *SegmentLogService) GenerateCSV(user_id uint, date1, date2 string) (string, error) {
	logs, err := sl.GetPeriod(user_id, date1, date2)
	if err != nil {
		return "", err
	}
	name := "User" + fmt.Sprint(user_id) + "_" + date1 + "_" + date2
	err = CreateCSVSegmentLogs(name, logs)
	if err != nil {
		return "", err
	}
	return name + ".csv", nil
}
