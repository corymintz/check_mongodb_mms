// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"time"
)

type Metric struct {
	MetricName string      `json:"metricName"`
	Units      string      `json:"units"`
	DataPoints []DataPoint `json:"dataPoints"`
}

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

var metricUnits = map[string]string{
	"RAW":          "",
	"BITS":         "b",
	"BYTES":        "B",
	"KILOBITS":     "kb",
	"KILOBYTES":    "KB",
	"MEGABITS":     "mb",
	"MEGABYTES":    "MB",
	"GIGABITS":     "gb",
	"GIGABYTES":    "GB",
	"TERABYTES":    "TB",
	"PETABYTES":    "PB",
	"MILLISECONDS": "ms",
	"SECONDS":      "secs",
	"MINUTES":      "mins",
	"HOURS":        "hours",
	"DAYS":         "days",
}

var metricFormaters = map[string]string{
	"ASSERT_MSG":                          "%v message asserts since process started",
	"ASSERT_REGULAR":                      "%v regular asserts since process started",
	"ASSERT_USER":                         "%v user asserts since process started",
	"ASSERT_WARNING":                      "%v warnings raised since process started",
	"BACKGROUND_FLUSH_AVG":                "%v millisecond background flush average",
	"COMPUTED_MEMORY":                     "%v megabytes non-mapped virtual memory",
	"CONNECTIONS":                         "%v active connections opened",
	"CURSORS_TOTAL_OPEN":                  "%v active cursors",
	"CURSORS_TOTAL_TIMED_OUT":             "%v cursor timeouts since process started",
	"DB_STORAGE_TOTAL":                    "%v bytes of on-disk storage used",
	"EFFECTIVE_LOCK_PERCENTAGE":           "%v effective lock percentage",
	"EXTRA_INFO_PAGE_FAULTS":              "%v page faults per second",
	"GLOBAL_ACCESSES_NOT_IN_MEMORY":       "%v not in memory page accesses per second",
	"GLOBAL_LOCK_CURRENT_QUEUE_READERS":   "%v queued readers",
	"GLOBAL_LOCK_CURRENT_QUEUE_TOTAL":     "%v queued total requests",
	"GLOBAL_LOCK_CURRENT_QUEUE_WRITERS":   "%v queued writers",
	"GLOBAL_PAGE_FAULT_EXCEPTIONS_THROWN": "%v page fault exceptions per second",
	"INDEX_COUNTERS_BTREE_ACCESSES":       "%v btree accesses per second",
	"INDEX_COUNTERS_BTREE_HITS":           "%v btree hits per second",
	"INDEX_COUNTERS_BTREE_MISSES":         "%v btree misses per second",
	"INDEX_COUNTERS_BTREE_MISS_RATIO":     "%v btree miss ratio",
	"JOURNALING_COMMITS_IN_WRITE_LOCK":    "%v journal commits in write lock",
	"JOURNALING_MB":                       "%v megabytes writen to journal per second",
	"MEMORY_MAPPED":                       "%v megabytes of mapped datafiles",
	"MEMORY_RESIDENT":                     "%v megabytes of resident memory used",
	"MEMORY_VIRTUAL":                      "%v megabyes of virtual memory used",
	"NETWORK_BYTES_IN":                    "%v incoming bytes per second",
	"NETWORK_BYTES_OUT":                   "%v outgoing bytes per second",
	"NETWORK_NUM_REQUESTS":                "%v requests per second",
	"OPCOUNTERS_CMD":                      "%v commands per second",
	"OPCOUNTERS_DELETE":                   "%v deletes per second",
	"OPCOUNTERS_GETMORE":                  "%v getmores per second",
	"OPCOUNTERS_INSERT":                   "%v inserts per second",
	"OPCOUNTERS_QUERY":                    "%v queries per second",
	"OPCOUNTERS_UPDATE":                   "%v updates per second",
	"OPCOUNTERS_REPL_CMD":                 "%v replicated commands per second",
	"OPCOUNTERS_REPL_DELETE":              "%v replicated deletes per second",
	"OPCOUNTERS_REPL_INSERT":              "%v replicated inserts per second",
	"OPCOUNTERS_REPL_UPDATE":              "%v replicated updates per second",
	"OPLOG_SLAVE_LAG_MASTER_TIME":         "%v seconds of replication lag",
	"OPLOG_MASTER_LAG_TIME_DIFF":          "%v seconds of replication headroom",
}

func (metric *Metric) ToStringLastDataPoint() string {
	if len(metric.DataPoints) == 0 {
		return "Metric has no datapoints"
	}

	return metric.ToStringDataPoint(len(metric.DataPoints) - 1)
}

func (metric *Metric) ToStringDataPoint(index int) string {
	metricFormater, ok := metricFormaters[metric.MetricName]
	if ok == false {
		return fmt.Sprintf("%v %v %v", metric.MetricName, metric.DataPoints[index].Value, metricUnits[metric.Units])
	}

	return fmt.Sprintf(metricFormater, metric.DataPoints[index].Value)
}
