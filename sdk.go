// Package dataexchange defines the deployment-neutral contract shared by the
// in-process Data Exchange Module and the remote SaaS binding.
package dataexchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/domainry/domainry-foundation/modulecapability"
)

const ProtocolVersionV1 = "data-exchange.v1"

var ErrSourceTooLarge = errors.New("Data Exchange source exceeds configured limit")
var ErrSourceUnreadable = errors.New("Data Exchange source cannot be read")
var ErrIdempotencyKeyReused = errors.New("Data Exchange idempotency key reused")
var ErrJobNotFound = errors.New("Data Exchange job was not found")
var ErrArtifactExpired = errors.New("Data Exchange artifact has expired")
var ErrContentCorrupt = errors.New("Data Exchange durable content failed integrity verification")

type DeploymentMode string

const (
	DeploymentModeModule DeploymentMode = "module"
	DeploymentModeSaaS   DeploymentMode = "saas"
)

type ApplicationRef struct {
	ApplicationID string `json:"application_id"`
	RuntimeID     string `json:"runtime_id"`
}

func (r ApplicationRef) Validate() error {
	if strings.TrimSpace(r.ApplicationID) == "" || strings.TrimSpace(r.RuntimeID) == "" {
		return fmt.Errorf("Data Exchange application identity is incomplete")
	}
	return nil
}

type Descriptor struct {
	ProtocolVersion string         `json:"protocol_version"`
	Mode            DeploymentMode `json:"mode"`
	Capabilities    []string       `json:"capabilities"`
}

func (d Descriptor) Validate() error {
	if d.ProtocolVersion != ProtocolVersionV1 || d.Mode != DeploymentModeModule && d.Mode != DeploymentModeSaaS {
		return fmt.Errorf("Data Exchange descriptor is incompatible")
	}
	return nil
}

type Scope struct {
	WorkspaceID string `json:"workspace_id"`
	ActorID     string `json:"actor_id"`
	RoleKey     string `json:"role_key,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
}

func (s Scope) Validate() error {
	if strings.TrimSpace(s.WorkspaceID) == "" || strings.TrimSpace(s.ActorID) == "" {
		return fmt.Errorf("Data Exchange workspace and actor are required")
	}
	return nil
}

type Job struct {
	ID, Provider, Operation, Status string
	WorkspaceID, ObjectKey          string
	ActorID, RoleKey                string
	ReferenceID                     string
	// Options is the immutable provider-owned request payload. It is returned
	// only through the owning Binding so the application can project and
	// authorize its own job without duplicating queue state.
	Options               []byte
	Checkpoint, Total     int
	Cursor                string
	ResultChunks          int
	FencingToken          int64
	LeaseOwner            string
	ArtifactID, ErrorCode string
	CreatedAt, UpdatedAt  time.Time
}

type ImportRequest struct {
	Scope          Scope
	Provider       string
	ObjectKey      string
	IdempotencyKey string
	Filename       string
	ContentType    string
	Options        []byte
	Source         io.Reader
	MaxBytes       int64
}

type ExportRequest struct {
	Scope          Scope
	Provider       string
	ObjectKey      string
	IdempotencyKey string
	// ReferenceID binds the winning owner request (for example, an audit row)
	// without changing the canonical export condition used for idempotency.
	ReferenceID string
	Options     []byte
}

type JobRequest struct {
	Scope Scope
	JobID string
	// Provider and Operation are optional ownership constraints. When supplied,
	// Job, Cancel, and Download must match them in the same persistence
	// operation; callers must not authorize a job with a separate read first.
	Provider  string
	Operation string
}

// JobListRequest lists only jobs owned by the authenticated actor in one
// workspace. Provider, operation, and status are optional server-side filters;
// Limit is bounded so the HTTP route cannot become an unbounded export path.
type JobListRequest struct {
	Scope     Scope
	Provider  string
	Operation string
	Status    string
	Limit     int
}

func (r JobListRequest) Validate() error {
	if err := r.Scope.Validate(); err != nil {
		return err
	}
	operation := strings.TrimSpace(r.Operation)
	if operation != "" && operation != "import" && operation != "export" {
		return fmt.Errorf("Data Exchange job operation is invalid")
	}
	status := strings.TrimSpace(r.Status)
	if status != "" && status != "queued" && status != "running" && status != "completed" && status != "failed" && status != "cancelled" {
		return fmt.Errorf("Data Exchange job status is invalid")
	}
	if r.Limit < 0 || r.Limit > 200 {
		return fmt.Errorf("Data Exchange job list limit is invalid")
	}
	return nil
}

func (r JobRequest) Validate() error {
	if err := r.Scope.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.JobID) == "" {
		return fmt.Errorf("Data Exchange job ID is required")
	}
	operation := strings.TrimSpace(r.Operation)
	if operation != "" && operation != "import" && operation != "export" {
		return fmt.Errorf("Data Exchange job operation is invalid")
	}
	return nil
}

type Artifact struct {
	ID, Filename, ContentType, SHA256 string
	Size                              int64
	ExpiresAt                         time.Time
	Content                           io.ReadCloser
}

type WorkerConfig struct {
	Enabled      bool
	PollInterval time.Duration
	BatchSize    int
	LeaseTTL     time.Duration
}

type Factory interface {
	Open(context.Context, ApplicationRef) (Binding, error)
}

type Binding interface {
	modulecapability.Binding
	Descriptor() Descriptor
	SubmitImport(context.Context, ImportRequest) (Job, bool, error)
	SubmitExport(context.Context, ExportRequest) (Job, bool, error)
	Job(context.Context, JobRequest) (Job, error)
	Cancel(context.Context, JobRequest) (Job, error)
	Download(context.Context, JobRequest) (Artifact, error)
	Start(context.Context, WorkerConfig) <-chan struct{}
	Close(context.Context) error
}

// SubjectLifecycleBinding is a privileged owner capability. It is deliberately
// separate from user job management and is never mounted on the job HTTP routes.
type SubjectLifecycleBinding interface {
	SubjectLifecycle() SubjectLifecycle
}

type SubjectErasureRequest struct {
	WorkspaceID string          `json:"workspace_id"`
	SubjectID   string          `json:"subject_id"`
	RequestID   string          `json:"request_id"`
	LegalHolds  json.RawMessage `json:"legal_holds,omitempty"`
}

type SubjectLifecycle interface {
	PreviewSubject(context.Context, string, string) (json.RawMessage, error)
	ExportSubject(context.Context, string, string) (json.RawMessage, error)
	PrepareSubjectErasure(context.Context, SubjectErasureRequest) (json.RawMessage, error)
	ErasePreparedSubject(context.Context, SubjectErasureRequest, json.RawMessage) (json.RawMessage, error)
}

type ImportBatch struct {
	Scope                     Scope
	ObjectKey, JobID, ChunkID string
	// Attempt is one-based and changes whenever the durable worker retries the
	// job. Final marks the last validation/apply batch in that pass so providers
	// can release attempt-scoped state without guessing from chunk IDs.
	Attempt int
	Final   bool
	Headers []string
	Rows    []ImportRow
}
type ImportRow struct {
	Number int
	Values []string
}
type ImportBatchResult struct {
	Accepted, Rejected int
	Receipt            string
}

// ImportArtifact is the replayable, provider-owned alternative to row-batch
// import. It is intended for canonical bundles whose integrity and atomicity
// would be lost if the engine decoded them as CSV rows.
type ImportArtifact struct {
	Scope                 Scope
	ObjectKey, JobID      string
	Filename, ContentType string
	Options               []byte
	Content               io.Reader
}

type ImportArtifactResult struct {
	Records int
	Receipt string
}

type ExportPageRequest struct {
	Scope             Scope
	ObjectKey         string
	ReferenceID       string
	Options           []byte
	Cursor            string
	PageSize          int
	JobID             string
	ArtifactExpiresAt time.Time
}
type ExportPage struct {
	Columns    []string
	Rows       [][]string
	NextCursor string
	Total      int
}

// ExportPlan fixes artifact identity before the first page is read. Providers
// may choose a governed filename and expiry; the engine supplies safe defaults
// when the optional planning capability is not implemented.
type ExportPlanRequest struct {
	Scope       Scope
	ObjectKey   string
	ReferenceID string
	Options     []byte
	JobID       string
	CreatedAt   time.Time
}

type ExportPlan struct {
	Filename    string
	ContentType string
	ExpiresAt   time.Time
}

// ExportArtifactRequest asks a provider to produce one canonical artifact.
// The Data Exchange engine remains responsible for durable storage, leases,
// retries, artifact identity and lifecycle.
type ExportArtifactRequest struct {
	Scope       Scope
	ObjectKey   string
	ReferenceID string
	Options     []byte
	JobID       string
	CreatedAt   time.Time
}

type ExportArtifact struct {
	Filename    string
	ContentType string
	ExpiresAt   time.Time
	Content     io.ReadCloser
	Records     int
}

// ExportCompletion is delivered after all result chunks have a stable
// identity and before the job becomes completed. Provider finalization must be
// replay-safe because a lost lease or failed terminal transition can retry it.
type ExportCompletion struct {
	Scope        Scope
	ObjectKey    string
	ReferenceID  string
	Options      []byte
	JobID        string
	Artifact     Artifact
	Rows         int
	ResultChunks int
}
