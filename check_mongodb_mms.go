// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"./util"
	"flag"
	"fmt"
	"github.com/fractalcat/nagiosplugin"
	"math"
	"os"
)

const (
	CredFile = ".mongodb_mms"
)

var groupId string
var server string
var metric string
var hostname string
var warning float64
var critical float64
var timeout int

func main() {
	setupFlags()
	if metric == "" || server == "" || groupId == "" {
		flag.Usage()
		os.Exit(2)
		return
	}

	check := nagiosplugin.NewCheck()
	defer check.Finish()

	config, err := util.LoadConfigFromHome(CredFile)
	if err != nil {
		check.AddResultf(nagiosplugin.UNKNOWN, "%v", err)
		return
	}

	username, apikey := config.GetCredentials()
	api, err := util.NewMMSAPI(hostname, username, apikey)
	if err != nil {
		check.AddResultf(nagiosplugin.UNKNOWN, "Failed to create API. Err: %v", err)
		return
	}

	host, err := api.GetHostByName(groupId, server)
	if err != nil {
		check.AddResultf(nagiosplugin.UNKNOWN, "%v", err)
		return
	}

	check.AddResult(nagiosplugin.OK, fmt.Sprintf("%v", host))
}

func setupFlags() {
	const (
		groupIdDefault  = ""
		groupIdUsage    = "The MMS/Ops Manager group ID that contains the server"
		serverDefault   = ""
		serverUsage     = "hostname of the mongod/s to check"
		metricDefault   = ""
		metricUsage     = "metric to query"
		hostnameDefault = "https://mms.mongodb.com"
		hostnameUsage   = "hostname and port of the MMS/Ops Manager service"
		warningDefault  = math.MaxFloat64
		warningUsage    = "warning threshold for given metric"
		criticalDefault = math.MaxFloat64
		criticalUsage   = "critical threshold for given metric"
		timeoutDefault  = 10
		timeoutUsage    = "connection timeout connecting MMS/Ops Manager service"
	)

	flag.StringVar(&groupId, "groupid", groupIdDefault, groupIdUsage)
	flag.StringVar(&groupId, "g", groupIdDefault, groupIdUsage+" (shorthand)")

	flag.StringVar(&server, "server", serverDefault, serverUsage)
	flag.StringVar(&server, "s", serverDefault, serverUsage+" (shorthand)")

	flag.StringVar(&metric, "metric", metricDefault, metricUsage)
	flag.StringVar(&metric, "m", metricDefault, metricUsage+" (shorthand)")

	flag.StringVar(&hostname, "hostname", hostnameDefault, hostnameUsage)
	flag.StringVar(&hostname, "H", hostnameDefault, hostnameUsage+" (shorthand)")

	flag.Float64Var(&warning, "warning", warningDefault, warningUsage)
	flag.Float64Var(&warning, "w", warningDefault, warningUsage+" (shorthand)")

	flag.Float64Var(&critical, "critical", criticalDefault, criticalUsage)
	flag.Float64Var(&critical, "c", criticalDefault, criticalUsage+" (shorthand)")

	flag.IntVar(&timeout, "timeout", timeoutDefault, timeoutUsage)
	flag.IntVar(&timeout, "t", timeoutDefault, timeoutUsage+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: check_mongodb_mms  -g groupid -s server -m metric [-H hostname] [-t timeout] [-w warning_level] [-c critica_level]\n")
		fmt.Fprintf(os.Stdout, "     -g, --groupid  %v\n", groupIdUsage)
		fmt.Fprintf(os.Stdout, "     -s, --server  %v\n", serverUsage)
		fmt.Fprintf(os.Stdout, "     -m, --metric %v\n", metricUsage)
		fmt.Fprintf(os.Stdout, "     -H, --hostname (default: %v) %v\n", hostnameDefault, hostnameUsage)
		fmt.Fprintf(os.Stdout, "     -w, --warning (default: %v) %v\n", warningDefault, warningUsage)
		fmt.Fprintf(os.Stdout, "     -c, --critical (default: %v) %v\n", criticalDefault, criticalUsage)
		fmt.Fprintf(os.Stdout, "     -t, --timeout (default: %v) %v\n", timeoutDefault, timeoutUsage)
	}
	flag.Parse()
}
