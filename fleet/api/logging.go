/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

var _ fleet.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    fleet.Service
}

func (l loggingMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (_ fleet.Page, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_agents",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_agents",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListAgents(ctx, token, pm)
}

func (l loggingMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (_ fleet.Agent, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_agent",
				zap.String("name", a.Name.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_agent",
				zap.String("name", a.Name.String()),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateAgent(ctx, token, a)
}

func (l loggingMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_agent_group",
				zap.String("name", s.Name.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_agent_group",
				zap.String("name", s.Name.String()),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateAgentGroup(ctx, token, s)
}

func NewLoggingMiddleware(svc fleet.Service, logger *zap.Logger) fleet.Service {
	return &loggingMiddleware{logger, svc}
}
