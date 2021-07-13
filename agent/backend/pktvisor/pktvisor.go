/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/ns1labs/orb/agent/backend"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger     *zap.Logger
	binary     string
	configFile string
	proc       *cmd.Cmd
	statusChan <-chan cmd.Status

	adminAPIHost     string
	adminAPIPort     uint16
	adminAPIProtocol string
}

// AppMetrics represents server application information
type AppMetrics struct {
	App struct {
		Version   string  `json:"version"`
		UpTimeMin float64 `json:"up_time_min"`
	} `json:"app"`
}

func (p *pktvisorBackend) request(url string, payload interface{}, method string, body io.Reader, contentType string) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	alive, err := p.checkAlive()
	if !alive {
		return err
	}

	URL := fmt.Sprintf("%s://%s:%d/api/v1/%s", p.adminAPIProtocol, p.adminAPIHost, p.adminAPIPort, url)

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return err
	}
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Add("Content-Type", contentType)

	res, getErr := client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("non 200 HTTP error code from pktvisord, no or invalid body: %d", res.StatusCode))
		}
		if body[0] == '{' {
			var jsonBody map[string]interface{}
			err := json.Unmarshal(body, &jsonBody)
			if err == nil {
				if errMsg, ok := jsonBody["error"]; ok {
					return errors.New(fmt.Sprintf("%d %s", res.StatusCode, errMsg))
				}
			}
		}
		return errors.New(fmt.Sprintf("%d %s", res.StatusCode, body))
	}

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return err
	}
	return nil
}

func (p *pktvisorBackend) checkAlive() (bool, error) {
	status := p.proc.Status()

	if status.Error != nil {
		p.logger.Error("pktvisor process error", zap.Error(status.Error))
		return false, status.Error
	}

	if status.Complete {
		p.proc.Stop()
		// TODO auto restart
		return false, errors.New("pktvisor process ended")
	}

	return true, nil
}

func (p *pktvisorBackend) ApplyPolicy(policy interface{}) error {

	p.logger.Info("pktvisor policy", zap.Any("data", policy))

	pyaml, err := yaml.Marshal(policy)
	if err != nil {
		return err
	}

	var resp map[string]interface{}
	err = p.request("policies", &resp, http.MethodPost, bytes.NewBuffer(pyaml), "application/x-yaml")
	if err != nil {
		p.logger.Debug("yaml policy failure", zap.ByteString("policy", pyaml))
		return err
	}

	return nil

}

func (p *pktvisorBackend) Version() (string, error) {
	var appMetrics AppMetrics
	err := p.request("metrics/app", &appMetrics, http.MethodGet, nil, "")
	if err != nil {
		return "", err
	}
	return appMetrics.App.Version, nil
}

func (p *pktvisorBackend) Write(payload []byte) (n int, err error) {
	p.logger.Info("pktvisord", zap.ByteString("log", payload))
	return len(payload), nil
}

func (p *pktvisorBackend) Start() error {
	p.logger.Info("pktvisor starting")

	_, err := exec.LookPath(p.binary)
	if err != nil {
		p.logger.Error("pktvisor startup error: binary not found", zap.Error(err))
		return err
	}

	p.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, p.binary, "--admin-api", "--config", p.configFile)
	p.statusChan = p.proc.Start()

	// log STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for p.proc.Stdout != nil || p.proc.Stderr != nil {
			select {
			case line, open := <-p.proc.Stdout:
				if !open {
					p.proc.Stdout = nil
					continue
				}
				p.logger.Info("pktvisor stdout", zap.String("log", line))
			case line, open := <-p.proc.Stderr:
				if !open {
					p.proc.Stderr = nil
					continue
				}
				p.logger.Info("pktvisor stderr", zap.String("log", line))
			}
		}
	}()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := p.proc.Status()

	if status.Error != nil {
		p.logger.Error("pktvisor startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		p.proc.Stop()
		return errors.New("pktvisor startup error, check log")
	}

	p.logger.Info("pktvisor process started", zap.Int("pid", status.PID))

	return nil
}

func (p *pktvisorBackend) Stop() error {
	p.logger.Info("pktvisor stopping")
	err := p.proc.Stop()
	finalStatus := <-p.statusChan
	if err != nil {
		p.logger.Error("pktvisor shutdown error", zap.Error(err))
		return err
	}
	p.logger.Info("pktvisor process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (p *pktvisorBackend) Configure(logger *zap.Logger, config map[string]string) error {
	p.logger = logger

	var prs bool
	if p.binary, prs = config["binary"]; !prs {
		return errors.New("you must specify pktvisor binary")
	}
	if p.configFile, prs = config["config_file"]; !prs {
		return errors.New("you must specify pktvisor configuration file")
	}
	return nil
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{
		adminAPIHost:     "localhost",
		adminAPIPort:     10853,
		adminAPIProtocol: "http",
	})
	return true
}
