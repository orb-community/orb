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

package kafkaexporter

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestValidate_err_compression(t *testing.T) {
	config := &Config{
		Producer: Producer{
			Compression: "idk",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "producer.compression should be one of 'none', 'gzip', 'snappy', 'lz4', or 'zstd'. configured value idk")
}

func Test_saramaProducerCompressionCodec(t *testing.T) {
	tests := map[string]struct {
		compression         string
		expectedCompression sarama.CompressionCodec
		expectedError       error
	}{
		"none": {
			compression:         "none",
			expectedCompression: sarama.CompressionNone,
			expectedError:       nil,
		},
		"gzip": {
			compression:         "gzip",
			expectedCompression: sarama.CompressionGZIP,
			expectedError:       nil,
		},
		"snappy": {
			compression:         "snappy",
			expectedCompression: sarama.CompressionSnappy,
			expectedError:       nil,
		},
		"lz4": {
			compression:         "lz4",
			expectedCompression: sarama.CompressionLZ4,
			expectedError:       nil,
		},
		"zstd": {
			compression:         "zstd",
			expectedCompression: sarama.CompressionZSTD,
			expectedError:       nil,
		},
		"unknown": {
			compression:         "unknown",
			expectedCompression: sarama.CompressionNone,
			expectedError:       fmt.Errorf("producer.compression should be one of 'none', 'gzip', 'snappy', 'lz4', or 'zstd'. configured value unknown"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := saramaProducerCompressionCodec(test.compression)
			assert.Equal(t, c, test.expectedCompression)
			assert.Equal(t, err, test.expectedError)
		})
	}
}
