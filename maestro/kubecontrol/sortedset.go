package kubecontrol

const CollectorStatusKey = "orb.sinks.collectors"

type CollectorStatusSortedSetEntry struct {
	OwnerID      string  `json:"ownerID"`
	SinkID       string  `json:"sinkID"`
	Status       string  `json:"status"`
	ErrorMessage *string `json:"error_message,omitempty"`
}
