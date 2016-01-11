package logic

import (
	"fmt"
	"strings"

	"github.com/outbrain/golib/log"
	"github.com/outbrain/golib/sqlutils"
)

const (
	maxConcurrentConnections = 128
)

// queryHost connects to a given host, issues the given set of queries, and outputs the results
// line per row in tab delimited format
func queryHost(host string, user string, password string, queries []string, timeout int) error {
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

//  QueryHosts will issue concurrent queries on given list of hosts
func QueryHosts(hosts []string, user string, password string, queries []string, timeout int) {
	concurrentHosts := make(chan bool, maxConcurrentConnections)
	completedHosts := make(chan bool)

	for _, host := range hosts {
		go func(host string) {
			concurrentHosts <- true
			queryHost(host, user, password, queries, timeout)
			<-concurrentHosts

			completedHosts <- true
		}(host)
	}
	// Barrier. Wait for all to complete
	for range hosts {
		<-completedHosts
	}
}
