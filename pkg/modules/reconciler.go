package modules

import (
	"context"
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/spf13/afero"
)

// Reconcile ensures all modules specified in cluster.yaml are installed into the IAC repository
func (m modulesReconciler) Reconcile(_ context.Context, meta reconciliation.Metadata, _ *clientCore.StateHandlers) (reconciliation.Result, error) {
	for _, module := range meta.ClusterDeclaration.Modules {
		err := InstallModule(m.fs, module, m.modulesBaseDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("scaffolding module: %w", err)
		}
	}

	return reconciliation.Result{Requeue: false}, nil
}

func (m modulesReconciler) String() string {
	return reconcilerID
}

// NewReconciler returns an initialized modules reconciler
func NewReconciler(fs *afero.Afero, modulesBaseDir string) reconciliation.Reconciler {
	return &modulesReconciler{
		fs:             fs,
		modulesBaseDir: modulesBaseDir,
	}
}

const reconcilerID = "Modules"

type modulesReconciler struct {
	fs             *afero.Afero
	modulesBaseDir string
}
