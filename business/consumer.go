package business

import (
	"net/http"
	"time"

	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"storyicon.visualstudio.com/golang-proxy/dao"
	"storyicon.visualstudio.com/golang-proxy/model"
)

const (
	// HTTP defines an proxy type
	HTTP = "http"
	// HTTPS defines an proxy type
	HTTPS = "https"
)

var (
	// ConsumerStackLength stores the number of proxy to be pre evaluated in the current memory.
	ConsumerStackLength = 0
)

// StartConsumer is used to start the consumer
func StartConsumer() {
	scheduler := cron.New()
	scheduler.AddFunc("@every 1s", func() {
		if ConsumerStackLength < ConsumerStackCapacity {
			log.Infoln("[Comsumer]No proxy in stack, start to extract proxy from database to pre assess")
			proxies := dao.PopCrudeProxy(0, ConsumerPerExtract)
			ConsumerStackLength += len(proxies)
			for _, proxy := range proxies {
				go (func(proxy *model.CrudeProxy) {
					PreAssess(proxy)
				})(proxy)
			}
		}
	})
	log.Infoln("Start Comsumer")
	scheduler.Start()
}

// PreAssess is used to pre assess an proxy.
func PreAssess(proxy *model.CrudeProxy) {
	IP, port := proxy.IP, proxy.Port
	httpOK := HTTPBinTester(IP, port, HTTP)
	httpsOK := HTTPBinTester(IP, port, HTTPS)

	ConsumerStackLength--

	var schemeType int64
	if httpOK && httpsOK {
		schemeType = typeBOTH
	} else if httpsOK {
		schemeType = typeHTTPS
	} else if !httpsOK && !httpOK {
		return
	}
	dao.SaveProxy(&model.Proxy{
		IP:         proxy.IP,
		Port:       proxy.Port,
		SchemeType: schemeType,
		Content:    proxy.Content,
	})
}

// HTTPBinTester is used to use httpbin test agent.
func HTTPBinTester(IP string, port string, schemeTest string) bool {
	switch schemeTest {
	case HTTPS:
	default:
		schemeTest = HTTP
	}
	url, proxy := fmt.Sprintf("%s://httpbin.org/ip", schemeTest),
		fmt.Sprintf("%s://%s:%s", schemeTest, IP, port)
	httpBin := &model.HTTPBinIP{}
	request := gorequest.New().Proxy(proxy).Timeout(RequestTimeout * time.Second)
	response, _, errs := request.Get(url).
		Retry(
			ConsumerRetryTimes,
			RequestTimeout,
			http.StatusBadRequest,
			http.StatusInternalServerError,
		).
		EndStruct(httpBin)
	if len(errs) == 0 && response.StatusCode == 200 {
		if IP == httpBin.Origin {
			log.Infof(`[Consumer][%s]Proxy Pre Assess Successful: %s`, schemeTest, proxy)
			return true
		}
		log.Println(IP, httpBin.Origin)
		log.Warnf(`[Consumer]Proxy %s Pre Assess Failed: Not Highly Anonymous`, proxy)
		return false
	}
	log.Warnf(`[Consumer]Proxy %s Pre Assess Failed: Connection Timeout or Refused`, proxy)
	return false
}
