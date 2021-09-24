/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/db"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/opentracing/opentracing-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	auth        mainflux.AuthServiceClient
	db          *sqlx.DB
	Backend     string
	Description string
}

type BackendTaps struct {
	Name             string
	InputType        string
	ConfigPredefined []string
	TotalAgents      uint64
}

func (p pktvisorBackend) Metadata() interface{} {
	return p
}

func (p pktvisorBackend) MakeHandler(tracer opentracing.Tracer, opts []kithttp.ServerOption, r *bone.Mux) {
	MakePktvisorHandler(tracer, p, opts, r)
}

func (p pktvisorBackend) handlers() (_ types.Metadata, err error) {

	wd := getWorkDirectory()
	jsonFile, err := ioutil.ReadFile(fmt.Sprintf("%s/fleet/backend/pktvisor/handlers.json", wd))
	if err != nil {
		return nil, err
	}
	var handlers types.Metadata
	err = json.Unmarshal([]byte(jsonFile), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) inputs() (_ types.Metadata, err error) {
	wd := getWorkDirectory()
	jsonFile, err := ioutil.ReadFile(fmt.Sprintf("%s/fleet/backend/pktvisor/inputs.json", wd))
	if err != nil {
		return nil, err
	}
	var handlers types.Metadata
	err = json.Unmarshal([]byte(jsonFile), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) retrieveAgentMetadataByOwner(ctx context.Context, ownerID string, db *sqlx.DB) ([]types.Metadata, error) {
	q := `SELECT agent_metadata
		FROM agents
		WHERE mf_owner_id = :mf_owner_id;`

	params := map[string]interface{}{
		"mf_owner_id": ownerID,
	}

	rows, err := db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSelectEntity, err)
	}
	defer rows.Close()

	var items []types.Metadata
	for rows.Next() {
		dbmd := dbAgent{}
		if err := rows.StructScan(&dbmd); err != nil {
			return nil, errors.Wrap(errors.ErrSelectEntity, err)
		}
		items = append(items, types.Metadata(dbmd.AgentMetadata))
	}
	return items, nil
}

type dbAgent struct {
	AgentMetadata db.Metadata `db:"agent_metadata"`
}

func getWorkDirectory() string {
	// When you works with tests, the path it's different from the prod running
	// So here I'm getting the right working directory, no matter if its test or prod
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "orb") {
		wd = filepath.Dir(wd)
	}
	return wd
}

func Register(auth mainflux.AuthServiceClient, db *sqlx.DB) bool {
	backend.Register("pktvisor", &pktvisorBackend{
		Backend:     "pktvisor",
		Description: "pktvisor observability agent from pktvisor.dev",
		auth:        auth,
		db:          db,
	})
	return true
}
