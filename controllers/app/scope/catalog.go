package scope

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api/util/patch"
)

type CatalogScope struct {
	ControllerScope
	catalogPatchHelper *patch.Helper
}

func NewCatalogScope(ctx context.Context, params CatalogScopeParams) (*CatalogScope, error) {
	scope := &CatalogScope{}

	ctrlScope, err := NewControllerScope(ctx, params.ControllerScopeParams)
	if err != nil {
		return scope, fmt.Errorf("failed to init controller scope: %w", err)
	}

	scope.ControllerScope = *ctrlScope

	catalogHelper, err := patch.NewHelper(params.Catalog, params.Client)
	if err != nil {
		return scope, fmt.Errorf("failed to init patch helper: %w", err)
	}
	scope.catalogPatchHelper = catalogHelper

	return scope, nil
}

// PatchObject persists the catalog/service configuration and status.
func (m *CatalogScope) PatchCatalogObject() error {
	return m.catalogPatchHelper.Patch(context.TODO(), m.Catalog)
}
