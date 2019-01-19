package business

const (
	RequestTimeout = 4
	UserAgent      = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"

	ConsumerRetryTimes     = 3
	ConsumerStackCapacity  = 500
	ConsumerPerExtract     = 10
	ConsumerProxyInitScore = 1

	AssessorAllowSuccessRateMin       = 0.5
	AssessorStackCapacity             = 500
	AssessorPerExtract                = 10
	AssessorInterval            int64 = 60

	ServiceListenAddress = ":9999"
)
