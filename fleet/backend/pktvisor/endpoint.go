package pktvisor

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

func listAgentBackendsEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	//svc := isvc.(fleet.Service)
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listAgentBackendsReq)
		if err := req.validate(); err != nil {
			return agentBackendsRes{}, err
		}
		pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		backends := backend.GetList()
		if err != nil {
			return agentBackendsRes{}, err
		}
		var res []interface{}
		for _, be := range backends {
			mt := backend.GetBackend(be).Metadata()
			if err != nil {
				return agentBackendsRes{}, err
			}
			res = append(res, mt)
		}

		return agentBackendsRes{
			Backends: res,
		}, nil
	}
}

func viewAgentBackendHandlerEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	//svc := isvc.(fleet.Service)
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		if backend.HaveBackend(req.id) {
			bk, err := Handlers()
			if err != nil {
				return nil, err
			}
			return bk, nil
		}
		return Inputs()
		//return svc.ViewAgentBackendHandler(ctx, req.token, req.id)
	}
}

func viewAgentBackendInputEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	//svc := isvc.(fleet.Service)
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		if backend.HaveBackend(req.id) {
			bk, err := Inputs()
			if err != nil {
				return nil, err
			}
			return bk, nil
		}
		return Inputs()
		//return svc.ViewAgentBackendInput(ctx, req.token, req.id)
	}
}

func viewAgentBackendTapsEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	//svc := isvc.(fleet.Service)
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		r, err := pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		if !backend.HaveBackend("pktvisor") {
			return nil, errors.ErrNotFound
		}
		metadataList, err := RetrieveAgentMetadataByOwner(ctx, r.Id, pkt.db)
		if err != nil {
			return nil, err
		}

		var list []types.Metadata
		for _, mt := range metadataList {
			extractTaps(mt, &list)
		}

		res, err := toBackendTaps(list)
		if err != nil {
			return nil, err
		}
		tapsGroup := groupTaps(res)

		var tpRes []agentBackendTapsRes
		for _, v := range tapsGroup {
			tpRes = append(tpRes, agentBackendTapsRes{
				Name:             v.Name,
				InputType:        v.InputType,
				ConfigPredefined: v.ConfigPredefined,
				TotalAgents:      totalAgents{Total: v.TotalAgents},
			})
		}
		return tpRes, nil
	}
}

// Used to get the taps from policy json
func extractTaps(mt map[string]interface{}, list *[]types.Metadata) {
	for k, v := range mt {
		if k == "taps" {
			m, _ := v.(map[string]interface{})
			*list = append(*list, m)
		} else {
			m, _ := v.(map[string]interface{})
			extractTaps(m, list)
		}
	}
}

// Used to cast the map[string]interface for a concrete struct
func toBackendTaps(list []types.Metadata) ([]BackendTaps, error) {
	var bkTaps []BackendTaps
	for _, tc := range list {
		bkTap := BackendTaps{}
		var idx int
		for k, v := range tc {
			bkTap.Name = k
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("Error to group taps")
			}
			for k, v := range m {
				switch k {
				case "config":
					m, ok := v.(map[string]interface{})
					if !ok {
						return nil, errors.New("Error to group taps")
					}
					for k, _ := range m {
						bkTap.ConfigPredefined = append(bkTap.ConfigPredefined, []string{k}...)
					}
				case "input_type":
					bkTap.InputType = k
				}
			}
			idx++
			bkTap.TotalAgents += uint64(idx)
			bkTaps = append(bkTaps, bkTap)
		}
	}
	return bkTaps, nil
}

// Used to aggregate and sumarize the taps and return a slice of BackendTaps
func groupTaps(taps []BackendTaps) []BackendTaps {
	//TODO sort taps before group
	tapsMap := make(map[string]BackendTaps)
	for _, tap := range taps {
		key := key(tap.Name, tap.InputType)
		if v, ok := tapsMap[key]; ok {
			v.ConfigPredefined = append(v.ConfigPredefined, tap.ConfigPredefined...)
			v.TotalAgents += 1
			tapsMap[key] = v
		} else {
			tapsMap[key] = BackendTaps{
				Name:             tap.Name,
				InputType:        tap.InputType,
				ConfigPredefined: tap.ConfigPredefined,
				TotalAgents:      tap.TotalAgents,
			}
		}
	}
	var bkTaps []BackendTaps
	for _, v := range tapsMap {
		bkTaps = append(bkTaps, v)
	}
	return bkTaps
}

func key(name string, inputType string) string {
	return fmt.Sprintf("%s-%s", name, inputType)
}
