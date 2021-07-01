package reconciliation

import (
	"time"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/constant"
)

// ClusterMetaAsID knows how to convert cluster declaration metadata to an api.ID struct
func ClusterMetaAsID(meta v1alpha1.ClusterMeta) api.ID {
	return api.ID{
		Region:       meta.Region,
		AWSAccountID: meta.AccountID,
		ClusterName:  meta.Name,
	}
}

// DetermineUserIndication knows how to interpret what operation the user wants for the certain reconciler
func DetermineUserIndication(metadata Metadata, componentFlag bool) Action {
	if metadata.Purge || !componentFlag {
		return ActionDelete
	}

	return ActionCreate
}

// DefaultDelayFunction defines a sane default reconciliation loop delay function
func DefaultDelayFunction() {
	time.Sleep(constant.DefaultReconciliationLoopDelayDuration)
}
