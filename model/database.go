package model

// CrudeProxy stores the agents that are crawled out, and cannot guarantee their quality.
type CrudeProxy struct {
	// ID is the ID value of the current record, which is unique among all proxies.
	ID int64 `gorm:"AUTO_INCREMENT;" json:"id"`
	// IP is the IP address of the proxy. e.g 127.0.0.1
	IP string `json:"ip"`
	// Port is the Port of the proxy. e.g 3306
	Port string `json:"port"`
	// Content is the ip:port of the proxy. e.g 127.0.0.1:3306
	Content string `gorm:"unique_index:unique_crude_content;" json:"content"`
	// InsertTime is the insertion time of the proxy
	InsertTime int64 `json:"insert_time"`
	// UpdateTime is the update time of the proxy
	UpdateTime int64 `json:"update_time"`
}

func (CrudeProxy) TableName() string {
	return CrudeProxyTableName
}

// Proxy stores the proxy filtered from CrudeProxy
type Proxy struct {
	// ID is the ID value of the current record, which is unique among all proxies.
	ID int64 `gorm:"AUTO_INCREMENT;" json:"id"`
	// IP is the IP address of the proxy. e.g 127.0.0.1
	IP string `json:"ip"`
	// Port is the Port of the proxy. e.g 3306
	Port string `json:"port"`
	// SchemeType represents the protocol type supported by the proxy.
	// 0: http
	// 1: https
	// 2: http & https
	SchemeType int64 `json:"scheme_type"`
	// Content is the ip:port of the proxy. e.g 127.0.0.1:3306
	Content string `gorm:"unique_index:unique_content;" json:"content"`

	// AssessTimes is the number of evaluations of the proxy
	AssessTimes int64 `json:"assess_times"`
	// SuccessTimes is the number of successful evaluations of the proxy
	SuccessTimes int64 `json:"success_times"`
	// AvgResponseTime is the average response time of the proxy
	AvgResponseTime float64 `json:"avg_response_time"`
	// ContinuousFailedTimes is the number of consecutive failures during the proxy evaluation process
	ContinuousFailedTimes int64 `json:"continuous_failed_times"`
	// Score is the rating of the proxy
	Score float64 `json:"score"`
	// InsertTime is the insertion time of the proxy
	InsertTime int64 `json:"insert_time"`
	// UpdateTime is the update time of the proxy, can also reflect the last evaluation time
	UpdateTime int64 `json:"update_time"`
}

func (Proxy) TableName() string {
	return ProxyTableName
}
