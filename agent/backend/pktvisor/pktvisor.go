/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"errors"
	"github.com/go-cmd/cmd"
	"github.com/ns1labs/orb/agent/backend"
	"go.uber.org/zap"
	"os/exec"
	"time"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger     *zap.Logger
	binary     string
	proc       *cmd.Cmd
	statusChan <-chan cmd.Status
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
	}, p.binary, "--admin-api")
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
		return err
	}

	if status.Complete {
		p.logger.Error("pktvisor startup error, check log", zap.Int("exit_code", status.Exit))
		return err
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

	return nil
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}
