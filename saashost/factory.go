package saashost

import (
	"context"
	"encoding/json"

	dataexchange "github.com/domainry/domainry-data-exchange-sdk"
	"github.com/domainry/domainry-data-exchange-sdk/modulehost"
)

// Transport carries the deployment-neutral Data Exchange contract to a SaaS
// owner. Cross-process implementations must decode an idempotency-key reuse
// response into dataexchange.ErrIdempotencyKeyReused before returning from
// SubmitImport or SubmitExport.
type Transport interface {
	Connect(context.Context, dataexchange.ApplicationRef, modulehost.Host) error
	Descriptor(context.Context, dataexchange.ApplicationRef) (dataexchange.Descriptor, error)
	SubmitImport(context.Context, dataexchange.ApplicationRef, dataexchange.ImportRequest) (dataexchange.Job, bool, error)
	SubmitExport(context.Context, dataexchange.ApplicationRef, dataexchange.ExportRequest) (dataexchange.Job, bool, error)
	Job(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Job, error)
	Cancel(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Job, error)
	Download(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Artifact, error)
	Close(context.Context, dataexchange.ApplicationRef) error
}

// SubjectLifecycleTransport is the authenticated owner-to-owner transport. A
// deployment must provide it before it can participate in account erasure.
type SubjectLifecycleTransport interface {
	PreviewSubject(context.Context, dataexchange.ApplicationRef, string, string) (json.RawMessage, error)
	ExportSubject(context.Context, dataexchange.ApplicationRef, string, string) (json.RawMessage, error)
	PrepareSubjectErasure(context.Context, dataexchange.ApplicationRef, dataexchange.SubjectErasureRequest) (json.RawMessage, error)
	ErasePreparedSubject(context.Context, dataexchange.ApplicationRef, dataexchange.SubjectErasureRequest, json.RawMessage) (json.RawMessage, error)
}

type Factory interface {
	OpenSaaS(context.Context, dataexchange.ApplicationRef, modulehost.Host) (dataexchange.Binding, error)
}
