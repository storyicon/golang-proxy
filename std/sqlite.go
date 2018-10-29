package std

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewSQLite(dbname string) (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", dbname)
	if err != nil {
		return
	}
	db.SingularTable(true)
	return db, nil
}
