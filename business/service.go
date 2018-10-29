package business

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/dao"
)

// Response is the response struct of the http service
type Response struct {
	Error   string      `json:"error"`
	Message interface{} `json:"message"`
}

// StartService used to start the http service
func StartService() {
	router := gin.Default()
	router.GET("/all", func(c *gin.Context) {
		tableName := queryWithDefault(c, "table", "proxy")
		sql := fmt.Sprintf("SELECT * FROM %s ", tableName)
		redirect(c, sql)
	})
	router.GET("/random", func(c *gin.Context) {
		tableName := queryWithDefault(c, "table", "proxy")
		sql := fmt.Sprintf("SELECT * FROM %s ORDER BY RANDOM() limit 1", tableName)
		redirect(c, sql)
	})
	router.GET("/sql", func(c *gin.Context) {
		query := c.Query("query")
		tableName := getTableNameBySQL(query)
		response, statusCode := Response{}, http.StatusOK
		if tableName != "" {
			record, err := dao.GetSQLResult(tableName, query)
			response.Message = record
			if err != nil {
				statusCode = http.StatusInternalServerError
				response.Error = fmt.Sprint(err)
			}
		} else {
			statusCode = http.StatusInternalServerError
			response.Error = "Unable to resolve table name"
		}
		c.JSON(statusCode, response)
	})
	log.Infof("[S]Start Service on %s", ServiceListenAddress)
	router.Run(ServiceListenAddress)
}

func redirect(context *gin.Context, sql string) {
	context.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/sql?query=%s", sql))
}

func getTableNameBySQL(s string) string {
	words := strings.Split(strings.ToLower(s), " ")
	length := len(words)
	for i := 0; i < length; i++ {
		if words[i] == "from" {
			if i < length-1 {
				return words[i+1]
			}
			break
		}
	}
	return ""
}

func queryWithDefault(c *gin.Context, key string, defaultValue string) string {
	if conseq := c.Query(key); conseq != "" {
		return conseq
	}
	return defaultValue
}
