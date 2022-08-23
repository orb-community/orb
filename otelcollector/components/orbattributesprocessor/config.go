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

package orbattributesprocessor

import (
	"go.opentelemetry.io/collector/config"

	"github.com/ns1labs/orb/otelcollector/components/internal/attraction"
	"github.com/ns1labs/orb/otelcollector/components/internal/filterconfig"
)

// Config specifies the set of attributes to be inserted, updated, upserted and
// deleted and the properties to include/exclude a span from being processed.
// This processor handles all forms of modifications to attributes within a span, log, or metric.
// Prior to any actions being applied, each span is compared against
// the include properties and then the exclude properties if they are specified.
// This determines if a span is to be processed or not.
// The list of actions is applied in order specified in the configuration.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	filterconfig.MatchConfig `mapstructure:",squash"`

	// Specifies the list of attributes to act on.
	// The set of actions are {INSERT, UPDATE, UPSERT, DELETE, HASH, EXTRACT}.
	// This is a required field.
	attraction.Settings `mapstructure:",squash"`
}

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}

func (cfg *Config) appendAction(newAction attraction.ActionKeyValue) []attraction.ActionKeyValue {
	if cfg.Actions == nil {
		cfg.Actions = []attraction.ActionKeyValue{}
	}
	return append(cfg.Actions, newAction)
}

func (cfg *Config) keyValueEntry(key string, value interface{}, action attraction.Action) bool {
	newAction := attraction.ActionKeyValue{
		Key:    key,
		Value:  value,
		Action: action,
	}
	cfg.Actions = cfg.appendAction(newAction)
	return true
}

func (cfg *Config) keyFromContextEntry(key string, fromContext string, action attraction.Action) bool {
	newAction := attraction.ActionKeyValue{
		Key:         key,
		FromContext: fromContext,
		Action:      action,
	}
	cfg.Actions = cfg.appendAction(newAction)
	return true
}

func (cfg *Config) keyFromAttributeEntry(key string, fromAttribute string, action attraction.Action) bool {
	newAction := attraction.ActionKeyValue{
		Key:           key,
		FromAttribute: fromAttribute,
		Action:        action,
	}
	cfg.Actions = cfg.appendAction(newAction)
	return true
}

func (cfg *Config) AddInsertActionKeyValue(key string, value interface{}) (ok bool) {
	return cfg.keyValueEntry(key, value, attraction.INSERT)
}

func (cfg *Config) AddInsertActionFromContext(key, fromContext string) (ok bool) {
	return cfg.keyFromContextEntry(key, fromContext, attraction.INSERT)
}

func (cfg *Config) AddInsertActionFromAttribute(key, fromAttribute string) (ok bool) {
	return cfg.keyFromAttributeEntry(key, fromAttribute, attraction.INSERT)
}

func (cfg *Config) AddUpsertActionKeyValue(key string, value interface{}) (ok bool) {
	return cfg.keyValueEntry(key, value, attraction.UPSERT)
}

func (cfg *Config) AddUpsertActionFromAttribute(key, fromAttribute string) (ok bool) {
	return cfg.keyFromAttributeEntry(key, fromAttribute, attraction.UPSERT)
}

func (cfg *Config) AddUpsertActionFromContext(key, fromContext string) (ok bool) {
	return cfg.keyFromContextEntry(key, fromContext, attraction.UPSERT)
}

func (cfg *Config) AddUpdateActionKeyValue(key string, value interface{}) (ok bool) {
	return cfg.keyValueEntry(key, value, attraction.UPDATE)
}

func (cfg *Config) AddUpdateActionFromAttribute(key, fromAttribute string) (ok bool) {
	return cfg.keyFromAttributeEntry(key, fromAttribute, attraction.UPDATE)
}

func (cfg *Config) AddUpdateActionFromContext(key, fromContext string) (ok bool) {
	return cfg.keyFromContextEntry(key, fromContext, attraction.UPDATE)
}

func (cfg *Config) AddDeleteActionKey(key string) (ok bool) {
	newAction := attraction.ActionKeyValue{
		Key:    key,
		Action: attraction.DELETE,
	}
	cfg.Actions = cfg.appendAction(newAction)

	return true
}

func (cfg *Config) AddConvertAction(key, convertedType string) (ok bool) {
	newAction := attraction.ActionKeyValue{
		Key:           key,
		ConvertedType: convertedType,
		Action:        attraction.CONVERT,
	}
	cfg.Actions = cfg.appendAction(newAction)
	return true
}

func (cfg *Config) AddHashAction(key string) (ok bool) {
	newAction := attraction.ActionKeyValue{
		Key:    key,
		Action: attraction.HASH,
	}
	cfg.Actions = cfg.appendAction(newAction)
	return true
}
