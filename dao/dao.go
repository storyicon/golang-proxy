package dao

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang-proxy/model"
	"golang-proxy/std"
)

var (
	SourceFolderPath string
	ConfigFilePath   string
)

var (
	Database *gorm.DB
	Config   *model.Config
	Sources  model.Sources
)

func GetDatabase() *gorm.DB {
	if Database == nil {
		Database = getDatabase()
	}
	return Database
}

func GetSources() model.Sources {
	if Sources == nil {
		Sources = getSources()
	}
	return Sources
}

func GetConfig() *model.Config {
	if Config == nil {
		config, err := getConfig()
		if err != nil {
			log.Panicf("Failed to load config file, %v", err)
		}
		Config = config
	}
	return Config
}

func getDatabase() (db *gorm.DB) {
	config, err := getConfig()
	if err == nil {
		log.Infoln("The configuration file has been read")
		log.Infoln("Try to use the MySQL database")
		db, err = std.NewMySQL(config.MySQL)
		if err == nil {
			return db
		}
		log.Errorf("An error occurred while connecting to the MySQL database: %v", err)

		log.Infoln("Try to use the PostgreSQL database")
		db, err = std.NewPostgreSQL(config.PostgreSQL)
		if err == nil {
			return db
		}
		log.Errorf("An error occurred while connecting to the PostgreSQL database: %v", err)

	}
	log.Infof("Start to use sqlite database, because %v", err)
	db, err = std.NewSQLite(model.SQLiteDatabase)
	if err != nil {
		log.Panicf("An error occurred while initializing database SQLite: %v", err)
	}
	return
}

func getConfig() (*model.Config, error) {
	config := model.Config{}
	if err := std.LoadYaml(getConfigFilePath(), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getSources() model.Sources {
	sources := model.Sources{}
	filepath.Walk(getSourceFolderPath(), func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if filename := info.Name(); filepath.Ext(filename) == ".yml" && !strings.HasPrefix(filename, ".") {
			source := model.Source{}
			std.LoadYaml(path, &source)
			source.Name = filename
			sources = append(sources, &source)
			log.Infof("Successfully load source: %s", filename)
		}
		return nil
	})
	return sources
}

func getSourceFolderPath() string {
	if SourceFolderPath == "" {
		path := filepath.Join(std.GetCurrentDirectory(), "source")
		if !std.IsDirExists(path) {
			log.Errorf(`Source folder "%s" does not exist`, path)
			os.Exit(1)
		}
		SourceFolderPath = path
	}
	return SourceFolderPath
}

func getConfigFilePath() string {
	if ConfigFilePath == "" {
		ConfigFilePath = filepath.Join(std.GetCurrentDirectory(), "config.yml")
	}
	return ConfigFilePath
}

func init() {
	session := GetDatabase()
	session.AutoMigrate(
		&model.CrudeProxy{},
		&model.Proxy{},
	)
}
