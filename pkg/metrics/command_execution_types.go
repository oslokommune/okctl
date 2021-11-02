package metrics

import metricsapi "github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics"

// CategoryCommandExecution represents the context of running commands
const CategoryCommandExecution = metricsapi.CategoryCommandExecution

const (

	// ActionScaffoldCluster represents running the command `okctl scaffold cluster`
	ActionScaffoldCluster = metricsapi.ActionScaffoldCluster
	// ActionApplyCluster represents running the command `okctl apply cluster`
	ActionApplyCluster = metricsapi.ActionApplyCluster
	// ActionDeleteCluster represents running the command `okctl delete cluster`
	ActionDeleteCluster = metricsapi.ActionDeleteCluster

	// ActionScaffoldApplication represents running the command `okctl scaffold application`
	ActionScaffoldApplication = metricsapi.ActionScaffoldApplication
	// ActionApplyApplication represents running the command `okctl apply application`
	ActionApplyApplication = metricsapi.ActionApplyApplication

	// ActionForwardPostgres represents running the command `okctl forward postgres`
	ActionForwardPostgres = metricsapi.ActionForwardPostgres
	// ActionAttachPostgres represents running the command `okctl attach postgres`
	ActionAttachPostgres = metricsapi.ActionAttachPostgres

	// ActionShowCredentials represents running the command `okctl show credentials`
	ActionShowCredentials = metricsapi.ActionShowCredentials
	// ActionUpgrade represents running the command `okctl upgrade`
	ActionUpgrade = metricsapi.ActionUpgrade
	// ActionVenv represents running the command `okctl venv`
	ActionVenv = metricsapi.ActionVenv
	// ActionVersion represents running the command `okctl version`
	ActionVersion = metricsapi.ActionVersion
)

const (
	LabelPhaseKey = metricsapi.LabelPhaseKey
	// LabelPhaseStart represents the start of a command
	LabelPhaseStart = metricsapi.LabelPhaseStart
	// LabelPhaseEnd represents the end of the command
	LabelPhaseEnd = metricsapi.LabelPhaseEnd
)
