package model

type ValidProxy struct {
	ID                    int64   `gorm:"AUTO_INCREMENT;" json:"id"`
	Content               string  `gorm:"unique_index:unique_vp;" json:"content"`
	AssessTimes           int     `json:"assess_times"`
	SuccessTimes          int     `json:"success_times"`
	AvgResponseTime       float64 `json:"avg_response_time"`
	ContinuousFailedTimes int     `json:"continuous_failed_times"`
	Score                 float64 `json:"score"`
	InsertTime            int64   `json:"insert_time"`
	UpdateTime            int64   `json:"update_time"`
}

type CrawlProxy struct {
	ID         int64  `gorm:"AUTO_INCREMENT;" json:"id"`
	Content    string `gorm:"unique_index:unique_ip;" json:"content"`
	InsertTime int64  `json:"insert_time"`
	UpdateTime int64  `json:"update_time"`
}
