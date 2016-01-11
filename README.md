# ccmysql
Multi server MySQL client

`ccmysql` is a simple utility which executes a given a set of queries on a given set of MySQL hosts
 in parallel.

Quick example:
```
echo "my.srv1.com my.srv2.com my.srv3.com" | ccmysql -q "show master status; select @@server_id" -u myuser -p 123456
```

## Usage
```
Usage of ccmysql:
  -H string
    	Hosts file, hostname[:port] comma or space or newline delimited format. If not given, hosts read from stdin
  -Q string
    	Query/queries input file
  -h string
    	Comma or space delimited list of hosts in hostname[:port] format. If not given, hosts read from stdin
  -p string
    	MySQL password
  -q string
    	Query/queries to execute
  -t int
    	Connect timeout seconds
  -u string
    	MySQL username (default OS user)
```

#### Hosts input

You may provide a list of hosts in the following ways:
- via `-h my.srv1.com:3307 my.srv2.com my.srv3.com`
- via `-H /path/to/hosts.txt`
- via _stdin_, as in `echo "my.srv1.com:3307 my.srv2.com my.srv3.com" | ccmysql ...`

Hostnames can be separated by spaces, commas, newline characters or all the above.
They may indicate a port. The default port, if unspecified, is `3306`

#### Queries input

You may provide a query or a list of queries in the following ways:
- single query, `-q "select @@global.server_id"`
- multiple queries, semicolon delimited: `-q "select @@global.server_id; set global slave_net_timeout:=10"`
- single or mutiple queries from text file: `-Q /path/to/queries.sql`

Queries are delimited by a semicolon (`;`). The last query may, but does not have to, be terminated by a semicolon.
Quotes are respected, up to a reasonable level. It is valid to include a semicolon in a quoted text, as in `select 'single;query'`. However `ccmysql` does not employ a full blown parser, so please don't overdo it. For example, the following may not be parsed correctly: `select '\';\''`. You get it.

#### Credentials

You may provide credentials in the following ways:
- via `-u myusername -p mypassword` (default username is your OS user; default password is empty)
- via credentials file: `-C /path/to/.my.cnf`. File must be in the following format:
  ```
  [client]
  user=myuser
  password=mypassword
  ```
