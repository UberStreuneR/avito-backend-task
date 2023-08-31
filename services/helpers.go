package services

import (
	"avito-task/entity"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func pickPercentOfUsers(percent int, arr []*entity.User) []*entity.User {
	rand.Seed(time.Now().UnixNano())
	count := len(arr) * percent / 100
	if count == 0 {
		count = 1
	}
	taken := make(map[int]bool)
	var res []*entity.User
	if len(arr) == 0 {
		return res
	}
	for count > 0 {
		ind := rand.Intn(len(arr))
		if taken[ind] {
			continue
		}
		res = append(res, arr[ind])
		taken[ind] = true
		count--
	}
	return res
}

func GetTimeDate(date string) (time.Time, error) {
	d := strings.Split(date, "-")
	year, err := strconv.Atoi(d[0])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(d[1])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

func CreateCSVSegmentLogs(name string, logs []*entity.SegmentLog) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	rootPath := exPath[:strings.Index(exPath, "/tmp")]
	file, err := os.Create(rootPath + "/static/" + name + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, sl := range logs {
		writer.Write([]string{fmt.Sprint(sl.UserID), sl.SegmentName, sl.Operation, sl.CreatedAt.String()})
	}
	return nil
}
