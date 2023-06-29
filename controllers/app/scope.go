/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/PDeXchange/pac/apis/app/v1alpha1"
	"github.com/PDeXchange/pac/controllers/util"
	"github.com/PDeXchange/pac/internal/pkg/client/platform"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/IBM/go-sdk-core/v5/core"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/PDeXchange/pac/internal/pkg/client/powervs"
)

const (
	catalogController = "catalog"
	serviceController = "service"
)

type ControllerScopeParams struct {
	Logger  logr.Logger
	Client  client.Client
	Type    string
	Catalog *v1alpha1.Catalog
	Debug   bool
}

type ControllerScope struct {
	logr.Logger
	Client         client.Client
	Catalog        *v1alpha1.Catalog
	PowerVSClient  *powervs.Client
	PlatformClient *platform.Client
}

type CatalogScopeParams struct {
	ControllerScopeParams
}

type ServiceScopeParams struct {
	ControllerScopeParams
	Service *v1alpha1.Service
}

type CatalogScope struct {
	ControllerScope
	catalogPatchHelper *patch.Helper
}

type ServiceScope struct {
	ControllerScope
	servicePatchHelper *patch.Helper
	Service            *v1alpha1.Service
}

func NewCatalogScope(ctx context.Context, params CatalogScopeParams) (scope *CatalogScope, err error) {
	scope = &CatalogScope{}

	ctrlScope, err := NewControllerScope(ctx, params.ControllerScopeParams)
	if err != nil {
		err = errors.Wrap(err, "failed to init controller scope")
		return nil, err
	}

	scope.ControllerScope = *ctrlScope

	catalogHelper, err := patch.NewHelper(params.Catalog, params.Client)
	if err != nil {
		err = errors.Wrap(err, "failed to init patch helper")
		return nil, err
	}
	scope.catalogPatchHelper = catalogHelper

	return scope, nil
}

func NewServiceScope(ctx context.Context, params ServiceScopeParams) (scope *ServiceScope, err error) {
	scope = &ServiceScope{}

	ctrlScope, err := NewControllerScope(ctx, params.ControllerScopeParams)
	if err != nil {
		err = errors.Wrap(err, "failed to init controller scope")
		return nil, err
	}

	scope.ControllerScope = *ctrlScope

	serviceHelper, err := patch.NewHelper(params.Service, params.Client)
	if err != nil {
		err = errors.Wrap(err, "failed to init patch helper")
		return nil, err
	}
	scope.servicePatchHelper = serviceHelper

	if params.Service == nil {
		err = errors.New("service is required when creating a ServiceScope")
		return
	}
	scope.Service = params.Service

	return scope, nil
}

func NewControllerScope(ctx context.Context, params ControllerScopeParams) (scope *ControllerScope, err error) {
	scope = &ControllerScope{}

	if params.Client == nil {
		err = errors.New("client is required when creating a CatalogScope")
		return
	}
	scope.Client = params.Client

	if params.Logger == (logr.Logger{}) {
		params.Logger = zap.New()
	}
	scope.Logger = params.Logger

	if params.Catalog == nil {
		err = errors.New("catalog is required when creating a scope for catalog and service controller")
		return
	}
	scope.Catalog = params.Catalog

	platformClient, err := platform.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating platform services client")
	}
	scope.PlatformClient = platformClient

	var cloudInstanceID, zone, accountID string
	switch params.Catalog.Spec.Type {
	case v1alpha1.CatalogTypeVM:
		cloudInstanceID, zone, accountID, err = util.ParsePowerVSCRN(params.Catalog.Spec.VM.CRN)
		if err != nil {
			return nil, err
		}
	}

	powerVSClient, err := powervs.NewClient(ctx, powervs.Options{
		AccountID:       accountID,
		CloudInstanceID: cloudInstanceID,
		Zone:            zone,
		Debug:           params.Debug})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create powervs client")
	}
	scope.PowerVSClient = powerVSClient

	if params.Debug {
		core.SetLoggingLevel(core.LevelDebug)
	}

	return scope, nil
}

// PatchObject persists the catalog/service configuration and status.
func (m *CatalogScope) PatchCatalogObject() error {
	return m.catalogPatchHelper.Patch(context.TODO(), m.Catalog)
}

// PatchObject persists the catalog/service configuration and status.
func (m *ServiceScope) PatchServiceObject() error {
	return m.servicePatchHelper.Patch(context.TODO(), m.Service)
}
