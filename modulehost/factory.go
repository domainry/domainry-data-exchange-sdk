package modulehost

import (
	"context"

	dataexchange "github.com/domainry/domainry-data-exchange-sdk"
)

type Factory interface {
	OpenModule(context.Context, dataexchange.ApplicationRef, ModuleHost) (dataexchange.Binding, error)
}
