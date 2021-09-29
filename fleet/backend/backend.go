/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/opentracing/opentracing-go"
)

type Backend interface {
	Metadata() interface{}
	MakeHandler(tracer opentracing.Tracer, opts []kithttp.ServerOption, r *bone.Mux)
}

type ConfigFeature struct {
	Type     string `json:"type"`
	Input    string `json:"input"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

var registry = make(map[string]Backend)

func Register(name string, b Backend) {
	registry[name] = b
}

func GetList() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

func HaveBackend(name string) bool {
	_, prs := registry[name]
	return prs
}

func GetBackend(name string) Backend {
	return registry[name]
}
