// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	const (
		timeoutDefault  = 10
		timeoutUsage    = "connection timeout connecting MMS/Ops Manager service"
		hostnameDefault = "https://mms.mongodb.com"
		hostnameUsage   = "hostname and port of the MMS/Ops Manager service"
		warningDefault  = math.MaxFloat64
		warningUsage    = "warning threshold for given metric"
		criticalDefault = math.MaxFloat64
		criticalUsage   = "critical threshold for given metric"
		metricDefault   = ""
		metricUsage     = "metric to query"
	)

	var timeout int
	flag.IntVar(&timeout, "timeout", timeoutDefault, timeoutUsage)
	flag.IntVar(&timeout, "t", timeoutDefault, timeoutUsage+" (shorthand)")

	var hostname string
	flag.StringVar(&hostname, "hostname", hostnameDefault, hostnameUsage)
	flag.StringVar(&hostname, "H", hostnameDefault, hostnameUsage+" (shorthand)")

	var warning float64
	flag.Float64Var(&warning, "warning", warningDefault, warningUsage)
	flag.Float64Var(&warning, "w", warningDefault, warningUsage+" (shorthand)")

	var critical float64
	flag.Float64Var(&critical, "critical", criticalDefault, criticalUsage)
	flag.Float64Var(&critical, "c", criticalDefault, criticalUsage+" (shorthand)")

	var metric string
	flag.StringVar(&metric, "metric", metricDefault, metricUsage)
	flag.StringVar(&metric, "m", metricDefault, metricUsage+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: check_mongodb_mms [-H hostname] [-t timeout] [-m metric] [-w warning_level] [-c critica_level]\n")
		fmt.Fprintf(os.Stdout, "    Options:\n")
		fmt.Fprintf(os.Stdout, "      -H, --hostname (default: %v) %v\n", hostnameDefault, hostnameUsage)
		fmt.Fprintf(os.Stdout, "      -t, --timeout (default: %v) %v\n", timeoutDefault, timeoutUsage)
		fmt.Fprintf(os.Stdout, "      -m, --metric (required) %v\n", metricUsage)
		fmt.Fprintf(os.Stdout, "      -w, --warning (default: %v) %v\n", warningDefault, warningUsage)
		fmt.Fprintf(os.Stdout, "      -c, --critical (default: %v) %v\n", criticalDefault, criticalUsage)
	}
	flag.Parse()

	if metric == "" {
		flag.Usage()
		os.Exit(2)
		return
	}
}
