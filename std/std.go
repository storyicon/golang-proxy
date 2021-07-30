package std

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

// TemplateRender is used to replace all the "key" in the template with "value"
func TemplateRender(template string, key string, value interface{}) string {
	return strings.Replace(template, "{"+key+"}", fmt.Sprint(value), -1)
}

// GetCurrentDirectory is used to get the directory where the current binary is running.
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	return dir
}

// LoadYaml is used to load yaml file
func LoadYaml(path string, data interface{}) (err error) {
	if _, err = os.Stat(path); err == nil {
		var bytes []byte
		if bytes, err = ioutil.ReadFile(path); err == nil {
			if err = yaml.Unmarshal(bytes, data); err == nil {
				return nil
			}
		}
	}
	return
}

func IsDirExists(path string) bool {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir()
	}
	return false
}

func Dump(variable interface{}) {
	if str, err := json.MarshalIndent(variable, "", "    "); err == nil {
		log.Print(string(str))
		return
	}
	log.Printf("%+v", variable)
}
