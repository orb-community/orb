package orbreceiver

import (
	"fmt"
	"os"
	"testing"
)

func Test_decodeTestDataFiles(t *testing.T) {
	dir := "/home/lpegoraro/workspace/orb/sinker/otel/orbreceiver/testdata"
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, dirEntry := range files {
		fileName := fmt.Sprintf("%s/%s", dir, dirEntry.Name())
		file, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatal(err)
		}
		e := protoEncoder{}
		got, err := e.unmarshalMetricsRequest(file)
		if err != nil {
			t.Errorf("unmarshalMetricsRequest() error = %v", err)
			return
		}
		//md := req.Metrics()
		//dataPointCount := md.DataPointCount()
		//if dataPointCount == 0 {
		//	return pmetricotlp.NewResponse(), nil
		//}
		//
		//ctx = r.obsrecv.StartMetricsOp(ctx)
		//err := r.nextConsumer.ConsumeMetrics(ctx, md)
		//r.obsrecv.EndMetricsOp(ctx, dataFormatProtobuf, dataPointCount, err)

		t.Log("succeeded reading and unmarshalling got: ", got, string(file))
	}
}
