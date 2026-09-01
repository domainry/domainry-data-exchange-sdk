package dataexchange

import (
	"testing"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

func TestDataExchangeHTTPSurfaceOwnsJobManagementContract(t *testing.T) {
	contract := DataExchangeHTTPSurfaceContract()
	operations := contract.OpenAPIOperations()
	if contract.Owner != "data_exchange" || contract.ContractVersion != DataExchangeHTTPSurfaceContractVersion || len(contract.Routes) != 3 || len(operations) != 3 {
		t.Fatalf("Data Exchange HTTP contract=%+v", contract)
	}
	wantKeys := []string{ActionDataExchangeJobGet, ActionDataExchangeJobCancel, ActionDataExchangeJobDownload}
	for index, route := range contract.Routes {
		if route.Action.Key != wantKeys[index] || operations[route.Pattern()]["operationId"] == nil || route.Action.Authorization.Strategy != actioncontract.AuthorizationAuthenticatedPrincipal || route.Action.Permission != nil {
			t.Fatalf("route %q has incomplete Action/OpenAPI ownership", route.Pattern())
		}
	}
}
