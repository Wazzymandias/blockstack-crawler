# Blockstack Crawler

Blockstack Crawler queries the Blockstack API and persists relevant data into database and storage.

It currently supports [Badger](https://github.com/dgraph-io/badger) or [Bolt](https://github.com/coreos/bbolt) as the underlying database.
Storage type is limited to local for now.

> Note: Current functionality is limited, but the repository was designed with extensibility in mind.
> Additional commands, as well as different database and storage types are a WIP.

## Commands
Below is a list of currently supported commands:

#### Global Flags
```
Flags:
      --api-host string     host api to query (default "core.blockstack.org")
  -p, --api-port uint16     Port to access API (default 80)
      --api-scheme string   URL scheme for API (default "https")
      --datadir string      data directory that stores embedded database and storage information (default "$HOME/.blockstack-crawler/data")
      --db string           type of database to store results in (default "bolt")
  -h, --help                help for blockstack-crawler
      --store string        type of storage to use for persisting file data (default "local")
```
#### Names
```
Usage:
  blockstack-crawler names [OPTIONS] [flags]

Flags:
  -b, --batch uint         number of concurrent requests to make to API (default 50)
  -f, --format string      output names in prettified json or text format (default "json")
  -h, --help               help for names
  -o, --outfile string     write results to file rather than printing to standard output
  -s, --since string       ISO 8601 formatted date [YYYY-MM-DD]
  -t, --timeout duration   timeout for API requests (default 2m0s)

Global Flags:
      --api-host string     host api to query (default "core.blockstack.org")
  -p, --api-port uint16     Port to access API (default 80)
      --api-scheme string   URL scheme for API (default "https")
      --datadir string      data directory that stores embedded database and storage information (default "$HOME/.blockstack-crawler/data")
      --db string           type of database to store results in (default "bolt")
      --store string        type of storage to use for persisting file data (default "local")
```
#### Users
`./blockstack-crawler names`
>Note: Users are treated as names to conform to Blockstack API

  - Returns the set of users for each app at the current date. Database and storage are first queried for the data. If not found, a worker will query the remote Blockstack API for that information. 
  
  - Multiple requests are made and processed concurrently to improve performance. 
  
  - The aggregated data is then transformed, persisted in database and storage, and returned to the user.
  
#### New Users
`./blockstack-crawler names --since YYYY-MM-DD`

      -s, --since string    ISO 8601 formatted date [YYYY-MM-DD]
      
New users since a given date can be displayed and follows the general steps outlined below:
  - Retrieve set of users at current date. The database and storage will first be queried to check if the data already exists.
  If not, the remote Blockstack API is queried for the data.
  - Get set of users from database or storage at the date specified by the `since` flag. 
  If no users exist at that date, the command will return an error and exit.
  - The set of new users can be printed to standard output, or written to file in text or JSON format.
  
  >Note: Fetching the latest users may take anywhere from 3 seconds to 1-2 minutes depending on API response times, network load, etc.

##### Daily New Users
It is recommended to use an external process manager such as `systemd`, `init.d` or `upstart` to run this program as a background daemon.

`cron` can be used with the program to query and report new users each day.

An example Dockerfile and scripts are provided in `tools/docker` and can be used to launch a container that reports daily new users.

To use the image:
```
# clone the repo to $GOPATH
git clone https://github.com/Wazzymandias/blockstack-crawler $GOPATH/src/github.com/Wazzymandias/blockstack-crawler

# build the docker image
cd $GOPATH/src/github.com/Wazzymandias/blockstack-crawler/tools/docker

./build-user-watcher.sh

# start the docker image
./watch-users.sh
```

### Install

#### Requirements
  - [Go 1.10+](https://golang.org/dl/)
  - [Glide](https://github.com/Masterminds/glide)
```
# clone the repo
git clone https://github.com/Wazzymandias/blockstack-crawler $GOPATH/src/github.com/Wazzymandias/blockstack-crawler

# enter directory 
cd $GOPATH/src/github.com/Wazzymandias/blockstack-crawler/

# install glide if it doesn't exist
go get -u github.com/Masterminds/glide 

# install dependencies
glide install 

# build or install binary 
go build -v # OR 
go install -v
```