package business

import (
	"math"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/dao"
	"github.com/storyicon/golang-proxy/model"
	"github.com/storyicon/golang-proxy/toolkit"
)

type Assessor struct {
	Database  *gorm.DB
	Scheduler *cron.Cron
}

func NewAssessor(db *gorm.DB) *Assessor {
	return &Assessor{
		Database:  db,
		Scheduler: cron.New(),
	}
}

func (s *Assessor) Assess(p *model.ValidProxy) {
	r := &model.HTTPBinIP{}
	req := gorequest.New().Proxy(p.Content).Timeout(RequestTimeout * time.Second)
	timeStart := time.Now().UnixNano() / 1e6
	res, _, errs := req.Get("http://httpbin.org/ip").EndStruct(r)
	AssessorStackLength--
	if len(errs) == 0 && res.StatusCode == 200 {
		if toolkit.GetHostNameByIP(p.Content) == r.Origin {
			timeCost := time.Now().UnixNano()/1e6 - timeStart
			s.FeedBack(p, 1, float64(timeCost)/1e3)
			log.Infof("[A]Proxy %s Assess Pass: (%dms)", p.Content, timeCost)
			return
		}
	}
	log.Warnf("[A]Proxy %s Assess Failed", p.Content)
	s.FeedBack(p, 0, float64(RequestTimeout)*1.5)
}

func (s *Assessor) FeedBack(p *model.ValidProxy, isSucc int, responseTime float64) {
	p.AssessTimes++
	times := float64(p.AssessTimes)
	p.AvgResponseTime = (p.AvgResponseTime*(times-1.0) + responseTime) / times
	if isSucc == 1 {
		p.ContinuousFailedTimes = 0
	} else {
		p.ContinuousFailedTimes++
	}
	p.SuccessTimes += isSucc
	p.UpdateTime = time.Now().Unix()
	p.Score = GetScore(p)
	s.UpdateValidProxy(p)
}

/**
Set 4 impact factors, namely AssessTimes, SuccessTimes, Speed, Mutation
Continuously increasing Mutation value will lead to a sharp drop in Score
Formula affected by SuccessRate and AssessTimes at the same time.
Formulas can be derived by yourself
*/
func GetScore(p *model.ValidProxy) float64 {
	times := float64(p.AssessTimes)
	success := float64(p.SuccessTimes)
	speed := math.Sqrt(float64(RequestTimeout)) / p.AvgResponseTime
	mutation := 1 / math.Pow(float64(p.ContinuousFailedTimes)+1, 2.0)
	return success * speed * mutation / math.Sqrt(times)
}

func (s *Assessor) UpdateValidProxy(p *model.ValidProxy) {
	db := dao.GetSQLite()
	succRate := float64(p.SuccessTimes) / float64(p.AssessTimes)
	if succRate < AllowAssessSuccessRateMin {
		log.Warnf("[A]Proxy %s Deleted: score too low", p.Content)
		db.Delete(p)
	} else {
		db.Save(p)
	}
}

func (s *Assessor) Start() {
	s.Scheduler.AddFunc("@every 3s", func() {
		if AssessorStackLength < AssessorStackCapacity {
			proxy := dao.GetValidProxy(AssessorInterval, AssessorPerExtract)
			for _, v := range *proxy {
				AssessorStackLength++
				go (func(p model.ValidProxy) {
					s.Assess(&p)
				})(v)
			}
		}
	})
	s.Scheduler.Start()
}
