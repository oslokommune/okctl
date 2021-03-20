package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"

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
	db.AdminSecretARN = database.SecretsManagerAdminSecretARN
	db.DatabaseConfigMapName = database.DatabaseConfigMapName
	db.RotaterLambdaRoleARN = database.LambdaRoleARN
	db.RotaterLambdaPolicyARN = database.LambdaPolicyARN

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

func (c *componentState) GetPostgresDatabase(applicationName string) (*client.PostgresDatabase, error) {
	db := c.state.GetDatabase(applicationName)

	return &client.PostgresDatabase{
		Namespace:             db.Namespace,
		AdminSecretName:       db.AdminSecretName,
		AdminSecretARN:        db.AdminSecretARN,
		DatabaseConfigMapName: db.DatabaseConfigMapName,
		PostgresDatabase: &api.PostgresDatabase{
			ApplicationName:         applicationName,
			EndpointAddress:         db.EndpointAddress,
			EndpointPort:            db.EndpointPort,
			OutgoingSecurityGroupID: db.SecurityGroupID,
			LambdaPolicyARN:         db.RotaterLambdaPolicyARN,
			LambdaRoleARN:           db.RotaterLambdaRoleARN,
		},
	}, nil
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
