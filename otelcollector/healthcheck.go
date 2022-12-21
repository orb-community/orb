package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hpcloud/tail"
)

type Answer struct {
	response map[string]string
}

func newHealthCheck(initStatus string, initMessage string) *Answer {
	resp := make(map[string]string)
	resp["status"] = initStatus
	resp["message"] = initMessage
	return &Answer{response: resp}
}

func (a *Answer) health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(a.response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (a *Answer) filterLogs(logFile string) {
	go func() {
		t, err := tail.TailFile(logFile, tail.Config{Follow: true, ReOpen: true})
		if err != nil {
			fmt.Println("Error reading log file")
		}
		for line := range t.Lines {
			fmt.Println(line.Text)
			if strings.Contains(line.Text, "Permanent error: remote write returned HTTP status 401 Unauthorized") {
				a.response["status"] = "error"
				a.response["message"] = "Permanent error: remote write returned HTTP status 401 Unauthorized"
			}
			if strings.Contains(line.Text, "Permanent error: remote write returned HTTP status 429 Too Many Requests") {
				a.response["status"] = "error"
				a.response["message"] = "Permanent error: remote write returned HTTP status 429 Too Many Requests"
			}
		}
	}()
}

func main() {
	// initial status
	as := newHealthCheck("ok", "healthy")
	// start filtering logs
	as.filterLogs("/tmp/otel-collector.log")
	http.HandleFunc("/health", as.health)
	// start webserver
	http.ListenAndServe(":8090", nil)
}
