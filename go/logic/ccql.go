package logic

import (
	"fmt"
	"log"
	"strings"

	"github.com/outbrain/golib/sqlutils"
	"sync"
)

// queryHost connects to a given host, issues the given set of queries, and outputs the results
// line per row in tab delimited format
func queryHost(host string, user string, password string, schema string, queries []string, timeout float64) error {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%fs", user, password, host, schema, timeout)
	fmt.Println(mysqlURI)
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
			output := []string{host, schema}
			for _, rowCell := range row {
				output = append(output, rowCell.String)
			}
			rowOutput := strings.Join(output, "\t")
			fmt.Println(rowOutput)
		}
	}
	return nil
}

func QuerySchemas(hosts []string, user string, password string, schemas []string, queries []string, maxConcurrency uint, timeout float64) (anyError error) {
	concurrentHosts := make(chan bool, maxConcurrency)
	completedHosts := make(chan bool)
	var wg sync.WaitGroup

	for _, host := range hosts {
		go func(host string) {
			wg.Add(len(schemas))
			concurrentHosts <- true
			// For each host, run all queries for the respective schema
			for _, schema := range schemas {
				go func(schema string) {
					if err := queryHost(host, user, password, schema, queries, timeout); err != nil {
						anyError = err
						log.Printf("%s %s", host, err.Error())
					}
					defer wg.Done()
				}(schema)
			}
			wg.Wait()
			<-concurrentHosts
			completedHosts <- true
		}(host)

	}

	// Barrier. Wait for all to complete
	for range hosts {
		<-completedHosts
	}

	return anyError
}
