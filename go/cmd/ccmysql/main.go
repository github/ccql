package main

import (
	"flag"
	"fmt"
	"github.com/github/ccmysql/go/logic"
	"github.com/outbrain/golib/log"
	"github.com/outbrain/golib/sqlutils"
	"os/user"
	"strings"
)

const (
	maxConcurrentConnections = 128
	defaultMySQLPort         = 3306
)

// main is the application's entry point. It will either spawn a CLI or HTTP itnerfaces.
func main() {
	osUser := ""
	// get os username as owner
	if usr, err := user.Current(); err == nil {
		osUser = usr.Username
	}

	user := flag.String("u", osUser, "MySQL username")
	password := flag.String("p", "", "MySQL password")
	//hostsFile := flag.String("H", "", "Hosts file (read from stdin if not given); expected hostname[:port] per line")
	query := flag.String("q", "", "Query/queries to execute")
	timeout := flag.Int("t", 0, "Connect timeout seconds")
	flag.Parse()

	if *query == "" {
		log.Fatalf("You must provide a query via -q")
	}

	concurrentHosts := make(chan bool, maxConcurrentConnections)
	completedHosts := make(chan bool)

	var hosts []string = []string{"shlomi-gh:22295", "shlomi-gh:22296", "", "shlomi-gh"}

	for _, host := range hosts {
		go func(host string) {
			concurrentHosts <- true
			queryHost(host, *user, *password, *query, *timeout)
			<-concurrentHosts

			completedHosts <- true
		}(host)
	}
	// Barrier. Wait for all to complete
	for range hosts {
		<-completedHosts
	}
}
