package saashost

import (
	"context"

	dataexchange "github.com/domainry/domainry-data-exchange-sdk"
	"github.com/domainry/domainry-data-exchange-sdk/modulehost"
	"github.com/domainry/domainry-foundation/modulecapability"
)

// Transport carries the deployment-neutral Data Exchange contract to a SaaS
// owner. Cross-process implementations must decode an idempotency-key reuse
// response into dataexchange.ErrIdempotencyKeyReused before returning from
// SubmitImport or SubmitExport.
type Transport interface {
	modulecapability.Binding
	Connect(context.Context, dataexchange.ApplicationRef, modulehost.Host) error
	Descriptor(context.Context, dataexchange.ApplicationRef) (dataexchange.Descriptor, error)
	SubmitImport(context.Context, dataexchange.ApplicationRef, dataexchange.ImportRequest) (dataexchange.Job, bool, error)
	SubmitExport(context.Context, dataexchange.ApplicationRef, dataexchange.ExportRequest) (dataexchange.Job, bool, error)
	Job(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Job, error)
	Cancel(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Job, error)
	Download(context.Context, dataexchange.ApplicationRef, dataexchange.JobRequest) (dataexchange.Artifact, error)
	Close(context.Context, dataexchange.ApplicationRef) error
}

type Factory interface {
	OpenSaaS(context.Context, dataexchange.ApplicationRef, modulehost.Host) (dataexchange.Binding, error)
}
