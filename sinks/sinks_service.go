// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinks

import (
	"context"
	"encoding/hex"

	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/authentication_type"
	"github.com/orb-community/orb/sinks/backend"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	ErrCreateSink                 = errors.New("failed to create Sink")
	ErrConflictSink               = errors.New("entity already exists")
	ErrUnsupportedContentTypeSink = errors.New("unsupported content type")
	ErrValidateSink               = errors.New("failed to validate Sink")
)

func (svc sinkService) CreateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}

	sink.MFOwnerID = mfOwnerID

	be, err := svc.validateBackend(&sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	at, err := validateAuthType(&sink)
	if err != nil {
		return Sink{}, err
	}
	cfg := Configuration{
		Authentication: at,
		Exporter:       be,
	}

	// encrypt data for the password
	sink, err = svc.encryptMetadata(cfg, sink)
	if err != nil {
		return Sink{}, err
	}

	id, err := svc.sinkRepo.Save(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	sink.ID = id

	// After creating, decrypt Metadata to send correct information to Redis
	sink, err = svc.decryptMetadata(cfg, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrCreateSink, err)
	}
	return sink, nil
}

func validateAuthType(s *Sink) (authentication_type.AuthenticationType, error) {
	var authMetadata types.Metadata
	if len(s.ConfigData) != 0 {
		var helper types.Metadata
		if s.Format == "yaml" {
			err := yaml.Unmarshal([]byte(s.ConfigData), &helper)
			if err != nil {
				return nil, err
			}
			authMetadata = helper.GetSubMetadata(authentication_type.AuthenticationKey)
		} else {
			return nil, errors.New("config format not supported")
		}
	} else {
		authMetadata = s.Config.GetSubMetadata(authentication_type.AuthenticationKey)
	}
	authTypeStr, ok := authMetadata["type"]
	if !ok {
		return nil, errors.Wrap(errors.ErrAuthTypeNotFound, errors.New("authentication type not found"))
	}

	if _, ok := authTypeStr.(string); !ok {
		return nil, errors.Wrap(errors.ErrInvalidAuthType, errors.New("invalid authentication type"))
	}

	authType, ok := authentication_type.GetAuthType(authTypeStr.(string))
	if !ok {
		return nil, errors.Wrap(errors.ErrInvalidAuthType, errors.New("invalid authentication type"))
	}

	err := authType.ValidateConfiguration("object", authMetadata)
	if err != nil {
		return nil, err
	}

	return authType, nil
}

func (svc sinkService) encryptMetadata(configSvc Configuration, sink Sink) (Sink, error) {
	var err error
	if sink.Config != nil {
		encodeMetadata, err := configSvc.Authentication.EncodeInformation("object", sink.Config)
		if err != nil {
			svc.logger.Error("error on parsing encrypted config in data")
			return sink, err
		}
		sink.Config = encodeMetadata.(types.Metadata)
	}
	if sink.ConfigData != "" {
		encodeMetadata, err := configSvc.Authentication.EncodeInformation("yaml", sink.ConfigData)
		if err != nil {
			svc.logger.Error("error on parsing encrypted config in data")
			return sink, err
		}
		sink.ConfigData = encodeMetadata.(string)
	}
	return sink, err
}

func (svc sinkService) ViewAuthenticationType(ctx context.Context, token string, key string) (authentication_type.AuthenticationTypeConfig, error) {
	_, err := svc.identify(token)
	if err != nil {
		return authentication_type.AuthenticationTypeConfig{}, err
	}

	value, ok := authentication_type.GetAuthType(key)
	if !ok {
		return authentication_type.AuthenticationTypeConfig{}, errors.New("invalid authentication type given name")
	}
	return value.Metadata(), nil
}

func (svc sinkService) ListAuthenticationTypes(ctx context.Context, token string) ([]authentication_type.AuthenticationTypeConfig, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}

	value := authentication_type.GetList()

	return value, nil
}

func (svc sinkService) decryptMetadata(configSvc Configuration, sink Sink) (Sink, error) {
	var err error
	if sink.Config != nil {
		decodeMetadata, err := configSvc.Authentication.DecodeInformation("object", sink.Config)
		if err != nil {
			svc.logger.Error("error on parsing encrypted config in data")
			return sink, err
		}
		sink.Config = decodeMetadata.(types.Metadata)
	}
	if sink.ConfigData != "" {
		decodeMetadata, err := configSvc.Authentication.DecodeInformation("yaml", sink.ConfigData)
		if err != nil {
			svc.logger.Error("error on parsing encrypted config in data")
			return sink, err
		}
		sink.ConfigData = decodeMetadata.(string)
	}
	return sink, err
}

func (svc sinkService) UpdateSinkInternal(ctx context.Context, sink Sink) (Sink, error) {
	var currentSink Sink
	currentSink, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	var cfg Configuration
	if sink.Config == nil && sink.ConfigData == "" {
		// No config sent, keep the previous
		sink.Config = currentSink.Config
		authType, _ := authentication_type.GetAuthType(sink.GetAuthenticationTypeName())
		be := backend.GetBackend(currentSink.Backend)
		cfg = Configuration{
			Authentication: authType,
			Exporter:       be,
		}

		// get the decrypted config, otherwise the password would be encrypted again
		sink, err = svc.decryptMetadata(cfg, sink)
		if err != nil {
			return Sink{}, errors.Wrap(ErrUpdateEntity, err)
		}
	} else {
		sink.Backend = currentSink.Backend
		// we still need to validate here
		be, err := svc.validateBackend(&sink)
		if err != nil {
			return Sink{}, errors.Wrap(ErrMalformedEntity, err)
		}
		at, err := validateAuthType(&sink)
		if err != nil {
			return Sink{}, errors.Wrap(ErrMalformedEntity, err)
		}
		cfg = Configuration{
			Authentication: at,
			Exporter:       be,
		}
		sink.State = Unknown
		sink.Error = ""
		if sink.Format == "yaml" {
			configDataByte, err := yaml.Marshal(sink.Config)
			if err != nil {
				return Sink{}, errors.Wrap(ErrMalformedEntity, err)
			}
			sink.ConfigData = string(configDataByte)
		}
	}

	if sink.Tags == nil {
		sink.Tags = currentSink.Tags
	}

	if sink.Description == nil {
		sink.Description = currentSink.Description
	}

	if newName := sink.Name.String(); newName == "" {
		sink.Name = currentSink.Name
	}

	sink.MFOwnerID = currentSink.MFOwnerID
	if sink.Backend == "" && currentSink.Backend != "" {
		sink.Backend = currentSink.Backend
	}
	sink, err = svc.encryptMetadata(cfg, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	err = svc.sinkRepo.Update(ctx, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	sinkEdited, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	sinkEdited, err = svc.decryptMetadata(cfg, sinkEdited)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}

	return sinkEdited, nil
}

func (svc sinkService) UpdateSink(ctx context.Context, token string, sink Sink) (Sink, error) {
	skOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	currentSink, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, err
	}

	authType, _ := authentication_type.GetAuthType(currentSink.GetAuthenticationTypeName())
	be := backend.GetBackend(currentSink.Backend)
	cfg := Configuration{
		Authentication: authType,
		Exporter:       be,
	}

	// get the decrypted config, otherwise the password would be encrypted again
	currentSink, err = svc.decryptMetadata(cfg, currentSink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}

	if sink.Config == nil && sink.ConfigData == "" {
		// No config sent, keep the previous
		sink.Config = currentSink.Config
		sink.ConfigData = currentSink.ConfigData
	} else {
		sink.Backend = currentSink.Backend
		be, err := svc.validateBackend(&sink)
		if err != nil {
			return Sink{}, errors.Wrap(errors.New("incorrect backend and exporter configuration"), err)
		}
		at, err := validateAuthType(&sink)
		if err != nil {
			return Sink{}, errors.Wrap(errors.New("incorrect authentication configuration"), err)
		}
		cfg = Configuration{
			Authentication: at,
			Exporter:       be,
		}

		// check if the password is encrypted and decrypt it if it is
		if existingAuth := sink.Config.GetSubMetadata(authentication_type.AuthenticationKey); existingAuth != nil {
			if password, ok := existingAuth["password"]; ok {
				// if the password is encrypted, it will be a hex string
				if _, err := hex.DecodeString(password.(string)); err == nil {
					if sink, err = svc.decryptMetadata(cfg, sink); err != nil {
						return Sink{}, errors.Wrap(ErrUpdateEntity, err)
					}
				}
			}
		}


		if sink.Format == "yaml" {
			configDataByte, err := yaml.Marshal(sink.Config)
			if err != nil {
				svc.logger.Error("failed to marshal config data", zap.Error(err))
				return Sink{}, errors.Wrap(errors.New("configuration is invalid for yaml format"), err)
			}
			sink.ConfigData = string(configDataByte)
		}
	}

	if sink.Tags == nil {
		sink.Tags = currentSink.Tags
	}

	if sink.Description == nil {
		sink.Description = currentSink.Description
	}

	if newName := sink.Name.String(); newName == "" {
		sink.Name = currentSink.Name
	}

	sink.MFOwnerID = skOwnerID
	if sink.Backend == "" && currentSink.Backend != "" {
		sink.Backend = currentSink.Backend
	}
	sink, err = svc.encryptMetadata(cfg, sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}
	err = svc.sinkRepo.Update(ctx, sink)
	if err != nil {
		return Sink{}, err
	}
	sinkEdited, err := svc.sinkRepo.RetrieveById(ctx, sink.ID)
	if err != nil {
		return Sink{}, err
	}
	sinkEdited, err = svc.decryptMetadata(cfg, sinkEdited)
	if err != nil {
		return Sink{}, errors.Wrap(ErrUpdateEntity, err)
	}

	return sinkEdited, nil
}

func (svc sinkService) ListBackends(ctx context.Context, token string) ([]string, error) {
	_, err := svc.identify(token)
	if err != nil {
		return []string{}, err
	}
	return backend.GetList(), nil
}

func (svc sinkService) ViewBackend(ctx context.Context, token string, key string) (backend.Backend, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	res := backend.GetBackend(key)
	if res == nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ViewSink(ctx context.Context, token string, key string) (Sink, error) {
	_, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}
	res, err := svc.sinkRepo.RetrieveById(ctx, key)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrNotFound, err)
	}
	return res, nil
}

func (svc sinkService) ViewSinkInternal(ctx context.Context, ownerID string, key string) (Sink, error) {
	res, err := svc.sinkRepo.RetrieveByOwnerAndId(ctx, ownerID, key)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrNotFound, err)
	}
	authType, _ := authentication_type.GetAuthType(res.GetAuthenticationTypeName())
	be := backend.GetBackend(res.Backend)
	cfg := Configuration{
		Authentication: authType,
		Exporter:       be,
	}
	res, err = svc.decryptMetadata(cfg, res)
	if err != nil {
		return Sink{}, errors.Wrap(errors.ErrViewEntity, err)
	}
	return res, nil
}

func (svc sinkService) ListSinksInternal(ctx context.Context, filter Filter) (sinksResp Page, err error) {
	sinks, err := svc.sinkRepo.SearchAllSinks(ctx, filter)
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrNotFound, err)
	}
	for _, sink := range sinks {
		authType, _ := authentication_type.GetAuthType(sink.GetAuthenticationTypeName())
		be := backend.GetBackend(sink.Backend)
		cfg := Configuration{
			Authentication: authType,
			Exporter:       be,
		}
		sink, err = svc.decryptMetadata(cfg, sink)
		if err != nil {
			return Page{}, errors.Wrap(errors.ErrViewEntity, err)
		}
		sinksResp.Sinks = append(sinksResp.Sinks, sink)
	}

	return
}

func (svc sinkService) ListSinks(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.identify(token)
	if err != nil {
		svc.GetLogger().Error("got error on identifying token", zap.Error(err))
		return Page{}, err
	}

	return svc.sinkRepo.RetrieveAllByOwnerID(ctx, res, pm)
}

func (svc sinkService) DeleteSink(ctx context.Context, token string, id string) error {
	res, err := svc.identify(token)
	if err != nil {
		return err
	}

	return svc.sinkRepo.Remove(ctx, res, id)
}

func (svc sinkService) ValidateSink(ctx context.Context, token string, sink Sink) (Sink, error) {

	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Sink{}, err
	}

	sink.MFOwnerID = mfOwnerID

	_, err = svc.validateBackend(&sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrValidateSink, err)
	}

	_, err = validateAuthType(&sink)
	if err != nil {
		return Sink{}, errors.Wrap(ErrValidateSink, err)
	}

	return sink, nil
}

func (svc sinkService) ChangeSinkStateInternal(ctx context.Context, sinkID string, msg string, ownerID string, state State) error {
	return svc.sinkRepo.UpdateSinkState(ctx, sinkID, msg, ownerID, state)
}

func (svc sinkService) validateBackend(sink *Sink) (be backend.Backend, err error) {
	if !backend.HaveBackend(sink.Backend) {
		return nil, ErrInvalidBackend
	}
	sinkBe := backend.GetBackend(sink.Backend)
	if len(sink.ConfigData) == 0 {
		config := sink.Config.GetSubMetadata("exporter")
		if config == nil {
			return nil, errors.Wrap(ErrInvalidBackend, errors.New("missing exporter configuration"))
		}
		return sinkBe, sinkBe.ValidateConfiguration(config)
	} else {
		parseConfig, err := sinkBe.ParseConfig("yaml", sink.ConfigData)
		if err != nil {
			return nil, errors.Wrap(ErrInvalidBackend, err)
		}
		sink.Config = parseConfig
		config2 := sink.Config.GetSubMetadata("exporter")
		if config2 == nil {
			return nil, errors.Wrap(ErrInvalidBackend, errors.New("missing exporter configuration"))
		}
		return sinkBe, sinkBe.ValidateConfiguration(config2)
	}
}
