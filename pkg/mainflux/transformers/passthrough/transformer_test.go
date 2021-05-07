/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package passthrough_test

import (
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestPassThrough(t *testing.T) {
	msg := messaging.Message{
		Channel:   "channel-1",
		Subtopic:  "subtopic-1",
		Publisher: "publisher-1",
		Protocol:  "protocol",
		Payload:   []byte(`some payload`),
		Created:   time.Now().Unix(),
	}
	tr := passthrough.New()
	m, err := tr.Transform(msg)
	assert.Equal(t, err, nil, "unexpected error")
	assert.Equal(t, msg.Channel, m.(messaging.Message).Channel, "passthrough failed")
}
