package library

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/nsqio/go-nsq"
	yaml "gopkg.in/yaml.v2"
)

func loadYaml(path string, data interface{}) {
	_, err := os.Stat(path)
	if err != nil {
		log.Panicln("Can not load yaml file at " + path)
	}
	yamlFile, _ := ioutil.ReadFile(path)
	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		log.Panicln("Can not parse yaml file", err)
	}
}

func GetConfig() *Config {
	c := Config{}
	environment := os.Getenv("SCHEDULER_ENV")
	if environment == "" {
		environment = "local"
		log.Println("Use Default Environment: local")
	}
	path := "./config/" + environment + ".yml"
	loadYaml(path, &c)
	return &c
}

func NewNSQProduer(address string) *nsq.Producer {
	p, err := nsq.NewProducer(address, nsq.NewConfig())
	if err == nil {
		err = p.Ping()
	}
	if err != nil {
		log.Panicln("Unable to initialize connection to NSQ")
	}
	return p
}

func NewNSQConsumer(addr string, topic string, chanel string, handler nsq.Handler) *nsq.Consumer {
	g := nsq.NewConfig()
	g.Set("max_in_flight", 2000)
	g.Set("rdy_redistribute_interval", time.Millisecond*5)
	c, err := nsq.NewConsumer(topic, chanel, g)
	if err != nil {
		panic(err)
	}
	c.AddHandler(handler)
	if err := c.ConnectToNSQLookupd(addr); err != nil {
		log.Panicln("Unable to initialize connection to NSQ", err)
	}
	return c
}

func GetSource() *[]Source {
	rootpath := "./source/"
	var s Source
	var r []Source
	filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		filename := info.Name()
		if strings.HasSuffix(filename, ".yml") && !strings.HasPrefix(filename, ".") {
			file := rootpath + filename
			loadYaml(file, &s)
			s.Name = filename
			r = append(r, s)
			log.Println("Successfully Load Source:", file)
		}
		return nil
	})
	return &r
}

func SleepRangeTime(delayRange []int) {
	delay := delayRange[0] + rand.Intn(delayRange[1]-delayRange[0])
	if delay != 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func JSONDecode(str string, variable interface{}) error {
	return json.Unmarshal([]byte(str), variable)
}

func JSONEncode(variable interface{}) ([]byte, error) {
	str, err := json.Marshal(variable)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func ProxyStringify(p *Proxy) string {
	return p.Protocal + "://" + p.IP + ":" + p.Port
}

type MySQL struct {
	Connection *gorm.DB
}

func NewMySQL(dsn string) *MySQL {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Panicln("Unable to initialize connection to MySQL", err)
	}
	db.SingularTable(true)
	db.LogMode(false)
	return &MySQL{
		Connection: db,
	}
}

func GetMysqlDsn(c *Config) string {
	m := c.Mysql
	dsn := m.User + ":" + m.Pass + "@(" + m.Host + ":" + m.Port + ")/" + m.Db + "?charset=" + m.Charset
	return dsn
}

func Contains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
