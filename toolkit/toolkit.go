package toolkit

import (
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func GetHostNameByIP(ip string) string {
	if u, err := url.Parse(ip); err == nil {
		return u.Hostname()
	}
	return ""
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panicln(err)
	}
	return filepath.FromSlash(dir)
}

func IsDirExists(path string) bool {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir()
	}
	return false
}

func LoadYaml(path string, data interface{}) {
	var err error
	if _, err = os.Stat(path); err == nil {
		var yamlFile []byte
		if yamlFile, err = ioutil.ReadFile(path); err == nil {
			if err = yaml.Unmarshal(yamlFile, data); err == nil {
				return
			}
		}
	}
	log.Errorf(`Can not import yml file at "%s": %v`, path, err)
	os.Exit(1)
}

func SleepRandomRangeTime(delayRange []int) {
	var t int
	switch len(delayRange) {
	case 0:
		t = 1
	case 1:
		t = delayRange[0]
	case 2:
		t = delayRange[0] + rand.Intn(delayRange[1]-delayRange[0])
	default:
		t = delayRange[rand.Intn(len(delayRange))]
	}
	time.Sleep(time.Duration(t) * time.Second)
}
