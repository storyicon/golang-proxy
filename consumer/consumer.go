package main

import (
	"golang-proxy/library"
	"log"
	"net/http"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/parnurzeal/gorequest"
)

type Consumer struct {
	Config      *library.Config
	NSQConsumer *nsq.Consumer
}

type ProxyHandler struct {
	MySQL *library.MySQL
}

func (h *ProxyHandler) HandleMessage(msg *nsq.Message) error {
	s := string(string(msg.Body))
	p := &library.Proxy{}
	library.JSONDecode(s, p)
	go h.ProxyPreAssess(p)
	return nil
}

func (h *ProxyHandler) ProxyPreAssess(p *library.Proxy) {
	var r library.HTTPbinIP
	u := library.ProxyStringify(p)
	request := gorequest.New().Proxy(u).Timeout(library.ProxyAssessTimeOut * time.Second)
	resp, _, errs := request.Get("http://httpbin.org/ip").
		Retry(
			library.ProxyPreAssessTimes,
			library.ProxyAssessTimeOut*time.Second,
			http.StatusBadRequest,
			http.StatusInternalServerError).
		EndStruct(&r)
	if len(errs) == 0 && resp.StatusCode == 200 && r.Origin == p.IP {
		log.Println("Successfully passed the first trial:", u)
		h.InsertProxy(p)
		return
	} else {
		log.Printf("Failed trail: %s", u)
	}
}

func (h *ProxyHandler) InsertProxy(p *library.Proxy) {
	h.MySQL.Connection.Create(&library.ValidProxy{
		// ID: "",
		Content:               library.ProxyStringify(p),
		AssessTimes:           0,
		SuccessTimes:          0,
		AvgResponseTime:       0,
		ContinuousFailedTimes: 0,
		LastAssessTime:        0,
		Score:                 library.ProxyInitScore,
	})
}

func NewConsumer() *Consumer {
	c := library.GetConfig()
	d := library.GetMysqlDsn(c)

	return &Consumer{
		Config: c,
		NSQConsumer: library.NewNSQConsumer(
			c.Nsq.ConsumerAddr,
			library.NSQTopic,
			"Consumer",
			&ProxyHandler{
				MySQL: library.NewMySQL(d),
			},
		),
	}
}

func main() {
	NewConsumer()
	select {}
}
