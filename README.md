# Blockstack Crawler

Blockstack Crawler queries the Blockstack API and persists the data into database and storage.

It currently supports [Badger](https://github.com/dgraph-io/badger) or [Bolt](https://github.com/coreos/bbolt) as the underlying database.
Storage type is limited to local for now.

>Note: The current set of commands is limited, but the repository was designed with extensibility in mind.
Additional commands, as well as different database and storage types are a WIP.

##Commands
Below is a list of currently supported commands:

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
  - Get set of users from database or storage at the date specified by the `since` flag. 
  If no users exist at that date, the command will return an error and exit.
  - Retrieve set of users at current date. The database and storage will first be queried to check if the data already exists.
  If not, the remote Blockstack API is queried for the data.
  - The set of new users can be printed to standard output, or written to file in text or JSON format.
  
Currently, fetching the latest users takes a short period of time (3 seconds to 1-2 minutes).

##### Daily New Users
It is recommended to use an external process manager such as `systemd`, `init.d` or `upstart` to run this program as a background daemon.

`cron` can be used to query and report new users each day

An example Dockerfile is provided in `tools/docker` and can be used to launch a container that reports daily new users.


#### 