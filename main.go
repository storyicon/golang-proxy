package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	"storyicon.visualstudio.com/golang-proxy/business"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "all", "all/consumer/producer/assessor/service, default is all")
	flag.Parse()

	log.Infof("Operating Mode: %s, will start running after 3 seconds", mode)

	time.Sleep(3 * time.Second)

	switch mode {
	case "all":
		go business.StartConsumer()
		go business.StartProducer()
		go business.StartAssessor()
		go business.StartService()
	case "consumer":
		business.StartConsumer()
	case "producer":
		business.StartProducer()
	case "assessor":
		business.StartAssessor()
	case "service":
		business.StartService()
	default:
		log.Panicf("Unknown mode: %s", mode)
	}
	select {}
}
