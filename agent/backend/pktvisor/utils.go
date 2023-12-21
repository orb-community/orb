/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/orb-community/orb/agent/backend"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

// note this needs to be stateless because it is called for multiple go routines
func (p *pktvisorBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string, timeout int32) error {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	status, _, err := p.getProcRunningStatus()
	if status != backend.Running {
		p.logger.Warn("skipping pktvisor REST API request because process is not running or is unresponsive", zap.String("url", url), zap.String("method", method), zap.Error(err))
		return err
	}

	URL := fmt.Sprintf("%s://%s:%s/api/v1/%s", p.adminAPIProtocol, p.adminAPIHost, p.adminAPIPort, url)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		p.logger.Error("received error from payload", zap.Error(err))
		return err
	}

	req.Header.Add("Content-Type", contentType)
	res, getErr := client.Do(req)

	if getErr != nil {
		p.logger.Error("received error from payload", zap.Error(getErr))
		return getErr
	}

	if (res.StatusCode < 200) || (res.StatusCode > 299) {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("non 2xx HTTP error code from pktvisord, no or invalid body: %d", res.StatusCode))
		}
		if len(body) == 0 {
			return errors.New(fmt.Sprintf("%d empty body", res.StatusCode))
		} else if body[0] == '{' {
			var jsonBody map[string]interface{}
			err := json.Unmarshal(body, &jsonBody)
			if err == nil {
				if errMsg, ok := jsonBody["error"]; ok {
					return errors.New(fmt.Sprintf("%d %s", res.StatusCode, errMsg))
				}
			}
		}
	}

	if res.Body != nil {
		err = json.NewDecoder(res.Body).Decode(&payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pktvisorBackend) getProcRunningStatus() (backend.RunningStatus, string, error) {
	if p.proc == nil {
		return backend.Unknown, "backend not started yet", nil
	}
	status := p.proc.Status()

	if status.Error != nil {
		errMsg := fmt.Sprintf("pktvisor process error: %v", status.Error)
		return backend.BackendError, errMsg, status.Error
	}

	if status.Complete {
		err := p.proc.Stop()
		return backend.Offline, "pktvisor process ended", err
	}

	if status.StopTs > 0 {
		return backend.Offline, "pktvisor process ended", nil
	}
	return backend.Running, "", nil
}

// also used for HTTP REST API readiness check
func (p *pktvisorBackend) getAppInfo() (AppInfo, error) {
	var appInfo AppInfo
	err := p.request("metrics/app", &appInfo, http.MethodGet, http.NoBody, "application/json", VersionTimeout)
	return appInfo, err
}
