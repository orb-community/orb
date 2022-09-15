package orbreceiver

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"reflect"
	"testing"
)

func generateMetricsRequest() pmetricotlp.Request {
	md := pmetric.NewMetrics()
	m := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	m.SetName("test_metric")
	m.SetDataType(pmetric.MetricDataTypeGauge)
	m.Gauge().DataPoints().AppendEmpty()
	return pmetricotlp.NewRequestFromMetrics(md)
}

func Test_jsonEncoder_unmarshalMetricsRequest(t *testing.T) {
	type args struct {
		metric pmetricotlp.Request
	}
	tests := []struct {
		name    string
		args    args
		want    pmetricotlp.Request
		wantErr bool
	}{
		{
			name: "going back and forth",
			args: args{
				metric: generateMetricsRequest(),
			},
			want:    generateMetricsRequest(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodingBuf, err := tt.args.metric.MarshalProto()
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalMetricsRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			e := protoEncoder{}
			got, err := e.unmarshalMetricsRequest(encodingBuf)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalMetricsRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalMetricsRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
