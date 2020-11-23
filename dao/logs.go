package dao

import (
	"errors"
	"x-server/core/db"
	"x-server/model"
)

func GetLogsFirst(where string, order []string) (model.Logs, error) {
	logs := model.Logs{}

	dbs, err := db.Get("default", false)
	if err != true {
		return logs, errors.New("db not match")
	}
	dbs.Table("logs").Where(where)
	//dbs.LogMode(true)
	if len(order) > 0 {
		for v := range order {
			dbs.Order(v)
		}
	}
	dbs.Limit(1).Find(&logs)
	return logs, nil
}
