package business

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/dao"
	"github.com/storyicon/golang-proxy/model"
	"github.com/storyicon/golang-proxy/toolkit"
)

type Publisher struct {
	Sources   *model.Sources
	Database  *gorm.DB
	Scheduler *cron.Cron
}

type Proxy struct {
	IP     string
	Port   string
	Scheme string
}

func NewPublisher(s *model.Sources, db *gorm.DB) *Publisher {
	return &Publisher{
		Sources:   s,
		Database:  db,
		Scheduler: cron.New(),
	}
}

func getProxyString(p *Proxy) string {
	var s string
	if strings.Count(p.IP, ".") != 3 {
		return ""
	}
	if strings.Index(p.IP, "://") <= 0 {
		if p.Scheme == "" {
			p.Scheme = "http"
		}
		s += p.Scheme + "://"
	} else {
		p.Scheme = ""
	}
	s += p.IP
	if p.Port != "" && !strings.HasSuffix(p.IP, ":"+p.Port) {
		s += ":" + p.Port
	}
	return s
}

func fieldFormat(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (p *Publisher) Start() {
	var s *model.Source
	length := len(*p.Sources)
	log.Infof("[P]Totally %d source was found", length)
	for i := 0; i < length; i++ {
		s = &(*p.Sources)[i]
		(func(p *Publisher, s *model.Source) {
			go p.newSourceCrawler(s)
			log.Infof("[P]Periodical Task %s Was Assigned %s", s.Name, s.Category.Interval)
			p.Scheduler.AddFunc(s.Category.Interval, func() {
				log.Infof("[P]Periodical Task @%s is Running!", s.Name)
				p.newSourceCrawler(s)
			})
		})(p, s)
	}
	p.Scheduler.Start()
}

func (p *Publisher) newSourceCrawler(s *model.Source) {
	c := colly.NewCollector(
		colly.UserAgent(UserAgent),
		colly.Async(true),
	)
	c.SetRequestTimeout(RequestTimeout * time.Second)
	c.Limit(&colly.LimitRule{
		Parallelism: s.Category.ParallelNumber,
	})
	c.OnRequest(func(r *colly.Request) {
		if s.Debug {
			log.Infof("[P]Start visit: %s", r.URL)
		}
	})
	c.OnHTML(s.Selector.Iterator, func(e *colly.HTMLElement) {
		item := e.DOM
		if s.Selector.Filter == "" || item.Find(s.Selector.Filter).Length() > 0 {
			proxy := &Proxy{
				IP:     fieldFormat(item.Find(s.Selector.IP).Text()),
				Port:   fieldFormat(item.Find(s.Selector.Port).Text()),
				Scheme: fieldFormat(item.Find(s.Selector.Scheme).Text()),
			}
			if r := getProxyString(proxy); r != "" {
				dao.SaveCrawlProxy(r)
				if s.Debug {
					log.Infof("[P][%s] Dig a proxy %s", s.Name, r)
				}
			}

		}
	})
	c.OnScraped(func(r *colly.Response) {
		if s.Debug {
			log.Infof("[P]Finished visit: %s", r.Request.URL)
		}
	})

	c.Visit(s.Page.Entry)

	for i := s.Page.From; i <= s.Page.To; i++ {
		toolkit.SleepRandomRangeTime(s.Category.DelayRange)
		c.Visit(TemplateRender(s.Page.Template, "page", i))
	}
}

func TemplateRender(template string, key string, value interface{}) string {
	return strings.Replace(template, "{"+key+"}", fmt.Sprintf("%v", value), -1)
}
