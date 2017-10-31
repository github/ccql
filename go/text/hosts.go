package text

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	defaultMySQLPort = 3306
)

// ParseHosts will return list of hostnames from either given list,
// or file, or stdin.
func ParseHosts(hostsList string, hostsFile string) (hosts []string, err error) {
	if hostsFile != "" {
		bytes, err := ioutil.ReadFile(hostsFile)
		if err != nil {
			return hosts, err
		}
		hostsList = string(bytes)
	}
	if hostsList == "" {
		// Read from stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return hosts, err
		}
		hostsList = string(bytes)
	}
	parsedHosts := regexp.MustCompile("[,\\s]").Split(hostsList, -1)
	for _, host := range parsedHosts {
		host = strings.TrimSpace(host)
		if host != "" {
			if !strings.Contains(host, ":") {
				host = fmt.Sprintf("%s:%d", host, defaultMySQLPort)
			}
			hosts = append(hosts, host)
		}
	}

	return hosts, err
}

func SplitNonEmpty(s string, sep string) (result []string) {
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		if token != "" {
			result = append(result, strings.TrimSpace(token))
		}
	}
	return result
}
