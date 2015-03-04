// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"./util"
	"fmt"
	"os"
)

const (
	CredFile = ".mongodb_mms"
)

func main() {
	config, err := util.LoadConfigFromHome(CredFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}

	username, apikey := config.GetCredentials()
	api, err := util.NewMMSAPI("https://mms.mongodb.com", username, apikey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}

	hosts, err := api.GetAllHosts("5363cd319194bf134f77e6e0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}

	fmt.Fprintf(os.Stdout, "%v\n", hosts)
}
