package std

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"storyicon.visualstudio.com/golang-proxy/model"
)

func NewMySQL(options *model.MySQLOptions) (db *gorm.DB, err error) {
	log.Println(options.String())
	db, err = gorm.Open("mysql", options.String())
	if err != nil {
		return
	}
	db.SingularTable(true)
	return
}
