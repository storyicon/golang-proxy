package dao

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/model"
)

func GetSQLResult(tableName string, sql string) (conseq interface{}, err error) {
	session := GetDatabase()
	switch tableName {
	case model.CrudeProxyTableName:
		conseq = &[]model.CrudeProxy{}
	case model.ProxyTableName:
		conseq = &[]model.Proxy{}
	default:
		return nil, errors.New("Query unknown table")
	}
	err = session.Raw(sql).Scan(conseq).Error
	return
}

func SaveProxy(proxy *model.Proxy) error {
	timestamp := time.Now().Unix()
	_, err := Save(&model.Proxy{
		IP:                    proxy.IP,
		Port:                  proxy.Port,
		SchemeType:            proxy.SchemeType,
		Content:               proxy.Content,
		AssessTimes:           0,
		SuccessTimes:          0,
		AvgResponseTime:       0,
		ContinuousFailedTimes: 0,
		InsertTime:            timestamp,
		Score:                 model.ProxyInitScore,
	})
	return err
}

func GetProxy(interval int64, limit int) []*model.Proxy {
	proxy := []*model.Proxy{}
	session := GetDatabase()
	session.Where("update_time <= ?", time.Now().Unix()-interval).
		Order("update_time").
		Limit(limit).
		Find(&proxy)
	return proxy
}

func GetCrudeProxy(offset int, limit int) []*model.CrudeProxy {
	proxies := []*model.CrudeProxy{}
	session := GetDatabase()
	session.Model(proxies).
		Offset(offset).
		Limit(limit).
		Find(&proxies)
	return proxies
}

func PopCrudeProxy(offset int, limit int) []*model.CrudeProxy {
	proxies := GetCrudeProxy(offset, limit)
	DeleteCrudeProxy(proxies)
	return proxies
}

func SaveCrudeProxy(proxy *model.CrudeProxy) error {
	timestamp := time.Now().Unix()
	if proxy.InsertTime == 0 {
		proxy.InsertTime = timestamp
	}
	_, err := Save(&model.CrudeProxy{
		IP:         proxy.IP,
		Port:       proxy.Port,
		Content:    proxy.Content,
		InsertTime: proxy.InsertTime,
		UpdateTime: timestamp,
	})
	return err
}

func DeleteCrudeProxy(proxies []*model.CrudeProxy) error {
	var idList []int64
	session := GetDatabase()
	for _, proxy := range proxies {
		idList = append(idList, proxy.ID)
	}
	return session.Where("id in (?)", idList).Delete(&model.CrudeProxy{}).Error
}

func Save(data interface{}) (interface{}, error) {
	session := GetDatabase().Begin()
	if err := session.Create(data).Error; err != nil {
		log.Errorln(err)
		session.Rollback()
		return 0, err
	}
	session.Commit()
	return data, nil
}
