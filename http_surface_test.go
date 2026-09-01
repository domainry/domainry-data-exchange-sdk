package dataexchange

import "testing"

func TestDataExchangeHTTPSurfaceOwnsJobManagementContract(t *testing.T) {
	contract := DataExchangeHTTPSurfaceContract()
	if contract.Owner != "data_exchange" || contract.ContractVersion != DataExchangeHTTPSurfaceContractVersion || len(contract.Routes) != 3 || len(contract.OpenAPI) != 3 {
		t.Fatalf("Data Exchange HTTP contract=%+v", contract)
	}
	for _, route := range contract.Routes {
		if contract.OpenAPI[route.Pattern]["operationId"] == nil {
			t.Fatalf("route %q has no OpenAPI operation", route.Pattern)
		}
	}
}
