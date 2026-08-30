package modulehost

import (
	"context"
	"database/sql"

	dataexchange "github.com/domainry/domainry-data-exchange-sdk"
)

type Host interface {
	ImportProvider(string) (ImportProvider, bool)
	ExportProvider(string) (ExportProvider, bool)
}

type ModuleHost interface {
	Host
	Database() *sql.DB
	Migrations() MigrationRegistrar
	WorkspaceContext(context.Context, string, string) context.Context
}

type Migration struct {
	ID, SQL string
}

type MigrationRegistrar interface {
	Driver() string
	Schema() string
	ApplyOwnedMigrations(context.Context, string, []Migration) error
}

type ImportProvider interface {
	ValidateImportBatch(context.Context, dataexchange.ImportBatch) (dataexchange.ImportBatchResult, error)
	ApplyImportBatch(context.Context, dataexchange.ImportBatch) (dataexchange.ImportBatchResult, error)
}

// ImportArtifactProvider is an optional atomic-bundle capability. The engine
// opens the durable source independently for validation and apply, so both
// calls receive a fresh reader and remain replay-safe.
type ImportArtifactProvider interface {
	ValidateImportArtifact(context.Context, dataexchange.ImportArtifact) (dataexchange.ImportArtifactResult, error)
	ApplyImportArtifact(context.Context, dataexchange.ImportArtifact) (dataexchange.ImportArtifactResult, error)
}

type ExportProvider interface {
	ReadExportPage(context.Context, dataexchange.ExportPageRequest) (dataexchange.ExportPage, error)
}

// ExportArtifactProvider is an optional canonical-bundle capability for
// formats that cannot be represented as paged CSV without losing their
// domain integrity contract.
type ExportArtifactProvider interface {
	BuildExportArtifact(context.Context, dataexchange.ExportArtifactRequest) (dataexchange.ExportArtifact, error)
}

// ExportPlanningProvider is an optional capability for providers that own
// governed artifact naming or expiry semantics.
type ExportPlanningProvider interface {
	PlanExport(context.Context, dataexchange.ExportPlanRequest) (dataexchange.ExportPlan, error)
}

// ExportCompletionProvider is an optional capability for replay-safe business
// projection and audit finalization. The Data Exchange engine remains the only
// owner of file content and artifact persistence.
type ExportCompletionProvider interface {
	CompleteExport(context.Context, dataexchange.ExportCompletion) error
}
