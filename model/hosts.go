// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

type Host struct {
	Id       string    `json:"id"`
	LastPing time.Time `json:"lastPing"`
}

type HostsResponse struct {
	Hosts []Host `json:"results"`
}
