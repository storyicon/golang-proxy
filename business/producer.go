package business

import (
	"github.com/arkadybag/golang-proxy/dao"
	"github.com/arkadybag/golang-proxy/model"
	"github.com/arkadybag/golang-proxy/std"
	"math/rand"
	"time"

	"github.com/gocolly/colly"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// StartProducer used to start the producer
// Producer used to cralw source and store the crawled proxy in database
func StartProducer() {
	scheduler := cron.New()
	sources := dao.GetSources()

	log.Infof("[Producer]Totally %d source was found", len(sources))

	std.Dump(sources)
	for _, source := range sources {
		(func(scheduler *cron.Cron, source *model.Source) {
			go newSourceCrawler(source)
			name := source.Name
			interval := source.Category.Interval
			log.Infof("[Producer]Periodical Task %s Was Assigned %s", name, interval)
			scheduler.AddFunc(interval, func() {
				log.Infof("[Producer]Periodical Task @%s is Running!", name)
				newSourceCrawler(source)
			})
		})(scheduler, source)
	}

	scheduler.Start()
}

func newSourceCrawler(source *model.Source) {
	var (
		name           = source.Name
		debug          = source.Debug
		parallelNumber = source.Category.ParallelNumber
		iterator       = source.Selector.Iterator
		IPSelector     = source.Selector.IP
		portSelector   = source.Selector.Port
		startURL       = source.Page.Entry
		pageFrom       = source.Page.From
		pageTo         = source.Page.To
		delayRange     = source.Category.DelayRange
		template       = source.Page.Template
	)
	c := colly.NewCollector(
		colly.UserAgent(UserAgent),
		colly.Async(true),
	)
	c.SetRequestTimeout(RequestTimeout * time.Second)
	c.Limit(&colly.LimitRule{
		Parallelism: parallelNumber,
	})
	c.OnError(func(_ *colly.Response, err error) {
		if debug {
			log.Warnf("[Producer][%s]Visit error: %s", name, err)
		}
	})
	c.OnRequest(func(request *colly.Request) {
		if debug {
			log.Infof("[Producer][%s]Start visit: %s", name, request.URL)
		}
	})
	c.OnHTML(iterator, func(element *colly.HTMLElement) {
		item := element.DOM
		proxy := NewProxy(
			item.Find(IPSelector).Text(),
			item.Find(portSelector).Text(),
		)
		if proxy == nil {
			return
		}
		if content := proxy.String(); content != "" {
			dao.SaveCrudeProxy(&model.CrudeProxy{
				IP:      proxy.IP,
				Port:    proxy.Port,
				Content: content,
			})
			if debug {
				log.Infof("[Producer][%s]Proxy %s was mined", name, content)
			}
		}
	})
	c.OnScraped(func(response *colly.Response) {
		if debug {
			log.Infof("[Producer][%s]Finish visit: %s", name, response.Request.URL)
		}
	})

	c.Visit(startURL)

	for i := pageFrom; i < pageTo; i++ {
		sleep(delayRange)
		c.Visit(std.TemplateRender(template, "page", i))
	}

}

func sleep(delayRange []int) {
	delay := 1
	switch count := len(delayRange); count {
	case 0:
		break
	case 1:
		delay = delayRange[0]
	case 2:
		delay = delayRange[0] + rand.Intn(delayRange[1]-delayRange[0])
	default:
		delay = delayRange[rand.Intn(count)]
	}
	time.Sleep(time.Duration(delay) * time.Second)
}
