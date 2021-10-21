package pktvisor

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/mainflux"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
)

func viewAgentBackendHandlerEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		_, err = pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		bk, err := pkt.handlers()
		if err != nil {
			return nil, err
		}
		return bk, nil
	}
}

func viewAgentBackendInputEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		_, err = pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		bk, err := pkt.inputs()
		if err != nil {
			return nil, err
		}
		return bk, nil
	}
}

func viewAgentBackendTapsEndpoint(pkt pktvisorBackend) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(viewResourceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		r, err := pkt.auth.Identify(ctx, &mainflux.Token{Value: req.token})
		if err != nil {
			return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
		}

		metadataList, err := pkt.taps(ctx, r.Id)
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
					bkTap.InputType = fmt.Sprintf("%v",v)
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
	keys := make(map[string]bool)
	for _, v := range tapsMap {
		list := []string{}
		for _, config := range v.ConfigPredefined {
			if _, ok := keys[config]; !ok {
				keys[config] = true
				list = append(list, config)
			}
		}
		v.ConfigPredefined = list
		bkTaps = append(bkTaps, v)
	}
	return bkTaps
}

func key(name string, inputType string) string {
	return fmt.Sprintf("%s-%s", name, inputType)
}
