package main

import (
	"github.com/storyicon/golang-proxy/business"
	"github.com/storyicon/golang-proxy/dao"
)

func main() {
	sources := dao.GetSources()
	database := dao.GetSQLite()
	publisher := business.NewPublisher(sources, database)
	consumer := business.NewConsumer(database)
	assessor := business.NewAssessor(database)
	go publisher.Start()
	go consumer.Start()
	go assessor.Start()
	go business.StartService()
	select {}
}
