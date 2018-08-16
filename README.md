# check_mongodb_mms
A Nagios plugin to check MongoDB metrics that have been collected by [MMS] (https://mms.mongodb.com) or [Ops Manager](https://www.mongodb.com/products/mongodb-enterprise-advanced), the On-Premise version of MMS.

The checks are executed by calling the [Public API](https://docs.mms.mongodb.com/reference/api/).

# Authentication
The API takes a username (most likely your email address) and an API key for authentication. Information on enabling the API for a MMS/Ops Manager group, as well as generating an API key, can be found at https://docs.mms.mongodb.com/tutorial/enable-public-api/.

# Build
Build with:
`go build check_mongodb_mms.go`

# Usage
The supported list of metric names can be found at https://docs.opsmanager.mongodb.com/current/reference/api/metrics/#entity-fields.

#### Help Output
    Usage: check_mongodb_mms  -g groupid -H hostname [-m metric] [-d dbname] [-a age] [-s server] [-t timeout] [-w warning_level] [-c critica_level] [-r granularity] [-p period] [-u username] [-k apikey]
     -g, --groupid  The MMS/Ops Manager group ID that contains the server
     -H, --hostname hostname:port of the mongod/s to check
     -m, --metric (no metric means check last ping age in seconds) metric to query
     -d, --dbname (default ) database name for DB_ metrics
     -a, --maxage (default 360) the maximum number of seconds old a metric before it is considered stale
     -s, --server (default: https://mms.mongodb.com) hostname and port of the MMS/Ops Manager service
     -w, --warning (default: ~:) warning threshold for given metric
     -c, --critical (default: ~:) critical threshold for given metric
     -t, --timeout (default: 10) connection timeout connecting MMS/Ops Manager service
     -r, --granularity (default: MINUTE) the size of the epoch. Acceptable values are MINUTE HOUR DAY
     -p, --period (default: 1H) the ISO-8601 formatted time period that specifies how far back in the past to query.
     -u, --username (default: ) the username for auth
     -k, --apiKey (default: ) the api key for the user

     -w and -c support the standard nagios threshold formats.
     See https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT for more details.

## Example Command Line Usage
MMS/Ops Manager not receiving a ping from a host is a warning after 180 seconds and critical after 300 seconds.

    ./check_mongodb_mms -g 54f84f43e6ccc36e22eef700 -H my-server.example.com:27017 -w 180 -c 300 -u username -k apikey

Delete Operations / Sec is considered a warning at 10 and critical at 25.

    ./check_mongodb_mms -g 54f84f43e6ccc36e22eef700 -H my-server.example.com:27017 -m OPCOUNTERS_DELETE -w 10 -c 25 -u username -k apikey

Virtual Memory usage is considered a warning at 8000 MB and critical at 10000 MB.

    ./check_mongodb_mms -g 54f84f43e6ccc36e22eef700 -H my-server.example.com:27017 -m MEMORY_VIRTUAL -w 8000 -c 10000 -u username -k apikey

## Example Nagios Config
    define command {
      command_nam e  check_mongodb_mms
      command_line  /usr/local/bin/check_mongodb_mms -g $ARG1$ -H $HOSTNAME$:$_HOSTPORT$ -m $ARG2$ -w $ARG3$ -c $ARG4$ -u $ARG5$ -k apikey $ARG6$
    }

    define service {
        use                 generic-service
        host_name           my-server.example.com
        service_description Inserts/Sec
        check_command       check_mongodb_mms!54f84f43e6ccc36e22eef700!OPCOUNTERS_INSERT!1000!1500!monitoringuser!1da09000d-d5da-1650-9a09-457e14ff7ef0
    }

    define host {
        use                     generic-host
        host_name               my-server.example.com
        alias                   my-server.example.com
        address                 127.0.0.1
        _PORT                   27017
    }
