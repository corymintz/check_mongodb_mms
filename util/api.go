// Copyright 2015 MongoDB, Inc. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"../model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type MMSAPI struct {
	client   *http.Client
	hostname string
}

func NewMMSAPI(hostname string, timeout int, username string, apiKey string) (*MMSAPI, error) {
	t := NewTransport(username, apiKey)
	c, err := t.Client()
	if err != nil {
		return nil, err
	}

	// Setup our own Transport to ensure that the timeout is respected
	// It may seem excesive to set a dialer timeout, a deadline on the
	// connection, and a response header timeout, but we have seen
	// problems in the MongoDB MMS Backup Agent that required all three
	// of these.
	t.Transport = &http.Transport{
		Dial: func(network, addr string) (conn net.Conn, err error) {
			conn, err = net.DialTimeout(network, addr, time.Duration(timeout)*time.Second)
			if err != nil {
				return conn, err
			}

			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			return conn, nil
		},
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
	}

	return &MMSAPI{client: c, hostname: hostname}, nil
}

func (api *MMSAPI) GetAllHosts(groupId string) ([]model.Host, error) {
	body, err := api.doGet(fmt.Sprintf("/groups/%v/hosts", groupId))
	if err != nil {
		return nil, err
	}

	hostResp := &model.HostsResponse{}
	if err := json.Unmarshal([]byte(body), &hostResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response did not contain valid JSON. Error: %v, Body: %v", err, body))
	}

	return hostResp.Hosts, nil
}

func (api *MMSAPI) GetHostByName(groupId string, name string) (*model.Host, error) {
	body, err := api.doGet(fmt.Sprintf("/groups/%v/hosts/byName/%v", groupId, name))
	if err != nil {
		return nil, err
	}

	host := &model.Host{}
	if err := json.Unmarshal([]byte(body), &host); err != nil {
		return nil, errors.New(fmt.Sprintf("Response did not contain valid JSON. Error: %v, Body: %v", err, body))
	}

	return host, nil
}

func (api *MMSAPI) GetHostMetric(groupId string, hostId string, metricName string) (*model.Metric, error) {
	body, err := api.doGet(fmt.Sprintf("/groups/%v/hosts/%v/metrics/%v", groupId, hostId, metricName))
	if err != nil {
		return nil, err
	}

	metric := &model.Metric{}
	if err := json.Unmarshal([]byte(body), &metric); err != nil {
		return nil, errors.New(fmt.Sprintf("Response did not contain valid JSON. Error: %v, Body: %v", err, body))
	}

	return metric, nil
}

func (api *MMSAPI) doGet(path string) (string, error) {
	uri := fmt.Sprintf("%v/api/public/v1.0%v", api.hostname, path)

	response, err := api.client.Get(uri)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to make HTTP request. Error: %v", err))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read HTTP response body. Error: %v", err))
	}

	if response.StatusCode != 200 {
		return "", handleError(response.StatusCode, string(body[:]))
	}

	return string(body[:]), nil
}

func handleError(statusCode int, body string) error {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonBody); err != nil {
		return errors.New(fmt.Sprintf("API response did not contain valid JSON. Body: %v", body))
	}

	return errors.New(fmt.Sprintf("API Error: %v (%v)", jsonBody["reason"], jsonBody["detail"]))
}
