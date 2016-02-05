package logic

import (
	"fmt"
	"strings"

	"github.com/outbrain/golib/log"
	"github.com/outbrain/golib/sqlutils"
)

// queryHost connects to a given host, issues the given set of queries, and outputs the results
// line per row in tab delimited format
func queryHost(host string, user string, password string, defaultSchema string, queries []string, timeout uint) error {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%ds", user, password, host, defaultSchema, timeout)
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

// QueryHosts will issue concurrent queries on given list of hosts
func QueryHosts(hosts []string, user string, password string, defaultSchema string, queries []string, maxConcurrency uint, timeout uint) {
	concurrentHosts := make(chan bool, maxConcurrency)
	completedHosts := make(chan bool)

	for _, host := range hosts {
		go func(host string) {
			concurrentHosts <- true
			queryHost(host, user, password, defaultSchema, queries, timeout)
			<-concurrentHosts

			completedHosts <- true
		}(host)
	}
	// Barrier. Wait for all to complete
	for range hosts {
		<-completedHosts
	}
}
