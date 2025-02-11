package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Runs = (*runs)(nil)

// Runs describes all the run related methods that the Terraform Enterprise
// API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/run.html
type Runs interface {
	// List all the runs of the given workspace.
	List(ctx context.Context, workspaceID string, options *RunListOptions) (*RunList, error)

	// Create a new run with the given options.
	Create(ctx context.Context, options RunCreateOptions) (*Run, error)

	// Read a run by its ID.
	Read(ctx context.Context, runID string) (*Run, error)

	// ReadWithOptions reads a run by its ID using the options supplied
	ReadWithOptions(ctx context.Context, runID string, options *RunReadOptions) (*Run, error)

	// Apply a run by its ID.
	Apply(ctx context.Context, runID string, options RunApplyOptions) error

	// Cancel a run by its ID.
	Cancel(ctx context.Context, runID string, options RunCancelOptions) error

	// Force-cancel a run by its ID.
	ForceCancel(ctx context.Context, runID string, options RunForceCancelOptions) error

	// Discard a run by its ID.
	Discard(ctx context.Context, runID string, options RunDiscardOptions) error
}

// runs implements Runs.
type runs struct {
	client *Client
}

// RunStatus represents a run state.
type RunStatus string

// List all available run statuses.
const (
	RunApplied            RunStatus = "applied"
	RunApplying           RunStatus = "applying"
	RunApplyQueued        RunStatus = "apply_queued"
	RunCanceled           RunStatus = "canceled"
	RunConfirmed          RunStatus = "confirmed"
	RunCostEstimated      RunStatus = "cost_estimated"
	RunCostEstimating     RunStatus = "cost_estimating"
	RunDiscarded          RunStatus = "discarded"
	RunErrored            RunStatus = "errored"
	RunFetching           RunStatus = "fetching"
	RunFetchingCompleted  RunStatus = "fetching_completed"
	RunPending            RunStatus = "pending"
	RunPlanned            RunStatus = "planned"
	RunPlannedAndFinished RunStatus = "planned_and_finished"
	RunPlanning           RunStatus = "planning"
	RunPlanQueued         RunStatus = "plan_queued"
	RunPolicyChecked      RunStatus = "policy_checked"
	RunPolicyChecking     RunStatus = "policy_checking"
	RunPolicyOverride     RunStatus = "policy_override"
	RunPolicySoftFailed   RunStatus = "policy_soft_failed"
	RunPostPlanCompleted  RunStatus = "post_plan_completed"
	RunPostPlanRunning    RunStatus = "post_plan_running"
	RunPrePlanCompleted   RunStatus = "pre_plan_completed"
	RunPrePlanRunning     RunStatus = "pre_plan_running"
	RunQueuing            RunStatus = "queuing"
)

// RunSource represents a source type of a run.
type RunSource string

// List all available run sources.
const (
	RunSourceAPI                  RunSource = "tfe-api"
	RunSourceConfigurationVersion RunSource = "tfe-configuration-version"
	RunSourceUI                   RunSource = "tfe-ui"
)

// RunOperation represents an operation type of run.
type RunOperation string

// List all available run operations.
const (
	RunOperationPlanApply   RunOperation = "plan_and_apply"
	RunOperationPlanOnly    RunOperation = "plan_only"
	RunOperationRefreshOnly RunOperation = "refresh_only"
	RunOperationDestroy     RunOperation = "destroy"
	RunOperationEmptyApply  RunOperation = "empty_apply"
)

// RunList represents a list of runs.
type RunList struct {
	*Pagination
	Items []*Run
}

// Run represents a Terraform Enterprise run.
type Run struct {
	ID                     string               `jsonapi:"primary,runs"`
	Actions                *RunActions          `jsonapi:"attr,actions"`
	AutoApply              bool                 `jsonapi:"attr,auto-apply,omitempty"`
	AllowEmptyApply        bool                 `jsonapi:"attr,allow-empty-apply"`
	CreatedAt              time.Time            `jsonapi:"attr,created-at,iso8601"`
	ForceCancelAvailableAt time.Time            `jsonapi:"attr,force-cancel-available-at,iso8601"`
	HasChanges             bool                 `jsonapi:"attr,has-changes"`
	IsDestroy              bool                 `jsonapi:"attr,is-destroy"`
	Message                string               `jsonapi:"attr,message"`
	Permissions            *RunPermissions      `jsonapi:"attr,permissions"`
	PositionInQueue        int                  `jsonapi:"attr,position-in-queue"`
	PlanOnly               bool                 `jsonapi:"attr,plan-only"`
	Refresh                bool                 `jsonapi:"attr,refresh"`
	RefreshOnly            bool                 `jsonapi:"attr,refresh-only"`
	ReplaceAddrs           []string             `jsonapi:"attr,replace-addrs,omitempty"`
	Source                 RunSource            `jsonapi:"attr,source"`
	Status                 RunStatus            `jsonapi:"attr,status"`
	StatusTimestamps       *RunStatusTimestamps `jsonapi:"attr,status-timestamps"`
	TargetAddrs            []string             `jsonapi:"attr,target-addrs,omitempty"`
	TerraformVersion       string               `jsonapi:"attr,terraform-version"`
	Variables              []*RunVariableAttr   `jsonapi:"attr,variables"`

	// Relations
	Apply                *Apply                `jsonapi:"relation,apply"`
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`
	CostEstimate         *CostEstimate         `jsonapi:"relation,cost-estimate"`
	CreatedBy            *User                 `jsonapi:"relation,created-by"`
	Plan                 *Plan                 `jsonapi:"relation,plan"`
	PolicyChecks         []*PolicyCheck        `jsonapi:"relation,policy-checks"`
	TaskStages           []*TaskStage          `jsonapi:"relation,task-stages,omitempty"`
	Workspace            *Workspace            `jsonapi:"relation,workspace"`
	Comments             []*Comment            `jsonapi:"relation,comments"`
}

// RunActions represents the run actions.
type RunActions struct {
	IsCancelable      bool `jsonapi:"attr,is-cancelable"`
	IsConfirmable     bool `jsonapi:"attr,is-confirmable"`
	IsDiscardable     bool `jsonapi:"attr,is-discardable"`
	IsForceCancelable bool `jsonapi:"attr,is-force-cancelable"`
}

// RunPermissions represents the run permissions.
type RunPermissions struct {
	CanApply        bool `jsonapi:"attr,can-apply"`
	CanCancel       bool `jsonapi:"attr,can-cancel"`
	CanDiscard      bool `jsonapi:"attr,can-discard"`
	CanForceCancel  bool `jsonapi:"attr,can-force-cancel"`
	CanForceExecute bool `jsonapi:"attr,can-force-execute"`
}

// RunStatusTimestamps holds the timestamps for individual run statuses.
type RunStatusTimestamps struct {
	AppliedAt            time.Time `jsonapi:"attr,applied-at,rfc3339"`
	ApplyingAt           time.Time `jsonapi:"attr,applying-at,rfc3339"`
	ApplyQueuedAt        time.Time `jsonapi:"attr,apply-queued-at,rfc3339"`
	CanceledAt           time.Time `jsonapi:"attr,canceled-at,rfc3339"`
	ConfirmedAt          time.Time `jsonapi:"attr,confirmed-at,rfc3339"`
	CostEstimatedAt      time.Time `jsonapi:"attr,cost-estimated-at,rfc3339"`
	CostEstimatingAt     time.Time `jsonapi:"attr,cost-estimating-at,rfc3339"`
	DiscardedAt          time.Time `jsonapi:"attr,discarded-at,rfc3339"`
	ErroredAt            time.Time `jsonapi:"attr,errored-at,rfc3339"`
	FetchedAt            time.Time `jsonapi:"attr,fetched-at,rfc3339"`
	FetchingAt           time.Time `jsonapi:"attr,fetching-at,rfc3339"`
	ForceCanceledAt      time.Time `jsonapi:"attr,force-canceled-at,rfc3339"`
	PlannedAndFinishedAt time.Time `jsonapi:"attr,planned-and-finished-at,rfc3339"`
	PlannedAt            time.Time `jsonapi:"attr,planned-at,rfc3339"`
	PlanningAt           time.Time `jsonapi:"attr,planning-at,rfc3339"`
	PlanQueueableAt      time.Time `jsonapi:"attr,plan-queueable-at,rfc3339"`
	PlanQueuedAt         time.Time `jsonapi:"attr,plan-queued-at,rfc3339"`
	PolicyCheckedAt      time.Time `jsonapi:"attr,policy-checked-at,rfc3339"`
	PolicySoftFailedAt   time.Time `jsonapi:"attr,policy-soft-failed-at,rfc3339"`
	PostPlanCompletedAt  time.Time `jsonapi:"attr,post-plan-completed-at,rfc3339"`
	PostPlanRunningAt    time.Time `jsonapi:"attr,post-plan-running-at,rfc3339"`
	PrePlanCompletedAt   time.Time `jsonapi:"attr,pre-plan-completed-at,rfc3339"`
	PrePlanRunningAt     time.Time `jsonapi:"attr,pre-plan-running-at,rfc3339"`
	QueuingAt            time.Time `jsonapi:"attr,queuing-at,rfc3339"`
}

// RunIncludeOpt represents the available options for include query params.
// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
type RunIncludeOpt string

const (
	RunPlan             RunIncludeOpt = "plan"
	RunApply            RunIncludeOpt = "apply"
	RunCreatedBy        RunIncludeOpt = "created_by"
	RunCostEstimate     RunIncludeOpt = "cost_estimate"
	RunConfigVer        RunIncludeOpt = "configuration_version"
	RunConfigVerIngress RunIncludeOpt = "configuration_version.ingress_attributes"
	RunWorkspace        RunIncludeOpt = "workspace"
	RunTaskStages       RunIncludeOpt = "task_stages"
)

// RunListOptions represents the options for listing runs.
type RunListOptions struct {
	ListOptions

	// Optional: Searches runs that matches the supplied VCS username.
	User string `url:"search[user],omitempty"`

	// Optional: Searches runs that matches the supplied commit sha.
	Commit string `url:"search[commit],omitempty"`

	// Optional: Searches runs that matches the supplied VCS username, commit sha, run_id, and run message.
	// The presence of search[commit] or search[user] takes priority over this parameter and will be omitted.
	Search string `url:"search[basic],omitempty"`

	// Optional: Current status of the run.
	// Options are listed at https://www.terraform.io/cloud-docs/api-docs/run#run-states
	Status string `url:"filter[status],omitempty"`

	// Optional: Source that triggered the run.
	// Options are listed at https://www.terraform.io/cloud-docs/api-docs/run#run-sources
	Source string `url:"filter[source],omitempty"`

	// Optional: Operation type for the run.
	// Options are listed at https://www.terraform.io/cloud-docs/api-docs/run#run-operations
	Operation string `url:"filter[operation],omitempty"`

	// Optional: A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include []RunIncludeOpt `url:"include,omitempty"`
}

// RunReadOptions represents the options for reading a run.
type RunReadOptions struct {
	// Optional: A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include []RunIncludeOpt `url:"include,omitempty"`
}

// RunCreateOptions represents the options for creating a new run.
type RunCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,runs"`

	// AllowEmptyApply specifies whether Terraform can apply the run even when the plan contains no changes.
	// Often used to upgrade state after upgrading a workspace to a new terraform version.
	AllowEmptyApply *bool `jsonapi:"attr,allow-empty-apply,omitempty"`

	// TerraformVersion specifies the Terraform version to use in this run.
	// Only valid for plan-only runs; must be a valid Terraform version available to the organization.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// PlanOnly specifies if this is a speculative, plan-only run that Terraform cannot apply.
	// Often used in conjunction with terraform-version in order to test whether an upgrade would succeed.
	PlanOnly *bool `jsonapi:"attr,plan-only,omitempty"`

	// Specifies if this plan is a destroy plan, which will destroy all
	// provisioned resources.
	IsDestroy *bool `jsonapi:"attr,is-destroy,omitempty"`

	// Refresh determines if the run should
	// update the state prior to checking for differences
	Refresh *bool `jsonapi:"attr,refresh,omitempty"`

	// RefreshOnly determines whether the run should ignore config changes
	// and refresh the state only
	RefreshOnly *bool `jsonapi:"attr,refresh-only,omitempty"`

	// Specifies the message to be associated with this run.
	Message *string `jsonapi:"attr,message,omitempty"`

	// Specifies the configuration version to use for this run. If the
	// configuration version object is omitted, the run will be created using the
	// workspace's latest configuration version.
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`

	// Specifies the workspace where the run will be executed.
	Workspace *Workspace `jsonapi:"relation,workspace"`

	// If non-empty, requests that Terraform should create a plan including
	// actions only for the given objects (specified using resource address
	// syntax) and the objects they depend on.
	//
	// This capability is provided for exceptional circumstances only, such as
	// recovering from mistakes or working around existing Terraform
	// limitations. Terraform will generally mention the -target command line
	// option in its error messages describing situations where setting this
	// argument may be appropriate. This argument should not be used as part
	// of routine workflow and Terraform will emit warnings reminding about
	// this whenever this property is set.
	TargetAddrs []string `jsonapi:"attr,target-addrs,omitempty"`

	// If non-empty, requests that Terraform create a plan that replaces
	// (destroys and then re-creates) the objects specified by the given
	// resource addresses.
	ReplaceAddrs []string `jsonapi:"attr,replace-addrs,omitempty"`

	// AutoApply determines if the run should be applied automatically without
	// user confirmation. It defaults to the Workspace.AutoApply setting.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// Variables allows you to specify terraform input variables for
	// a particular run, prioritized over variables defined on the workspace.
	Variables []*RunVariable `jsonapi:"attr,variables,omitempty"`
}

// RunApplyOptions represents the options for applying a run.
type RunApplyOptions struct {
	// An optional comment about the run.
	Comment *string `json:"comment,omitempty"`
}

// RunCancelOptions represents the options for canceling a run.
type RunCancelOptions struct {
	// An optional explanation for why the run was canceled.
	Comment *string `json:"comment,omitempty"`
}

type RunVariableAttr struct {
	Key   string `jsonapi:"attr,key"`
	Value string `jsonapi:"attr,value"`
}

// RunVariableAttr represents a variable that can be applied to a run. All values must be expressed as an HCL literal
// in the same syntax you would use when writing terraform code. See https://www.terraform.io/docs/language/expressions/types.html#types
// for more details.
type RunVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RunForceCancelOptions represents the options for force-canceling a run.
type RunForceCancelOptions struct {
	// An optional comment explaining the reason for the force-cancel.
	Comment *string `json:"comment,omitempty"`
}

// RunDiscardOptions represents the options for discarding a run.
type RunDiscardOptions struct {
	// An optional explanation for why the run was discarded.
	Comment *string `json:"comment,omitempty"`
}

// List all the runs of the given workspace.
func (s *runs) List(ctx context.Context, workspaceID string, options *RunListOptions) (*RunList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/runs", url.QueryEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	rl := &RunList{}
	err = req.Do(ctx, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// Create a new run with the given options.
func (s *runs) Create(ctx context.Context, options RunCreateOptions) (*Run, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", "runs", &options)
	if err != nil {
		return nil, err
	}

	r := &Run{}
	err = req.Do(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Read a run by its ID.
func (s *runs) Read(ctx context.Context, runID string) (*Run, error) {
	return s.ReadWithOptions(ctx, runID, nil)
}

// Read a run by its ID with the given options.
func (s *runs) ReadWithOptions(ctx context.Context, runID string, options *RunReadOptions) (*Run, error) {
	if !validStringID(&runID) {
		return nil, ErrInvalidRunID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("runs/%s", url.QueryEscape(runID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	r := &Run{}
	err = req.Do(ctx, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Apply a run by its ID.
func (s *runs) Apply(ctx context.Context, runID string, options RunApplyOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/apply", url.QueryEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Cancel a run by its ID.
func (s *runs) Cancel(ctx context.Context, runID string, options RunCancelOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/cancel", url.QueryEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// ForceCancel is used to forcefully cancel a run by its ID.
func (s *runs) ForceCancel(ctx context.Context, runID string, options RunForceCancelOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/force-cancel", url.QueryEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Discard a run by its ID.
func (s *runs) Discard(ctx context.Context, runID string, options RunDiscardOptions) error {
	if !validStringID(&runID) {
		return ErrInvalidRunID
	}

	u := fmt.Sprintf("runs/%s/actions/discard", url.QueryEscape(runID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o RunCreateOptions) valid() error {
	if o.Workspace == nil {
		return ErrRequiredWorkspace
	}

	if validString(o.TerraformVersion) && (o.PlanOnly == nil || !*o.PlanOnly) {
		return ErrTerraformVersionValidForPlanOnly
	}

	return nil
}

func (o *RunReadOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateRunIncludeParam(o.Include); err != nil {
		return err
	}
	return nil
}

func (o *RunListOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateRunIncludeParam(o.Include); err != nil {
		return err
	}
	return nil
}

func validateRunIncludeParam(params []RunIncludeOpt) error {
	for _, p := range params {
		switch p {
		case RunPlan, RunApply, RunCreatedBy, RunCostEstimate, RunConfigVer, RunConfigVerIngress, RunWorkspace, RunTaskStages:
			// do nothing
		default:
			return ErrInvalidIncludeValue
		}
	}

	return nil
}
