package services

import (
	"avito-task/entity"
	"math/rand"
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
