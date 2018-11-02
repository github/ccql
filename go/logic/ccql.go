package logic

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/outbrain/golib/sqlutils"
)

// queryHost connects to a given host, issues the given set of queries, and outputs the results
// line per row in tab delimited format
func queryHost(
	host, user, password, schema string,
	queries []string,
	timeout float64, printSchema bool,
) error {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%fs", user, password, host, schema, timeout)
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
			if printSchema {
				output = append(output, schema)
			}
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
func QueryHosts(
	hosts, schemas, queries []string,
	user, password, defaultSchema string,
	maxConcurrency uint, timeout float64,
) (anyError error) {
	concurrentQueries := make(chan bool, maxConcurrency)
	printSchema := len(schemas) > 0
	if len(schemas) == 0 {
		schemas = []string{defaultSchema}
	}
	var wg sync.WaitGroup
	for _, host := range hosts {
		// For each host, run all queries for the respective schema
		for _, schema := range schemas {
			wg.Add(1)
			go func(host, schema string) {
				concurrentQueries <- true
				defer func() { <-concurrentQueries }()
				defer wg.Done()
				if err := queryHost(host, user, password, schema, queries, timeout, printSchema); err != nil {
					anyError = err
					log.Printf("%s %s", host, err.Error())
				}
			}(host, schema)
		}
	}

	// Barrier. Wait for all to complete
	wg.Wait()

	return anyError
}
