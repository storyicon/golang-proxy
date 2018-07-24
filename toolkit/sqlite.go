package toolkit

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func NewSQLite(dbname string) *gorm.DB {
	var db *gorm.DB
	var err error
	if db, err = gorm.Open("sqlite3", dbname); err != nil {
		log.Panicln(err)
	}
    db.SingularTable(true)
	return db
}
