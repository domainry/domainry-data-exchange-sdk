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

type ExportProvider interface {
	ReadExportPage(context.Context, dataexchange.ExportPageRequest) (dataexchange.ExportPage, error)
}
