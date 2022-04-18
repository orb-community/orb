/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package types

import (
    "encoding/json"
    "github.com/ns1labs/orb/pkg/errors"
)

// Tags A flat kv pair object
type Tags map[string]string

// Metadata Maybe a full object hierarchy
type Metadata map[string]interface{}

func (s *Metadata) Scan(src interface{}) error {
   switch v := src.(type) {
   case []byte:
       return json.Unmarshal(v, s)
   case string:
       return json.Unmarshal([]byte(v), s)
   }
   return errors.New("type assertion failed")
}