// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
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
