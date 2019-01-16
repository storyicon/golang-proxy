package std

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang-proxy/model"
	"log"
)

func NewPostgreSQL(options *model.PostgreSQLOptions) (db *gorm.DB, err error) {
	if options == nil {
		return nil, errors.New("mysql not set in config")
	}
	log.Println(options.String())
	db, err = gorm.Open("postgres", options.String())
	if err != nil {
		return
	}
	db.SingularTable(true)
	return
}
