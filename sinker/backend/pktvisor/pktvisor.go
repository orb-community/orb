/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/prometheus"
	"go.uber.org/zap"
	stdlog "log"
	"os"
	"strconv"
	"strings"
	"time"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	promClient prometheus.Client
	logger *zap.Logger
}

type labelList []prometheus.Label
type headerList []header
type dp prometheus.Datapoint

type header struct {
	name  string
	value string
}

func (p pktvisorBackend) ProcessMetrics(thingID string, channelID string, subtopic []string, payload []fleet.AgentMetricsRPCPayload) error {
	// process batch
	for _, data := range payload {
		// TODO check pktvisor version in data.BEVersion against PktvisorVersion
		// TODO policyID and datasetIDs are in data
		if data.Format != "json" {
			p.logger.Warn("ignoring non-json pktvisor payload", zap.String("format", data.Format))
			continue
		}
		// unmarshal pktvisor metrics
		var metrics map[string]map[string]interface{}
		err := json.Unmarshal(data.Data, &metrics)
		if err != nil {
			p.logger.Warn("unable to unmarshal pktvisor metric payload", zap.Any("payload", data.Data))
			continue
		}
		stats := StatSnapshot{}
		for _, handlerData := range metrics {
			if data, ok := handlerData["pcap"]; ok {
				err := mapstructure.Decode(data, &stats.Pcap)
				if err != nil {
					p.logger.Error("error decoding pcap handler", zap.Error(err))
					continue
				}
			} else if data, ok := handlerData["dns"]; ok {
				err := mapstructure.Decode(data, &stats.DNS)
				if err != nil {
					p.logger.Error("error decoding dns handler", zap.Error(err))
					continue
				}
			} else if data, ok := handlerData["packets"]; ok {
				err := mapstructure.Decode(data, &stats.Packets)
				if err != nil {
					p.logger.Error("error decoding packets handler", zap.Error(err))
					continue
				}
			}
		}
		// TODO turn StatSnapshot into format for prometheus remote_write
		test()
	}
	return nil
}

func Register(logger *zap.Logger, promClient prometheus.Client) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger, promClient: promClient})
	return true
}

func test() {
	var (
		log            = stdlog.New(os.Stderr, "promremotecli_log ", stdlog.LstdFlags)
		writeURLFlag   string
		labelsListFlag labelList
		headerListFlag headerList
		dpFlag         dp
	)

	//flag.StringVar(&writeURLFlag, "u", prometheus.DefaultRemoteWrite, "remote write endpoint")
	//flag.Var(&labelsListFlag, "t", "label pair to include in metric. specify as key:value e.g. status_code:200")
	//flag.Var(&headerListFlag, "h", "headers to set in the request, e.g. 'User-Agent: foo'")
	//flag.Var(&dpFlag, "d", "datapoint to add. specify as unixTimestamp(int),value(float) e.g. 1556026059,14.23. use `now` instead of timestamp for current time")
	//
	//flag.Parse()

	cfg := prometheus.NewConfig(
		prometheus.WriteURLOption(prometheus.DefaultRemoteWrite),
	)

	client, err := prometheus.NewClient(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to construct client: %v", err))
	}

	var headers map[string]string
	log.Println("writing datapoint", dpFlag.String())
	log.Println("labelled", labelsListFlag.String())
	if len(headerListFlag) > 0 {
		log.Println("with headers", headerListFlag.String())
		headers = make(map[string]string, len(headerListFlag))
		for _, header := range headerListFlag {
			headers[header.name] = header.value
		}
	}
	log.Println("writing to", writeURLFlag)

	timeSeriesList := []prometheus.TimeSeries{
		prometheus.TimeSeries{
			Labels: []prometheus.Label{
				{
					Name:  "__name__",
					Value: "foo_bar",
				},
				{
					Name:  "biz",
					Value: "baz",
				},
				{
					Name:  "pessoa",
					Value: "daniel",
				},
			},
			Datapoint: prometheus.Datapoint{
				Timestamp: time.Now(),
				Value:     1415.92,
			},
		},
	}

	result, writeErr := client.WriteTimeSeries(context.Background(), timeSeriesList,
		prometheus.WriteOptions{Headers: headers})
	if err := error(writeErr); err != nil {
		json.NewEncoder(os.Stdout).Encode(struct {
			Success    bool   `json:"success"`
			Error      string `json:"error"`
			StatusCode int    `json:"statusCode"`
		}{
			Success:    false,
			Error:      err.Error(),
			StatusCode: writeErr.StatusCode(),
		})
		os.Stdout.Sync()

		log.Fatal("write error", err)
	}

	json.NewEncoder(os.Stdout).Encode(struct {
		Success    bool `json:"success"`
		StatusCode int  `json:"statusCode"`
	}{
		Success:    true,
		StatusCode: result.StatusCode,
	})
	os.Stdout.Sync()

	log.Println("write success")
}

func (t *labelList) String() string {
	var labels [][]string
	for _, v := range []prometheus.Label(*t) {
		labels = append(labels, []string{v.Name, v.Value})
	}
	return fmt.Sprintf("%v", labels)
}

func (t *labelList) Set(value string) error {
	labelPair := strings.Split(value, ":")
	if len(labelPair) != 2 {
		return fmt.Errorf("incorrect number of arguments to '-t': %d", len(labelPair))
	}

	label := prometheus.Label{
		Name:  labelPair[0],
		Value: labelPair[1],
	}

	*t = append(*t, label)

	return nil
}

func (h *headerList) String() string {
	var headers [][]string
	for _, v := range []header(*h) {
		headers = append(headers, []string{v.name, v.value})
	}
	return fmt.Sprintf("%v", headers)
}

func (h *headerList) Set(value string) error {
	firstSplit := strings.Index(value, ":")
	if firstSplit == -1 {
		return fmt.Errorf("header missing separating colon: '%v'", value)
	}

	*h = append(*h, header{
		name:  strings.TrimSpace(value[:firstSplit]),
		value: strings.TrimSpace(value[firstSplit+1:]),
	})

	return nil
}

func (d *dp) String() string {
	return fmt.Sprintf("%v", []string{d.Timestamp.String(), fmt.Sprintf("%v", d.Value)})
}

func (d *dp) Set(value string) error {
	dp := strings.Split(value, ",")
	if len(dp) != 2 {
		return fmt.Errorf("incorrect number of arguments to '-d': %d", len(dp))
	}

	var ts time.Time
	if strings.ToLower(dp[0]) == "now" {
		ts = time.Now()
	} else {
		i, err := strconv.Atoi(dp[0])
		if err != nil {
			return fmt.Errorf("unable to parse timestamp: %s", dp[1])
		}
		ts = time.Unix(int64(i), 0)
	}

	val, err := strconv.ParseFloat(dp[1], 64)
	if err != nil {
		return fmt.Errorf("unable to parse value as float64: %s", dp[0])
	}

	d.Timestamp = ts
	d.Value = val

	return nil
}
