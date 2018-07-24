package business

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/storyicon/golang-proxy/dao"
)

type Response struct {
	Code    int
	Message interface{}
}

func StdExport(r http.ResponseWriter, data interface{}, err error) {
	var response Response
	r.Header().Set("Content-Type", "application/json")
	if err == nil {
		response.Code = ResponseCodeSuccess
		response.Message = data
		r.WriteHeader(http.StatusOK)
	} else {
		response.Code = ResponseCodeError
		response.Message = fmt.Sprintf("%v", err)
		r.WriteHeader(http.StatusInternalServerError)
	}
	var bytes []byte
	if bytes, err = json.Marshal(response); err == nil {
		r.Write(bytes)
	}
}

func SQLRedirect(res http.ResponseWriter, req *http.Request, sql string) {
	url := fmt.Sprintf("http://%s/sql?query=%s", req.Host, sql)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func getTableBySQL(s string) string {
	s = strings.ToLower(s)
	n := strings.Index(s, "from")
	if n < 0 {
		return ""
	}
	var w, c string
	for n = n + 4; n < len(s); n++ {
		c = s[n : n+1]
		if c != " " {
			w += c
			continue
		}
		if w != "" {
			break
		}
	}
	return w
}

func getRequestParam(req *http.Request, key string, _default string) string {
	if v := strings.Join(req.URL.Query()[key], ""); v != "" {
		return v
	}
	return _default
}

func StartService() {
	router := mux.NewRouter()
	router.HandleFunc("/all", func(res http.ResponseWriter, req *http.Request) {
		table := getRequestParam(req, "table", "valid_proxy")
		sql := fmt.Sprintf("SELECT * FROM %s ", table)
		SQLRedirect(res, req, sql)
	})
	router.HandleFunc("/random", func(res http.ResponseWriter, req *http.Request) {
		table := getRequestParam(req, "table", "valid_proxy")
		sql := fmt.Sprintf("SELECT * FROM %s ORDER BY RANDOM() limit 1", table)
		SQLRedirect(res, req, sql)
	})
	router.HandleFunc("/sql", func(res http.ResponseWriter, req *http.Request) {
		vars := req.URL.Query()
		if query := vars["query"]; len(query) == 1 {
			sql := strings.Join(query, "")
			if table := getTableBySQL(sql); table != "" {
				record, err := dao.GetSQLResult(table, sql)
				StdExport(res, record, err)
			}
		}
	})
	address := fmt.Sprintf("localhost:%d", ServiceListenPort)
	log.Infof("[S]Start Service on %s", address)
	http.ListenAndServe(address, router)
}
