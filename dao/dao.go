package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/storyicon/golang-proxy/model"
	"github.com/storyicon/golang-proxy/toolkit"
)

const (
	SQLiteDatabase = "sqlite3.db"
)

var (
	SourcesPath string
	Sources     *model.Sources
	SQLite      *gorm.DB
)

func GetSQLite() *gorm.DB {
	if SQLite == nil {
		SQLite = toolkit.NewSQLite(SQLiteDatabase)
	}
	return SQLite
}

func GetSources() *model.Sources {
	if Sources == nil {
		SourcesPath = toolkit.GetSourcePath()
		Sources = toolkit.GetSources(SourcesPath)
	}
	return Sources
}

func init() {
	sqlite := GetSQLite()
	sqlite.AutoMigrate(
		&model.ValidProxy{},
		&model.CrawlProxy{},
	)
}
