package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
)

const usersReconcilerIdentifier = "Users"

type usersReconciler struct {
	client client.IdentityManagerService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *usersReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	userPoolID := ""

	identityPoolExists, err := state.IdentityManager.HasIdentityPool()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("checking identity pool existence: %w", err)
	}

	if !identityPoolExists && !meta.ClusterDeclaration.Integrations.Cognito {
		meta.Purge = true
	}

	actionMap, err := z.determineActionsForUsers(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	hasCreate := hasCreateAction(actionMap)
	if hasCreate {
		if !identityPoolExists {
			return reconciliation.Result{Requeue: true}, nil
		}

		im, err := state.IdentityManager.GetIdentityPool(
			cfn.NewStackNamer().IdentityPool(meta.ClusterDeclaration.Metadata.Name),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting identity pool: %w", err)
		}

		userPoolID = im.UserPoolID
	}

	err = z.handleActionMap(ctx, handleOpts{
		actionMap:  actionMap,
		meta:       meta,
		userPoolID: userPoolID,
	})
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("handling action map: %w", err)
	}

	return reconciliation.Result{Requeue: false}, nil
}

type handleOpts struct {
	actionMap  map[string]reconciliation.Action
	meta       reconciliation.Metadata
	userPoolID string
}

func (z *usersReconciler) handleActionMap(ctx context.Context, opts handleOpts) error {
	clusterID := reconciliation.ClusterMetaAsID(opts.meta.ClusterDeclaration.Metadata)

	for email, action := range opts.actionMap {
		switch action {
		case reconciliation.ActionCreate:
			_, err := z.client.CreateIdentityPoolUser(ctx, client.CreateIdentityPoolUserOpts{
				ID:         clusterID,
				Email:      email,
				UserPoolID: opts.userPoolID,
			})
			if err != nil {
				return fmt.Errorf("creating identity pool user %s: %w", email, err)
			}

			continue
		case reconciliation.ActionDelete:
			err := z.client.DeleteIdentityPoolUser(ctx, client.DeleteIdentityPoolUserOpts{
				ClusterID: clusterID,
				UserEmail: email,
			})
			if err != nil {
				return fmt.Errorf("deleting identity pool user %s: %w", email, err)
			}

			continue
		case reconciliation.ActionNoop:
			continue
		default:
			return reconciliation.ErrIndecisive
		}
	}

	return nil
}

func (z *usersReconciler) determineActionsForUsers(meta reconciliation.Metadata, state *clientCore.StateHandlers) (map[string]reconciliation.Action, error) {
	currentlyIndicatedUsers := meta.ClusterDeclaration.Users

	existingUsers, err := state.IdentityManager.GetIdentityPoolUsers()
	if err != nil {
		return nil, fmt.Errorf("getting existing identity pool users: %w", err)
	}

	actionMap := make(map[string]reconciliation.Action)

	for _, user := range existingUsers {
		actionMap[user.Email] = reconciliation.ActionDelete
	}

	if meta.Purge {
		return actionMap, nil
	}

	for _, user := range currentlyIndicatedUsers {
		_, ok := actionMap[user.Email]
		if ok {
			actionMap[user.Email] = reconciliation.ActionNoop
		} else {
			actionMap[user.Email] = reconciliation.ActionCreate
		}
	}

	return actionMap, nil
}

func hasCreateAction(actionMap map[string]reconciliation.Action) bool {
	for _, action := range actionMap {
		if action == reconciliation.ActionCreate {
			return true
		}
	}

	return false
}

// String returns the identifier type
func (z *usersReconciler) String() string {
	return usersReconcilerIdentifier
}

// NewUsersReconciler creates a new reconciler for the Users resource
func NewUsersReconciler(client client.IdentityManagerService) reconciliation.Reconciler {
	return &usersReconciler{
		client: client,
	}
}
