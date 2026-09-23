package dataexchange

import (
	"testing"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

func TestDataExchangeHTTPAdapterOwnsJobManagementContract(t *testing.T) {
	contract := DataExchangeHTTPAdapterContract()
	if contract.Owner != "data_exchange" || contract.ContractVersion != DataExchangeHTTPAdapterContractVersion || len(contract.Routes) != 4 {
		t.Fatalf("Data Exchange HTTP contract=%+v", contract)
	}
	wantKeys := []string{ActionDataExchangeJobList, ActionDataExchangeJobGet, ActionDataExchangeJobCancel, ActionDataExchangeJobDownload}
	for index, route := range contract.Routes {
		permission := route.Action.Permission
		if route.Action.Key != wantKeys[index] || route.Action.Authorization.Strategy != actioncontract.AuthorizationAuthenticated || permission == nil || permission.Key != route.Action.Key || permission.Owner != "module:data_exchange" {
			t.Fatalf("route %q has incomplete Action ownership", route.Pattern())
		}
	}
}
