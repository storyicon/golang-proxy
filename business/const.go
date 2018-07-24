package business

const (
	RequestTimeout                  = 5
	AllowAssessSuccessRateMin       = 0.5
	ConsumerAssessTimes             = 3
	ConsumerStackCapacity           = 500
	ConsumerPerExtract              = 30
	AssessorStackCapacity           = 500
	AssessorPerExtract              = 30
	AssessorInterval          int64 = 60
	UserAgent                       = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"
	ResponseCodeError         int   = 1
	ResponseCodeSuccess       int   = 0
	ServiceListenPort               = 9999
)
