package kubecontrol

const CollectorStatusKey = "orb.sinks.collectors"

type CollectorStatusSortedSetEntry struct {
	SinkId       string  `json:"sinkId"`
	Status       string  `json:"status"`
	ErrorMessage *string `json:"error_message,omitempty"`
}
