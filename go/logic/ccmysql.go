package logic

import (
	"fmt"
	"strings"

	"github.com/outbrain/golib/log"
	"github.com/outbrain/golib/sqlutils"
)

func queryHost(host string, user string, password string, queries []string, timeout int) error {
	if host == "" {
		return nil
	}
	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:%d", host, defaultMySQLPort)
	}
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s)/?timeout=%ds", user, password, host, timeout)
	db, _, err := sqlutils.GetDB(mysqlURI)
	if err != nil {
		return log.Errore(err)
	}

	for _, query := range queries {
		resultData, err := sqlutils.QueryResultData(db, query)
		if err != nil {
			return log.Errore(err)
		}
		for _, row := range resultData {
			output := []string{host}
			for _, rowCell := range row {
				output = append(output, rowCell.String)
			}
			rowOutput := strings.Join(output, "\t")
			fmt.Println(rowOutput)
		}
	}
	return nil

}
