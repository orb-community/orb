// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpexporter_test

import (
	"context"
	"github.com/ns1labs/orb/agent/otel/otlpexporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"

	"go.opentelemetry.io/collector/model/otlpgrpc"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var (
	resourceAttributes1      = pcommon.NewResource()
	TestMetricStartTime      = time.Date(2020, 2, 11, 20, 26, 12, 321, time.UTC)
	TestMetricStartTimestamp = pcommon.NewTimestampFromTime(TestMetricStartTime)

	TestMetricExemplarTime      = time.Date(2020, 2, 11, 20, 26, 13, 123, time.UTC)
	TestMetricExemplarTimestamp = pcommon.NewTimestampFromTime(TestMetricExemplarTime)

	TestMetricTime      = time.Date(2020, 2, 11, 20, 26, 13, 789, time.UTC)
	TestMetricTimestamp = pcommon.NewTimestampFromTime(TestMetricTime)
)

const (
	TestSumIntMetricName = "sum-int"
	TestLabelKey1        = "label-1"
	TestLabelValue1      = "label-value-1"
	TestLabelKey2        = "label-2"
	TestLabelValue2      = "label-value-2"
)

type mockReceiver struct {
	srv          *grpc.Server
	requestCount int32
	totalItems   int32
	mux          sync.Mutex
	metadata     metadata.MD
}

func (r *mockReceiver) GetMetadata() metadata.MD {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.metadata
}

type mockMetricsReceiver struct {
	mockReceiver
	lastRequest pmetric.Metrics
}

func (r *mockMetricsReceiver) Export(ctx context.Context, req otlpgrpc.MetricsRequest) (otlpgrpc.MetricsResponse, error) {
	md := req.Metrics()
	atomic.AddInt32(&r.requestCount, 1)
	atomic.AddInt32(&r.totalItems, int32(md.DataPointCount()))
	r.mux.Lock()
	defer r.mux.Unlock()
	r.lastRequest = md
	r.metadata, _ = metadata.FromIncomingContext(ctx)
	return otlpgrpc.NewMetricsResponse(), nil
}

func (r *mockMetricsReceiver) GetLastRequest() pmetric.Metrics {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.lastRequest
}

func otlpMetricsReceiverOnGRPCServer(ln net.Listener) *mockMetricsReceiver {
	rcv := &mockMetricsReceiver{
		mockReceiver: mockReceiver{
			srv: grpc.NewServer(),
		},
	}

	// Now run it as a gRPC server
	otlpgrpc.RegisterMetricsServer(rcv.srv, rcv)
	go func() {
		_ = rcv.srv.Serve(ln)
	}()

	return rcv
}

func TestSendMetrics(t *testing.T) {
	// Start an OTLP-compatible receiver.
	ln, err := net.Listen("tcp", "localhost:")
	require.NoError(t, err, "Failed to find an available address to run the gRPC server: %v", err)
	rcv := otlpMetricsReceiverOnGRPCServer(ln)
	// Also closes the connection.
	defer rcv.srv.GracefulStop()

	// Start an OTLP exporter and point to the receiver.
	factory := otlpexporter.NewFactory()
	cfg := factory.CreateDefaultConfig().(*otlpexporter.Config)
	cfg.GRPCClientSettings = configgrpc.GRPCClientSettings{
		Endpoint: ln.Addr().String(),
		TLSSetting: configtls.TLSClientSetting{
			Insecure: true,
		},
		Headers: map[string]string{
			"header": "header-value",
		},
	}
	set := componenttest.NewNopExporterCreateSettings()
	set.BuildInfo.Description = "Collector"
	set.BuildInfo.Version = "1.2.3test"
	exp, err := factory.CreateMetricsExporter(context.Background(), set, cfg)
	require.NoError(t, err)
	require.NotNil(t, exp)
	defer func() {
		assert.NoError(t, exp.Shutdown(context.Background()))
	}()

	host := componenttest.NewNopHost()

	assert.NoError(t, exp.Start(context.Background(), host))

	// Ensure that initially there is no data in the receiver.
	assert.EqualValues(t, 0, atomic.LoadInt32(&rcv.requestCount))

	// Send empty metric.
	md := pmetric.NewMetrics()
	assert.NoError(t, exp.ConsumeMetrics(context.Background(), md))

	// Wait until it is received.
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&rcv.requestCount) > 0
	}, 10*time.Second, 5*time.Millisecond)

	// Ensure it was received empty.
	assert.EqualValues(t, 0, atomic.LoadInt32(&rcv.totalItems))

	// Send two metrics.
	md = GenerateMetricsTwoMetrics()

	err = exp.ConsumeMetrics(context.Background(), md)
	assert.NoError(t, err)

	// Wait until it is received.
	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&rcv.requestCount) > 1
	}, 10*time.Second, 5*time.Millisecond)

	expectedHeader := []string{"header-value"}

	// Verify received metrics.
	assert.EqualValues(t, 2, atomic.LoadInt32(&rcv.requestCount))
	assert.EqualValues(t, 4, atomic.LoadInt32(&rcv.totalItems))
	assert.EqualValues(t, md, rcv.GetLastRequest())

	mdata := rcv.GetMetadata()
	require.EqualValues(t, mdata.Get("header"), expectedHeader)
	require.Equal(t, len(mdata.Get("User-Agent")), 1)
	require.Contains(t, mdata.Get("User-Agent")[0], "Collector/1.2.3test")
}

func GenerateMetricsTwoMetrics() pmetric.Metrics {
	md := GenerateMetricsOneEmptyInstrumentationLibrary()
	rm0ils0 := md.ResourceMetrics().At(0).ScopeMetrics().At(0)
	initSumIntMetric(rm0ils0.Metrics().AppendEmpty())
	initSumIntMetric(rm0ils0.Metrics().AppendEmpty())
	return md
}

func GenerateMetricsOneEmptyInstrumentationLibrary() pmetric.Metrics {
	md := GenerateMetricsNoLibraries()
	md.ResourceMetrics().At(0).ScopeMetrics().AppendEmpty()
	return md
}

func GenerateMetricsNoLibraries() pmetric.Metrics {
	md := GenerateMetricsOneEmptyResourceMetrics()
	ms0 := md.ResourceMetrics().At(0)
	initResource1(ms0.Resource())
	return md
}

func GenerateMetricsOneEmptyResourceMetrics() pmetric.Metrics {
	md := pmetric.NewMetrics()
	md.ResourceMetrics().AppendEmpty()
	return md
}

func initResource1(r pcommon.Resource) {
	initResourceAttributes1(r)
}

func initResourceAttributes1(dest pcommon.Resource) {
	dest.Attributes().Clear()
	resourceAttributes1.CopyTo(dest)
}

func initSumIntMetric(im pmetric.Metric) {
	initMetric(im, TestSumIntMetricName, pmetric.MetricDataTypeSum)

	idps := im.Sum().DataPoints()
	idp0 := idps.AppendEmpty()
	initMetricAttributes1(idp0.Attributes())
	idp0.SetStartTimestamp(TestMetricStartTimestamp)
	idp0.SetTimestamp(TestMetricTimestamp)
	idp0.SetIntVal(123)
	idp1 := idps.AppendEmpty()
	initMetricAttributes2(idp1.Attributes())
	idp1.SetStartTimestamp(TestMetricStartTimestamp)
	idp1.SetTimestamp(TestMetricTimestamp)
	idp1.SetIntVal(456)
}

func initMetric(m pmetric.Metric, name string, ty pmetric.MetricDataType) {
	m.SetName(name)
	m.SetDescription("")
	m.SetUnit("1")
	m.SetDataType(ty)
	switch ty {
	case pmetric.MetricDataTypeSum:
		sum := m.Sum()
		sum.SetIsMonotonic(true)
		sum.SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	case pmetric.MetricDataTypeHistogram:
		histo := m.Histogram()
		histo.SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	case pmetric.MetricDataTypeExponentialHistogram:
		histo := m.ExponentialHistogram()
		histo.SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
	}
}

func initMetricAttributes1(dest pcommon.Map) {
	dest.Clear()
	dest.InsertString(TestLabelKey1, TestLabelValue1)
}

func initMetricAttributes2(dest pcommon.Map) {
	dest.Clear()
	dest.InsertString(TestLabelKey2, TestLabelValue2)
}
