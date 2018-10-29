package business

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	// PortRegexString is used for regular matching port numbers.
	PortRegexString = `([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])`
	// IPRegexString is used for regular matching IP address.
	IPRegexString = `((?:(?:25[0-5]|2[0-4]\d|(?:1\d{2}|[1-9]?\d))\.){3}(?:25[0-5]|2[0-4]\d|(?:1\d{2}|[1-9]?\d)))`
)

const (
	typeHTTP = iota
	typeHTTPS
	typeBOTH
)

var (
	// PortRegexp is the compiled port regular expression.
	PortRegexp = regexp.MustCompile(PortRegexString)
	// IPRegexp is the compiled IP regular expression.
	IPRegexp = regexp.MustCompile(IPRegexString)
)

// Proxy defines the data structure of proxy.
type Proxy struct {
	// IP is an ip address in the strict sense, in the form of xxx.xxx.xxx.xxx
	IP string
	// Port is a port number in the strict sense
	Port string
	// In order to avoid the burden of source configuration and untrusted scheme data,
	// we decided to hide this data item and automatically detect the proxy's scheme by the assessor module.
	// Scheme string
}

// String is used to convert Proxy to strings.
func (proxy *Proxy) String() string {
	if proxy.IP == "" || proxy.Port == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", proxy.IP, proxy.Port)
}

// NewProxy is used to return a Proxy object by parsing and formatting
func NewProxy(ipRaw string, portRaw string) *Proxy {
	ip, port := parseProxyByIPRaw(ipRaw)
	if ip == "" {
		return nil
	}
	if port == "" {
		port = parsePortByPortRaw(portRaw)
	}
	return &Proxy{
		IP:   ip,
		Port: port,
	}
}

// Parsing out the port from the Port string
func parsePortByPortRaw(portRaw string) string {
	return PortRegexp.FindString(portRaw)
}

// Parsing out the IP from the IP string
func parseIPByIPRaw(ipRaw string) string {
	return IPRegexp.FindString(ipRaw)
}

// parseProxyByIPRaw is used to parse out ip and port from ip string
// Accept parameters similar to "127.0.0.1", "127.0.0.1:1080", "http://127.0.0.1" ...
func parseProxyByIPRaw(ipRaw string) (ip string, port string) {
	if strings.Index(ipRaw, "://") == -1 {
		ipRaw = "http://" + ipRaw
	}
	if parser, err := url.Parse(ipRaw); err == nil {
		ip = parseIPByIPRaw(parser.Hostname())
		port = parser.Port()
	}
	return
}
