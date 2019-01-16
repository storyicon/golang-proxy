package business

import (
	"math"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"golang-proxy/dao"
	"golang-proxy/model"
)

var (
	// AssessorStackLength stores the number of proxy to be evaluated in the current memory.
	AssessorStackLength = 0
)

// StartAssessor is used to start the evaluation procedure.
func StartAssessor() {
	scheduler := cron.New()
	scheduler.AddFunc("@every 3s", func() {
		if AssessorStackLength < AssessorStackCapacity {
			proxies := dao.GetProxy(AssessorInterval, AssessorPerExtract)
			AssessorStackLength += len(proxies)
			for _, proxy := range proxies {
				go (func(proxy *model.Proxy) {
					Assess(proxy)
				})(proxy)
			}
		}
	})
	log.Infoln("Start Assessor")
	scheduler.Start()
}

// Assess is used to evaluate an proxy.
func Assess(proxy *model.Proxy) {
	schemeTest := HTTP
	timestamp := time.Now().UnixNano() / 1e6
	switch proxy.SchemeType {
	case typeHTTP:
		schemeTest = HTTP
	case typeHTTPS:
		schemeTest = HTTPS
	case typeBOTH:
		if timestamp%2 == 0 {
			schemeTest = HTTPS
		}
	default:
		log.Errorf("Unknown proxy scheme type: %d", proxy.SchemeType)
	}

	testOK := HTTPBinTester(proxy.IP, proxy.Port, schemeTest)
	AssessorStackLength--
	timeCost := float64(time.Now().UnixNano()/1e6-timestamp) / 1e3
	feedBack(proxy, testOK, timeCost)
}

func feedBack(proxy *model.Proxy, isOK bool, timeCost float64) {
	proxy.AssessTimes++
	assessTimes := float64(proxy.AssessTimes)
	proxy.AvgResponseTime = (proxy.AvgResponseTime*(assessTimes-1.0) + timeCost) / assessTimes
	if isOK {
		proxy.ContinuousFailedTimes = 0
		proxy.SuccessTimes++
		log.Infof("[Assessor]Proxy %s assess pass(%gms)", proxy.Content, timeCost)
	} else {
		proxy.ContinuousFailedTimes++
		log.Warnf("[Assessor]Proxy %s Assess Failed", proxy.Content)
	}
	proxy.UpdateTime = time.Now().Unix()
	proxy.Score = GetScore(proxy)
	UpdateProxy(proxy)
}

// UpdateProxy is used to update the evaluation information to the database.
func UpdateProxy(proxy *model.Proxy) {
	session := dao.GetDatabase()
	successRate := float64(proxy.SuccessTimes) / float64(proxy.AssessTimes)
	if successRate < AssessorAllowSuccessRateMin {
		log.Warnf("[Assessor]Proxy %s Deleted: score too low", proxy.Content)
		session.Delete(proxy)
	} else {
		session.Save(proxy)
	}
}

// GetScore uses association algorithm to evaluate the score of a Proxy.
// Set 4 impact factors, namely AssessTimes, SuccessTimes, Speed, Mutation
// Continuously increasing Mutation value will lead to a sharp drop in Score
// Formula affected by SuccessRate and AssessTimes at the same time.
// Formulas can be derived by yourself
func GetScore(p *model.Proxy) float64 {
	times := float64(p.AssessTimes)
	success := float64(p.SuccessTimes)
	speed := math.Sqrt(float64(RequestTimeout)) / p.AvgResponseTime
	mutation := 1 / math.Pow(float64(p.ContinuousFailedTimes)+1, 2.0)
	return success * speed * mutation / math.Sqrt(times)
}
