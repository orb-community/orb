package orbreceiver

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
)

const dataFormatProtobuf = "protobuf"

// Receiver is the type used to handle metrics from OpenTelemetry exporters.
type internalReceiver struct {
	pmetricotlp.UnimplementedGRPCServer
	nextConsumer consumer.Metrics
	obsrecv      *obsreport.Receiver
}

// New creates a new Receiver reference.
func InternalReceiverNew(nextConsumer consumer.Metrics, obsrecv *obsreport.Receiver) *internalReceiver {
	return &internalReceiver{
		nextConsumer: nextConsumer,
		obsrecv:      obsrecv,
	}
}

// Export implements the service Export metrics func.
func (r *internalReceiver) Export(ctx context.Context, req pmetricotlp.ExportRequest) (pmetricotlp.ExportResponse, error) {
	md := req.Metrics()
	dataPointCount := md.DataPointCount()
	if dataPointCount == 0 {
		return pmetricotlp.NewExportResponse(), nil
	}

	ctx = r.obsrecv.StartMetricsOp(ctx)
	err := r.nextConsumer.ConsumeMetrics(ctx, md)
	r.obsrecv.EndMetricsOp(ctx, dataFormatProtobuf, dataPointCount, err)

	return pmetricotlp.NewExportResponse(), err
}
