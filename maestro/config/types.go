package config

import (
	"database/sql/driver"
	"time"
)

type SinkData struct {
	SinkID          string          `json:"sink_id"`
	OwnerID         string          `json:"owner_id"`
	Url             string          `json:"remote_host"`
	User            string          `json:"username"`
	Password        string          `json:"password"`
	Token           string          `json:"token"`
	OpenTelemetry   string          `json:"opentelemetry"`
	State           PrometheusState `json:"state,omitempty"`
	Migrate         string          `json:"migrate,omitempty"`
	Msg             string          `json:"msg,omitempty"`
	LastRemoteWrite time.Time       `json:"last_remote_write,omitempty"`
}

const (
	Unknown PrometheusState = iota
	Active
	Error
	Idle
	Warning
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
	"warning",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
	"warning": Warning,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) SetFromString(value string) error {
	*p = promStateRevMap[value]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}
