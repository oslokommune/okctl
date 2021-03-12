package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type componentState struct {
	state state.Componenter
}

const dbTypePostgres = "postgres"

func (c *componentState) SavePostgresDatabase(database *client.PostgresDatabase) (*store.Report, error) {
	db := c.state.GetDatabase(database.ApplicationName)

	db.Namespace = database.Namespace
	db.SecurityGroupID = database.OutgoingSecurityGroupID
	db.EndpointPort = database.EndpointPort
	db.EndpointAddress = database.EndpointAddress
	db.Type = dbTypePostgres
	db.AdminSecretName = database.AdminSecretName
	db.DatabaseConfigMapName = database.DatabaseConfigMapName

	report, err := c.state.SaveDatabase(database.ApplicationName, db)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "Postgres",
			Path: fmt.Sprintf("address=%s, port=%d", database.EndpointAddress, database.EndpointPort),
			Type: "StateUpdate[add]",
		},
	}, report.Actions...)

	return report, nil
}

func (c *componentState) RemovePostgresDatabase(applicationName string) (*store.Report, error) {
	report, err := c.state.DeleteDatabase(applicationName)
	if err != nil {
		return nil, err
	}

	report.Actions = append([]store.Action{
		{
			Name: "Postgres",
			Path: fmt.Sprintf("application=%s", applicationName),
			Type: "StateUpdate[removed]",
		},
	}, report.Actions...)

	return report, nil
}

// NewComponentState returns an initialised state updated
func NewComponentState(state state.Componenter) client.ComponentState {
	return &componentState{
		state: state,
	}
}
