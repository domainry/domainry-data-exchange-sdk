package dataexchange

import (
	"testing"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

func TestDataExchangeHTTPAdapterOwnsJobManagementContract(t *testing.T) {
	contract := DataExchangeHTTPAdapterContract()
	operations := contract.OpenAPIOperations()
	if contract.Owner != "data_exchange" || contract.ContractVersion != DataExchangeHTTPAdapterContractVersion || len(contract.Routes) != 3 || len(operations) != 3 {
		t.Fatalf("Data Exchange HTTP contract=%+v", contract)
	}
	wantKeys := []string{ActionDataExchangeJobGet, ActionDataExchangeJobCancel, ActionDataExchangeJobDownload}
	for index, route := range contract.Routes {
		if route.Action.Key != wantKeys[index] || operations[route.Pattern()]["operationId"] == nil || route.Action.Authorization.Strategy != actioncontract.AuthorizationAuthenticated || route.Action.Permission != nil {
			t.Fatalf("route %q has incomplete Action/OpenAPI ownership", route.Pattern())
		}
	}
}
