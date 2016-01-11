package sql

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// ParseQueries parses multiple queries given a single string.
// Input may be a semicolon delimited set of queries, such as:
// `select 1; show variables; select '2;nd'`
// Empty queries are ignored and skipped. Result is an array of single queries.
func ParseQueries(queriesText string, queriesFile string) (queries []string, err error) {
	if queriesFile != "" {
		bytes, err := ioutil.ReadFile(queriesFile)
		if err != nil {
			return queries, err
		}
		queriesText = string(bytes)
	}

	// The following regexp makes for a reasonable (yet incomplete) parses that splits
	// by delimiter ";" yet ignores it when it is quoted
	r := regexp.MustCompile(`([^;']+|'([^']*)')+[;]?`)
	matches := r.FindAllString(queriesText, -1)
	for _, match := range matches {
		match = strings.TrimSuffix(match, ";")
		match = strings.TrimSpace(match)
		if match != "" {
			queries = append(queries, match)
		}
	}
	return queries, nil
}
