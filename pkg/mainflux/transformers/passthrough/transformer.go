/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package passthrough

import (
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/transformers"
)

type funcTransformer func(messaging.Message) (interface{}, error)

func New() transformers.Transformer {
	return funcTransformer(transformer)
}

func (fh funcTransformer) Transform(msg messaging.Message) (interface{}, error) {
	return fh(msg)
}

func transformer(msg messaging.Message) (interface{}, error) {
	return msg, nil
}
