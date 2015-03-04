// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

const (
	Delimiter = "="
)

type Config map[string]string

func LoadConfigFromHome(configFileName string) (Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to find home directory. Error: %v", err))
	}

	configFilePath := fmt.Sprintf("%v%c%v", usr.HomeDir, os.PathSeparator, configFileName)
	config, err := readConfig(configFilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to load %v. Error: %v", configFilePath, err))
	}

	return config, nil
}

func (config Config) GetCredentials() (string, string) {
	return config["username"], config["apikey"]
}

func readConfig(configFile string) (Config, error) {
	buffer, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	return doReadConfig(buffer), nil
}

func validConfigLine(line string) bool {
	return line != "" && line[0] != '#'
}

func partition(str string, delim string) (string, string, string) {
	idx := strings.Index(str, delim)
	if idx == -1 {
		return str, "", ""
	}

	return strings.TrimSpace(str[:idx]), delim, strings.TrimSpace(str[idx+len(delim):])
}

func readLines(input io.Reader) []string {
	ret := make([]string, 0, 0)
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret
}

func doReadConfig(buffer []byte) map[string]string {
	config := make(map[string]string)

	lines := readLines(bytes.NewBuffer(buffer))
	for _, line := range lines {
		if validConfigLine(line) {
			key, sep, value := partition(line, Delimiter)
			if sep != Delimiter {
				continue
			}

			config[key] = value
		}
	}

	return config
}
