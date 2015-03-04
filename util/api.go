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
	"net/http"
)

type MMSAPI struct {
	client   *http.Client
	hostname string
}

func NewMMSAPI(hostname string, username string, apiKey string) (*MMSAPI, error) {
	t := NewTransport(username, apiKey)
	c, err := t.Client()
	if err != nil {
		return nil, err
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
		return nil, errors.New(fmt.Sprintf("Response did not contain valid JSON. Body: %v", body))
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
		return nil, errors.New(fmt.Sprintf("Response did not contain valid JSON. Body: %v", body))
	}

	return host, nil
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
		return errors.New(fmt.Sprintf("Response did not contain valid JSON. Body: %v", body))
	}

	return errors.New(fmt.Sprintf("Error from API. %v (%v)", jsonBody["reason"], jsonBody["detail"]))
}
