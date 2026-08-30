// Package dataexchange defines the deployment-neutral contract shared by the
// in-process Data Exchange Module and the remote SaaS binding.
package dataexchange

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const ProtocolVersionV1 = "data-exchange.v1"

var ErrSourceTooLarge = errors.New("Data Exchange source exceeds configured limit")
var ErrSourceUnreadable = errors.New("Data Exchange source cannot be read")

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
	Checkpoint, Total               int
	Cursor                          string
	ResultChunks                    int
	FencingToken                    int64
	LeaseOwner                      string
	ArtifactID, ErrorCode           string
	CreatedAt, UpdatedAt            time.Time
}

type ImportRequest struct {
	Scope          Scope
	Provider       string
	ObjectKey      string
	IdempotencyKey string
	Filename       string
	ContentType    string
	Source         io.Reader
	MaxBytes       int64
}

type ExportRequest struct {
	Scope          Scope
	Provider       string
	ObjectKey      string
	IdempotencyKey string
	Options        []byte
}

type JobRequest struct {
	Scope Scope
	JobID string
}

type Artifact struct {
	ID, Filename, ContentType, SHA256 string
	Size                              int64
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
	Descriptor() Descriptor
	SubmitImport(context.Context, ImportRequest) (Job, bool, error)
	SubmitExport(context.Context, ExportRequest) (Job, bool, error)
	Job(context.Context, JobRequest) (Job, error)
	Cancel(context.Context, JobRequest) (Job, error)
	Download(context.Context, JobRequest) (Artifact, error)
	Start(context.Context, WorkerConfig) <-chan struct{}
	Close(context.Context) error
}

type ImportBatch struct {
	Scope                     Scope
	ObjectKey, JobID, ChunkID string
	Headers                   []string
	Rows                      []ImportRow
}
type ImportRow struct {
	Number int
	Values []string
}
type ImportBatchResult struct {
	Accepted, Rejected int
	Receipt            string
}

type ExportPageRequest struct {
	Scope     Scope
	ObjectKey string
	Options   []byte
	Cursor    string
	PageSize  int
	JobID     string
}
type ExportPage struct {
	Columns    []string
	Rows       [][]string
	NextCursor string
	Total      int
}
