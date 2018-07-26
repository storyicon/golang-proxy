package business

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/dao"
	"github.com/storyicon/golang-proxy/model"
	"github.com/storyicon/golang-proxy/toolkit"
)

var (
	ConsuemrStackLength = 0
	AssessorStackLength = 0
)

type Consumer struct {
	Database  *gorm.DB
	Scheduler *cron.Cron
}

func NewConsumer(db *gorm.DB) *Consumer {
	return &Consumer{
		Database:  db,
		Scheduler: cron.New(),
	}
}

func (c *Consumer) PreAssess(proxy string) {
	var r model.HTTPBinIP
	req := gorequest.New().Proxy(proxy).Timeout(RequestTimeout * time.Second)
	res, _, errs := req.Get("http://httpbin.org/ip").
		Retry(
			ConsumerAssessTimes,
			RequestTimeout,
			http.StatusBadRequest,
			http.StatusInternalServerError,
		).
		EndStruct(&r)
	ConsuemrStackLength--
	if len(errs) == 0 && res.StatusCode == 200 {
		if toolkit.GetHostNameByIP(proxy) == r.Origin {
			log.Infof("[C]Proxy Pre Assess Successful: %s", proxy)
			dao.SaveValidProxy(proxy)
			return
		}
		log.Warnf(`[C]Proxy %s Pre Assess Failed: Not Highly Anonymous`, proxy)
		return
	}
	log.Warnf(`[C]Proxy %s Pre Assess Failed: Connection Timeout or Refused`, proxy)
}

func (c *Consumer) Start() {
	c.Scheduler.AddFunc("@every 1s", func() {
		if ConsuemrStackLength < ConsumerStackCapacity {
			proxy := dao.PopCrawlProxy(0, ConsumerPerExtract)
			ConsuemrStackLength += ConsumerPerExtract
			for _, v := range *proxy {
				go (func(proxy string) {
					c.PreAssess(proxy)
				})(v.Content)
			}
		}
	})
	c.Scheduler.Start()
}
