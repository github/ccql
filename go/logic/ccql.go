package logic

import (
	"fmt"
	"log"
	"strings"

	"github.com/outbrain/golib/sqlutils"
)

// queryHost connects to a given host, issues the given set of queries, and outputs the results
// line per row in tab delimited format
func queryHost(host string, user string, password string, defaultSchema string, queries []string, timeout float64) error {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%fs", user, password, host, defaultSchema, timeout)
	db, _, err := sqlutils.GetDB(mysqlURI)
	if err != nil {
		return err
	}

	for _, query := range queries {
		resultData, err := sqlutils.QueryResultData(db, query)
		if err != nil {
			return err
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
func QueryHosts(hosts []string, user string, password string, defaultSchema string, queries []string, maxConcurrency uint, timeout float64) {
	concurrentHosts := make(chan bool, maxConcurrency)
	completedHosts := make(chan bool)

	for _, host := range hosts {
		go func(host string) {
			concurrentHosts <- true
			if err := queryHost(host, user, password, defaultSchema, queries, timeout); err != nil {
				log.Printf("%s %s", host, err.Error())
			}
			<-concurrentHosts

			completedHosts <- true
		}(host)
	}
	// Barrier. Wait for all to complete
	for range hosts {
		<-completedHosts
	}
}
