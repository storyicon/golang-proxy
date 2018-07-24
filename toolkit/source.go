package toolkit

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/model"
)

func GetSourcePath() (SourceFolderPath string) {
	var input string
	var err error
	flag.StringVar(&input, "source", "", "source foler`s path")
	flag.Parse()
	if input != "" {
		if input, err = filepath.Abs(input); err != nil {
			log.Panicln(err)
		} else if IsDirExists(input) {
			SourceFolderPath = input
			return
		}
		log.Warningf(`The entered resource path "%s" is invalid, use exec path`, input)
	}

	input = filepath.Join(getCurrentDirectory(), "source")
	if !IsDirExists(input) {
		log.Errorf(`Source path "%s" does not exist`, input)
		os.Exit(1)
	}
	SourceFolderPath = input
	return
}

func GetSources(path string) *model.Sources {
	var s model.Source
	var r model.Sources
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		filename := info.Name()
		if filepath.Ext(filename) == ".yml" && !strings.HasPrefix(filename, ".") {
			LoadYaml(path, &s)
			s.Name = filename
			r = append(r, s)
			log.Infoln("Successfully load source:", filename)
		}
		return nil
	})
	return &r
}
