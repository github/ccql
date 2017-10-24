# ccql

[![downloads](https://img.shields.io/github/downloads/github/ccql/total.svg)](https://github.com/github/ccql/releases/latest)

#### Concurrent, multi server MySQL client

`ccql` is a simple executable utility which executes a given set of queries on a given set of MySQL hosts
 in parallel.

Quick example:
```
echo "my.srv1.com my.srv2.com my.srv3.com" | ccql -q "show master status; select @@server_id" -u myuser -p 123456
```

## Usage

```
Usage of ccql:
  -C string
        Credentials file, expecting [client] scope, with 'user', 'password' fields. Overrides -u and -p
  -H string
    	Hosts file, hostname[:port] comma or space or newline delimited format. If not given, hosts read from stdin
  -Q string
    	Query/queries input file
  -d string
      Schemas to use (default "information_schema")
  -h string
    	Comma or space delimited list of hosts in hostname[:port] format. If not given, hosts read from stdin
  -m uint
    	Max concurrent connections (default 32)
  -p string
    	MySQL password
  -q string
    	Query/queries to execute
  -t float
    	Connect timeout seconds
  -u string
    	MySQL username (default OS user)
```

#### Hosts input

You may provide a list of hosts in the following ways:
- via `-h my.srv1.com:3307 my.srv2.com my.srv3.com`
- via `-H /path/to/hosts.txt`
- via _stdin_, as in `echo "my.srv1.com:3307 my.srv2.com my.srv3.com" | ccql ...`

Hostnames can be separated by spaces, commas, newline characters or all the above.
They may indicate a port. The default port, if unspecified, is `3306`

#### Queries input

You may provide a query or a list of queries in the following ways:
- single query, `-q "select @@global.server_id"`
- multiple queries, semicolon delimited: `-q "select @@global.server_id; set global slave_net_timeout:=10"`
- single or mutiple queries from text file: `-Q /path/to/queries.sql`

Queries are delimited by a semicolon (`;`). The last query may, but does not have to, be terminated by a semicolon.
Quotes are respected, up to a reasonable level. It is valid to include a semicolon in a quoted text, as in `select 'single;query'`. However `ccql` does not employ a full blown parser, so please don't overdo it. For example, the following may not be parsed correctly: `select '\';\''`. You get it.

#### Credentials input

You may provide credentials in the following ways:
- via `-u myusername -p mypassword` (default username is your OS user; default password is empty)
- via credentials file: `-C /path/to/.my.cnf`. File must be in the following format:
  ```
  [client]
  user=myuser
  password=mypassword
  ```

#### Execution

Hosts are executed in parallel, with up to `128` concurrent executions (otherwise more hosts are accepted but wait in queue).
For each host, the set of queries executes sequentially. Error on any query terminates execution of that host.
Errors are isolated to hosts; an error while connecting or executing on host1 should not affect execution on host2.

#### Output

There is only output generated for queries that provide an output, typically `SELECT` queries. Queries such as
`SET GLOBAL...` or `FLUSH BINARY LOGS` or `CREATE DATABASE ...` do not generate and output.

Output is written to _stdout_. It is tab delimited. There is one output line per row returning from either query.
The first printed token is the fully qualified `hostname:port` of the instance whose query output is printed.
Remember that execution happens concurrently on multiple hosts. Output rows are therefore ordered arbitrarily
in between hosts, though deterministically for any specific host.
Other tokens are whatever columns were returned by the queries.

## More examples

Some examples dealing with replication follow. Combining shell scripting we can have some real fun.
For brevity, we assume `/tmp/hosts.txt` contains a list of servers, as follows:
```
echo "localhost:22293, localhost:22294, localhost:22295, localhost:22296" > /tmp/hosts.txt
```
(note that hosts can be separated by spaces, commas, newlines or any combination)

We also assume credentials are stored in `/etc/ccql.cnf`:
```
[client]
user=msandbox
password=msandbox
```

Warmup: select some stuff
```
cat /tmp/hosts.txt | ccql -C /etc/ccql.cnf -q "select @@global.server_id, @@global.binlog_format, @@global.version"
```
A sample output is:
```
localhost:22296	103	STATEMENT	5.6.28
localhost:22294	101	STATEMENT	5.6.28-log
localhost:22293	1	STATEMENT	5.6.28-log
localhost:22295	102	STATEMENT	5.6.28-log
```
The output is tab delimited.

Show only servers that are configured as replicas:
```
cat /tmp/hosts.txt | ccql -C /etc/ccql.cnf -q "show slave status" | awk '{print $1}'
```
Apply `slave_net_timeout` only on replicas:
```
cat /tmp/hosts.txt | ccql -C /etc/ccql.cnf -q "show slave status;" | awk '{print $1}' | ccql -C /etc/ccql.cnf -q "set global slave_net_timeout := 10"
```

Getting tired of typing `ccql -C /etc/ccql.cnf`? Let's make a shortcut:
```
alias ccql="ccql -C /etc/ccql.cnf"
```

Which servers are acting as masters to someone?
```
cat /tmp/hosts.txt | ccql -q "show slave status;" | awk -F $'\t' '{print $3 ":" $5}'
```

Of those, which are also replicating? i.e. act as intermediate masters?
```
cat /tmp/hosts.txt | ccql -q "show slave status;" | awk -F $'\t' '{print $3 ":" $5}' | sort | uniq | ccql -q "show slave status" | awk '{print $1}'
```

Set `sync_binlog=0` on all intermediate masters:
```
cat /tmp/hosts.txt | ccql -q "show slave status;" | awk -F $'\t' '{print $3 ":" $5}' | sort | uniq | ccql -q "show slave status" | awk '{print $1}' | ccql -q "set global sync_binlog=0"
```

## LICENSE

See [LICENSE](LICENSE). _ccql_ imports and includes 3rd party libraries, which have their own license. These are found under [vendor](vendor).

## Binaries, downloads

Find precompiled binaries for linux (amd64) and Darwin (aka OS/X, amd64) under [Releases](https://github.com/github/ccql/releases)

## Build

_ccql_ is built with Go 1.6/1.7, and supports 1.5 with the [Go 1.5 vendor directories](https://golang.org/cmd/go/#hdr-Vendor_Directories), which requires setting `GO15VENDOREXPERIMENT=1`.
Please see the [build file](build.sh)

## What's in a name?

_ccql_ is an abbreviation for _Concurrent Client for MySQL_ or something. We had a few iterations with the name
but had to replace one and we were all like _yeah_ and _whoa_ and fun times. Eventually we came by this name
which upset [tomkrouper](https://github.com/tomkrouper) being "too much on the left-side of the keyboard when typing" and that settled the matter.
Tom uses `alias a='ccql'`.

## Notes

- Credits to Domas Mituzas for creating [pmysql](http://dom.as/2010/08/12/pmysql-multi-server-mysql-client/).
This project mostly reimplements `pmysql` and delivers it in an easy to redistribute format.

- Pronounce "see-sequel"

- This project is open to [contributions](CONTRIBUTING.md). Generally speaking it should be kept small
and simple.
