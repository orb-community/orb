package cloudprober

import (
	"github.com/ns1labs/orb/agent/policies"
	"reflect"
	"testing"
)

func Test_cloudproberBackend_buildConfigFile(t *testing.T) {

	type args struct {
		policyYaml []policies.PolicyData
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cloudproberBackend{
				logger:     nil,
				configFile: "",
			}
			got, err := c.buildConfigFile(tt.args.policyYaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildConfigFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
