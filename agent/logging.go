/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type agentLoggerDebug struct {
	a *orbAgent
}
type agentLoggerWarn struct {
	a *orbAgent
}
type agentLoggerCritical struct {
	a *orbAgent
}
type agentLoggerError struct {
	a *orbAgent
}

var _ mqtt.Logger = (*agentLoggerDebug)(nil)
var _ mqtt.Logger = (*agentLoggerWarn)(nil)
var _ mqtt.Logger = (*agentLoggerCritical)(nil)
var _ mqtt.Logger = (*agentLoggerError)(nil)

func (a *agentLoggerWarn) Println(v ...interface{}) {
	a.a.logger.Warn("WARN mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerWarn) Printf(format string, v ...interface{}) {
	a.a.logger.Warn("WARN mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerDebug) Println(v ...interface{}) {
	a.a.logger.Debug("DEBUG mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerDebug) Printf(format string, v ...interface{}) {
	a.a.logger.Debug("DEBUG mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerCritical) Println(v ...interface{}) {
	a.a.logger.Error("CRITICAL mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerCritical) Printf(format string, v ...interface{}) {
	a.a.logger.Error("CRITICAL mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerError) Println(v ...interface{}) {
	a.a.logger.Error("ERROR mqtt log", zap.Any("payload", v))
}
func (a *agentLoggerError) Printf(format string, v ...interface{}) {
	a.a.logger.Error("ERROR mqtt log", zap.Any("payload", v))
}
