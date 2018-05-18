package library

type Config struct {
	Mysql struct {
		Host    string
		Port    string
		User    string
		Pass    string
		Db      string
		Charset string
	}
	Nsq struct {
		ProducerAddr string `yaml:"producer_addr"`
		ConsumerAddr string `yaml:"consumer_addr"`
	}
	Debug bool
}

type Proxy struct {
	IP        string
	Port      string
	Protocal  string
	Proxytype int
}

type Source struct {
	Name string
	Page struct {
		Entry    string
		Template string
		From     int
		To       int
	}
	Selector struct {
		Iterator  string
		IP        string `yaml:"ip"`
		Port      string
		Protocal  string
		Proxytype string
		Filter    string
	}
	Category struct {
		Parallelnumber int
		DelayRange     []int `yaml:"delayRange"`
		Interval       string
	}
	Debug bool
}

type HTTPbinIP struct {
	Origin string `json:"origin"`
}

type ValidProxy struct {
	ID                    int     `gorm:"column:id"`
	Content               string  `gorm:"column:content"`
	AssessTimes           int     `gorm:"column:assess_times"`
	SuccessTimes          int     `gorm:"column:success_times"`
	AvgResponseTime       float64 `gorm:"column:avg_response_time"`
	ContinuousFailedTimes int     `gorm:"column:continuous_failed_times"`
	LastAssessTime        int64   `gorm:"column:last_assess_time"`
	Score                 float64 `gorm:"column:score"`
}
