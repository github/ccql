package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/github/ccql/go/logic"
	"github.com/github/ccql/go/sql"
	"github.com/github/ccql/go/text"

	golib_log "github.com/outbrain/golib/log"
	"gopkg.in/gcfg.v1"
)

var AppVersion string

const (
	maxAllowedConcurrentConnections uint = 128
)

// main is the application's entry point. It will either spawn a CLI or HTTP itnerfaces.
func main() {

	golib_log.SetLevel(golib_log.FATAL)

	osUser := ""
	// get os username as owner
	if usr, err := user.Current(); err == nil {
		osUser = usr.Username
	}

	help := flag.Bool("help", false, "Display usage")
	user := flag.String("u", osUser, "MySQL username")
	password := flag.String("p", "", "MySQL password")
	credentialsFile := flag.String("C", "", "Credentials file, expecting [client] scope, with 'user', 'password' fields. Overrides -u and -p")
	defaultSchema := flag.String("d", "information_schema", "Default schema to use")
	hostsList := flag.String("h", "", "Comma or space delimited list of hosts in hostname[:port] format. If not given, hosts read from stdin")
	hostsFile := flag.String("H", "", "Hosts file, hostname[:port] comma or space or newline delimited format. If not given, hosts read from stdin")
	queriesText := flag.String("q", "", "Query/queries to execute")
	queriesFile := flag.String("Q", "", "Query/queries input file")
	timeout := flag.Float64("t", 0, "Connect timeout seconds")
	maxConcurrency := flag.Uint("m", 32, "Max concurrent connections")
	flag.Parse()

	if AppVersion == "" {
		AppVersion = "local-build"
	}
	if *help {
		fmt.Fprintf(os.Stderr, "Usage of ccql (version: %s):\n", AppVersion)
		flag.PrintDefaults()
		return
	}
	if *queriesText == "" && *queriesFile == "" {
		fmt.Fprintf(os.Stderr, "You must provide a query via -q '<some query>' or via -Q <query-file>\n")
		fmt.Fprintf(os.Stderr, "Usage of ccql:\n")
		flag.PrintDefaults()
		return
	}
	if *hostsList != "" && *hostsFile != "" {
		log.Fatalf("Both -q and -Q given. Please specify exactly one")
	}
	queries, err := sql.ParseQueries(*queriesText, *queriesFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(queries) == 0 {
		log.Fatalf("No query/queries given")
	}

	if *hostsList != "" && *hostsFile != "" {
		log.Fatalf("Both -h and -H given. Please specify one of them, or none (in which case stdin is used)")
	}
	hosts, err := text.ParseHosts(*hostsList, *hostsFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(hosts) == 0 {
		log.Fatalf("No hosts given")
	}

	if *maxConcurrency > maxAllowedConcurrentConnections {
		log.Fatalf("Max concurrent connections (-m) may not exceed %d", maxAllowedConcurrentConnections)
	}
	if *maxConcurrency < 1 {
		*maxConcurrency = 1
	}

	if *credentialsFile != "" {
		mySQLConfig := struct {
			Client struct {
				User     string
				Password string
			}
		}{}
		gcfg.RelaxedParserMode = true
		err := gcfg.ReadFileInto(&mySQLConfig, *credentialsFile)
		if err != nil {
			log.Fatalf("Failed to parse gcfg data from file: %+v", err)
		} else {
			*user = mySQLConfig.Client.User
			*password = mySQLConfig.Client.Password
		}
	}

	if err := logic.QueryHosts(hosts, *user, *password, *defaultSchema, queries, *maxConcurrency, *timeout); err != nil {
		os.Exit(1)
	}
}
