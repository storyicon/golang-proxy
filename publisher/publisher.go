package main

import (
	"log"
	"strconv"
	"strings"

	"golang-proxy/library"

	"github.com/gocolly/colly"
	nsq "github.com/nsqio/go-nsq"
	"github.com/robfig/cron"
)

type Publisher struct {
	Config      *library.Config
	Source      *[]library.Source
	NSQProducer *nsq.Producer
	Scheduler   *cron.Cron
}

func (p *Publisher) NewCrawler(s *library.Source) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: s.Category.Parallelnumber,
	})
	c.OnRequest(func(r *colly.Request) {
		if s.Debug {
			log.Println("Start visit", r.URL)
		}
	})
	c.OnHTML(s.Selector.Iterator, func(e *colly.HTMLElement) {
		item := e.DOM
		if s.Selector.Filter == "" || item.Find(s.Selector.Filter).Length() > 0 {
			ip := item.Find(s.Selector.IP).Text()
			if strings.Count(ip, ".") == 3 {
				proxy := library.Proxy{
					IP:        ip,
					Port:      item.Find(s.Selector.Port).Text(),
					Protocal:  strings.ToLower(item.Find(s.Selector.Protocal).Text()),
					Proxytype: item.Find(s.Selector.Proxytype).Length(),
				}
				if proxy.Protocal == "" {
					proxy.Protocal = "http"
				}
				ps, _ := library.JSONEncode(proxy)
				p.NSQProducer.Publish(library.NSQTopic, ps)
				if s.Debug {
					log.Printf("[%s] Get a Proxy %s", s.Name, string(ps))
				}
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if s.Debug {
			log.Println("Finished visit", r.Request.URL)
		}
	})

	c.Visit(s.Page.Entry)

	for i := s.Page.From; i <= s.Page.To; i++ {
		library.SleepRangeTime(s.Category.DelayRange)
		url := strings.Replace(s.Page.Template, "{page}", strconv.Itoa(i), -1)
		c.Visit(url)
	}
}

func NewPublisher() *Publisher {
	c := library.GetConfig()
	return &Publisher{
		Config:      c,
		Source:      library.GetSource(),
		NSQProducer: library.NewNSQProduer(c.Nsq.ProducerAddr),
		Scheduler:   cron.New(),
	}
}

func (p *Publisher) Start() {
	var s *library.Source
	length := len(*p.Source)
	log.Printf("total %d source was found", length)
	for i := 0; i < length; i++ {
		s = &(*p.Source)[i]
		(func(pub *Publisher, sour *library.Source) {
			go pub.NewCrawler(sour)
			log.Printf("Periodical Task %s Was Assigned %s", sour.Name, sour.Category.Interval)
			pub.Scheduler.AddFunc(sour.Category.Interval, func() {
				pub.NewCrawler(sour)
			})
		})(p, s)
	}
	p.Scheduler.Start()
	select {}
}

func main() {
	publisher := NewPublisher()
	publisher.Start()
}
