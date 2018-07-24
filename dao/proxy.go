package dao

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/storyicon/golang-proxy/model"
)

const (
	ProxyInitScore = 1
)

func GetSQLResult(table string, sql string) (interface{}, error) {
	var r interface{}
	db := GetSQLite()
	switch table {
	case "valid_proxy":
		r = &[]model.ValidProxy{}
	case "crawl_proxy":
		r = &[]model.CrawlProxy{}
	default:
		return nil, errors.New("Query unknown table")
	}
	err := db.Raw(sql).Scan(r).Error
	return r, err
}

func GetValidProxy(interval int64, limit int) *[]model.ValidProxy {
	proxy := &[]model.ValidProxy{}
	db := GetSQLite()
	db.Where("update_time <= ?", time.Now().Unix()-interval).
		Order("update_time").
		Limit(limit).
		Find(proxy)
	return proxy
}

func SaveValidProxy(proxy string) error {
	t := time.Now().Unix()
	_, err := Save(&model.ValidProxy{
		Content:               proxy,
		AssessTimes:           0,
		SuccessTimes:          0,
		AvgResponseTime:       0,
		ContinuousFailedTimes: 0,
		InsertTime:            t,
		Score:                 ProxyInitScore,
	})
	return err
}

func PopCrawlProxy(offset int, limit int) *[]model.CrawlProxy {
	s := GetCrawlProxy(offset, limit)
	DeleteCrawlProxy(s)
	return s
}

func DeleteCrawlProxy(p *[]model.CrawlProxy) error {
	var ids []int64
	db := GetSQLite()
	for _, v := range *p {
		ids = append(ids, v.ID)
	}
	return db.Where("id in (?)", ids).Delete(&model.CrawlProxy{}).Error
}

func GetCrawlProxy(offset int, limit int) *[]model.CrawlProxy {
	proxy := &[]model.CrawlProxy{}
	db := GetSQLite()
	db.Model(proxy).
		Offset(offset).
		Limit(limit).Find(proxy)
	return proxy
}

func SaveCrawlProxy(proxy string) error {
	t := time.Now().Unix()
	_, err := Save(&model.CrawlProxy{
		Content:    proxy,
		InsertTime: t,
		UpdateTime: t,
	})
	return err
}

func Save(d interface{}) (interface{}, error) {
	tx := GetSQLite().Begin()
	if err := tx.Create(d).Error; err != nil {
		log.Errorln(err)
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return d, nil
}
