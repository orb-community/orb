/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cloud_config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type CloudConfigManager interface {
	GetCloudConfig() (config.MQTTConfig, error)
}

var _ CloudConfigManager = (*cloudConfigManager)(nil)

type cloudConfigManager struct {
	logger *zap.Logger
	config config.Config
	db     *sqlx.DB
}

func New(logger *zap.Logger, c config.Config, db *sqlx.DB) (CloudConfigManager, error) {
	return &cloudConfigManager{logger: logger, config: c, db: db}, nil
}

func (cc *cloudConfigManager) migrateDB() error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "cloud_config_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS cloud_config (
						address TEXT NOT NULL,
						id TEXT	NOT NULL,
						key TEXT	NOT NULL,
						channel TEXT	NOT NULL,
						ts_created INTEGER NOT NULL
						)`,
				},
				Down: []string{
					"DROP TABLE cloud_config",
				},
			},
		},
	}

	_, err := migrate.Exec(cc.db.DB, "sqlite3", migrations, migrate.Up)

	return err
}

func (cc *cloudConfigManager) request(address string, token string, response interface{}, method string, body []byte) error {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	URL := fmt.Sprintf("%s/api/v1/agents", address)

	req, err := http.NewRequest(method, URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	cc.logger.Debug("cloud api request", zap.String("url", req.URL.String()), zap.ByteString("body", body))
	req.Header.Add("Authorization", token)

	res, getErr := client.Do(req)
	if getErr != nil {
		return getErr
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("non 200 HTTP error code from API, no or invalid body: %d", res.StatusCode))
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

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}

func (cc *cloudConfigManager) autoProvision(apiAddress string, token string) (config.MQTTConfig, error) {

	type AgentRes struct {
		ID        string `json:"id"`
		Key       string `json:"key"`
		ChannelID string `json:"channel_id"`
	}

	type AgentReq struct {
		Name string `json:"name"`
	}

	aname, haveAname := cc.config.OrbAgent.Cloud["config"]["agent_name"]
	if !haveAname {
		hostname, err := os.Hostname()
		if err != nil {
			return config.MQTTConfig{}, err
		}
		aname = hostname
	}

	agentReq := AgentReq{Name: strings.Replace(aname, ".", "-", -1)}
	body, err := json.Marshal(agentReq)
	if err != nil {
		return config.MQTTConfig{}, err
	}

	cc.logger.Info("attempting auto provision", zap.String("address", apiAddress))

	var result AgentRes
	err = cc.request(apiAddress, token, result, http.MethodPost, body)
	if err != nil {
		return config.MQTTConfig{}, err
	}

	return config.MQTTConfig{}, errors.New("unable to auto provision agent")
}

func (cc *cloudConfigManager) GetCloudConfig() (config.MQTTConfig, error) {

	// if MQTT is specified in the config file, always use that
	mqttAddress, haveMqttAddress := cc.config.OrbAgent.Cloud["mqtt"]["address"]
	apiAddress, haveApiAddress := cc.config.OrbAgent.Cloud["api"]["address"]
	id, haveId := cc.config.OrbAgent.Cloud["mqtt"]["id"]
	key, haveKey := cc.config.OrbAgent.Cloud["mqtt"]["key"]
	channel, haveChannel := cc.config.OrbAgent.Cloud["mqtt"]["channel_id"]

	// currently we require address to be specified, it cannot be auto provisioned.
	// this may change in the future
	if !haveMqttAddress {
		return config.MQTTConfig{}, errors.New("cloud.mqtt.address must be specified in configuration")
	}

	if haveMqttAddress && haveId && haveKey && haveChannel {
		cc.logger.Info("using explicitly specified cloud configuration",
			zap.String("address", mqttAddress),
			zap.String("id", id))
		return config.MQTTConfig{
			Address:   mqttAddress,
			Id:        id,
			Key:       key,
			ChannelID: channel,
		}, nil
	}

	// if full config is not available, possibly attempt auto provision
	var ap bool
	var err error
	apConfig, haveAutoProvision := cc.config.OrbAgent.Cloud["config"]["auto_provision"]
	if haveAutoProvision {
		ap, err = strconv.ParseBool(apConfig)
		if err != nil {
			return config.MQTTConfig{}, err
		}
	} else {
		ap = true
	}

	if !ap {
		return config.MQTTConfig{}, errors.New("valid cloud MQTT config was not specified, and auto_provision was disabled")
	}
	if !haveApiAddress {
		return config.MQTTConfig{}, errors.New("wanted to auto provision, but no API address was available")
	}
	token, haveToken := cc.config.OrbAgent.Cloud["api"]["token"]
	if !haveToken {
		return config.MQTTConfig{}, errors.New("wanted to auto provision, but no API token was available")
	}

	err = cc.migrateDB()
	if err != nil {
		return config.MQTTConfig{}, err
	}
	return cc.autoProvision(apiAddress, token)

}
