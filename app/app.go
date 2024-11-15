package app

import (
	"tushuguanli/app/model"
	"tushuguanli/app/router"
)

func Start() {
	model.NewMysql()
	//model.NewRdb()
	//model.NewMongoDB()
	defer func() {
		model.Close()
	}()

	//schedule.Start()
	//tools.NewLogger()
	router.New()
}
