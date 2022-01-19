package metrics

import "github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics/types/commandexecution"

// CategoryCommandExecution represents the context of running commands
const CategoryCommandExecution = commandexecution.Category

const (

	// ActionScaffoldCluster represents running the command `okctl scaffold cluster`
	ActionScaffoldCluster = commandexecution.ActionScaffoldCluster
	// ActionApplyCluster represents running the command `okctl apply cluster`
	ActionApplyCluster = commandexecution.ActionApplyCluster
	// ActionDeleteCluster represents running the command `okctl delete cluster`
	ActionDeleteCluster = commandexecution.ActionDeleteCluster

	// ActionScaffoldApplication represents running the command `okctl scaffold application`
	ActionScaffoldApplication = commandexecution.ActionScaffoldApplication
	// ActionApplyApplication represents running the command `okctl apply application`
	ActionApplyApplication = commandexecution.ActionApplyApplication
	// ActionDeleteApplication represents running the command `okctl delete application`
	ActionDeleteApplication = commandexecution.ActionDeleteApplication

	// ActionForwardPostgres represents running the command `okctl forward postgres`
	ActionForwardPostgres = commandexecution.ActionForwardPostgres
	// ActionAttachPostgres represents running the command `okctl attach postgres`
	ActionAttachPostgres = commandexecution.ActionAttachPostgres

	// ActionShowCredentials represents running the command `okctl show credentials`
	ActionShowCredentials = commandexecution.ActionShowCredentials
	// ActionUpgrade represents running the command `okctl upgrade`
	ActionUpgrade = commandexecution.ActionUpgrade
	// ActionVenv represents running the command `okctl venv`
	ActionVenv = commandexecution.ActionVenv
	// ActionVersion represents running the command `okctl version`
	ActionVersion = commandexecution.ActionVersion

	// ActionMaintenanceStateAcquireLock represents running the command `okctl maintenance state-acquire-lock
	ActionMaintenanceStateAcquireLock = commandexecution.ActionMaintenanceStateAcquireLock
	// ActionMaintenanceStateReleaseLock represents running the command `okctl maintenance state-release-lock
	ActionMaintenanceStateReleaseLock = commandexecution.ActionMaintenanceStateReleaseLock
	// ActionMaintenanceStateDownload represents running the command `okctl maintenance state-download
	ActionMaintenanceStateDownload = commandexecution.ActionMaintenanceStateDownload
	// ActionMaintenanceStateUpload represents running the command `okctl maintenance state-upload
	ActionMaintenanceStateUpload = commandexecution.ActionMaintenanceStateUpload
)

const (
	// LabelPhaseKey represents the key for the phase labels
	LabelPhaseKey = commandexecution.LabelPhaseKey
	// LabelPhaseStart represents the start of a command
	LabelPhaseStart = commandexecution.LabelPhaseStart
	// LabelPhaseEnd represents the end of the command
	LabelPhaseEnd = commandexecution.LabelPhaseEnd
)
