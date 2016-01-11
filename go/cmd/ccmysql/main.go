package main

import (
	"flag"

	"github.com/github/ccmysql/go/logic"
	"github.com/github/ccmysql/go/sql"
	"github.com/github/ccmysql/go/text"

	"github.com/outbrain/golib/log"
	"os/user"
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
	hostsList := flag.String("h", "", "Comma or space delimited list of hosts in hostname[:port] format. If not given, hosts read from stdin")
	hostsFile := flag.String("H", "", "Hosts file, hostname[:port] comma or space or newline delimited format. If not given, hosts read from stdin")
	queriesText := flag.String("q", "", "Query/queries to execute")
	queriesFile := flag.String("Q", "", "Query/queries input file")
	timeout := flag.Int("t", 0, "Connect timeout seconds")
	flag.Parse()

	if *queriesText == "" && *queriesFile == "" {
		log.Fatalf(`You must provide a query via -q "<some query>" or via -Q <query-file>`)
	}
	if *hostsList != "" && *hostsFile != "" {
		log.Fatalf("Both -q and -Q given. Please specify exactly one")
	}
	queries, err := sql.ParseQueries(*queriesText, *queriesFile)
	if err != nil {
		log.Fatale(err)
	}
	if len(queries) == 0 {
		log.Fatalf("No query/queries given")
	}

	if *hostsList != "" && *hostsFile != "" {
		log.Fatalf("Both -h and -H given. Please choose either one, or none (in which case stdin is used)")
	}
	hosts, err := text.ParseHosts(*hostsList, *hostsFile)
	if err != nil {
		log.Fatale(err)
	}
	if len(hosts) == 0 {
		log.Fatalf("No hosts given")
	}

	logic.QueryHosts(hosts, *user, *password, queries, *timeout)
}
