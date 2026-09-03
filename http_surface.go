package dataexchange

import (
	"strings"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const DataExchangeHTTPSurfaceContractVersion = "domainry-data-exchange-http-surface-v2"

const (
	ActionDataExchangeJobGet      = "data_exchange.jobs.get"
	ActionDataExchangeJobCancel   = "data_exchange.jobs.cancel"
	ActionDataExchangeJobDownload = "data_exchange.jobs.download"
)

type HTTPRouteContract struct {
	Action           actioncontract.ActionDefinition `json:"action"`
	OpenAPIOperation map[string]any                  `json:"openapi_operation"`
}

func (route HTTPRouteContract) Pattern() string {
	if route.Action.HTTP == nil {
		return ""
	}
	return route.Action.HTTP.Method + " " + route.Action.HTTP.RouteTemplate
}

type HTTPSurfaceContract struct {
	ContractVersion string              `json:"contract_version"`
	Owner           string              `json:"owner"`
	Name            string              `json:"name"`
	Routes          []HTTPRouteContract `json:"routes"`
}

// OpenAPIOperations projects Method+URL keys from the same Action entries used
// for route mounting. The contract deliberately stores no second URL-keyed
// OpenAPI inventory.
func (contract HTTPSurfaceContract) OpenAPIOperations() map[string]map[string]any {
	operations := make(map[string]map[string]any, len(contract.Routes))
	for _, route := range contract.Routes {
		operations[route.Pattern()] = route.OpenAPIOperation
	}
	return operations
}

func DataExchangeHTTPSurfaceContract() HTTPSurfaceContract {
	routes := []HTTPRouteContract{
		jobHTTPRoute(ActionDataExchangeJobGet, "get", "GET /data-exchange/jobs/{jobID}", "Get Data Exchange job", "read", "not_applicable", "owner_read_audit_policy", dataExchangeJobOperation("getDataExchangeJob", "Get an actor-owned Data Exchange job")),
		jobHTTPRoute(ActionDataExchangeJobCancel, "cancel", "POST /data-exchange/jobs/{jobID}/cancel", "Cancel Data Exchange job", "write", "natural_key", "mutation_audit_required", dataExchangeJobOperation("cancelDataExchangeJob", "Cancel an actor-owned Data Exchange job")),
		jobHTTPRoute(ActionDataExchangeJobDownload, "download", "GET /data-exchange/jobs/{jobID}/download", "Download Data Exchange job", "read", "not_applicable", "business_export_download_audit", dataExchangeDownloadOperation()),
	}
	return HTTPSurfaceContract{
		ContractVersion: DataExchangeHTTPSurfaceContractVersion, Owner: "data_exchange", Name: "job_management", Routes: routes,
	}
}

func jobHTTPRoute(key, operation, pattern, label, effect, idempotency, audit string, openAPI map[string]any) HTTPRouteContract {
	method, route, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	risk := actioncontract.RiskLow
	if actioncontract.EffectClass(effect) == actioncontract.EffectWrite {
		risk = actioncontract.RiskMedium
	}
	definition, err := actioncontract.NormalizeDefinition(actioncontract.ActionDefinition{
		Key: key, Owner: "module:data_exchange", SourceKind: "module_surface",
		CapabilityKey: "data_exchange.jobs", CapabilityLabel: "Data Exchange jobs",
		OperationKey: operation, OperationLabel: label, Label: label,
		Exposures:     []actioncontract.Exposure{actioncontract.ExposurePublic},
		Authorization: actioncontract.Authorization{Strategy: actioncontract.AuthorizationAuthenticated},
		HTTP:          &actioncontract.HTTPBinding{Method: method, RouteTemplate: route},
		EffectClass:   actioncontract.EffectClass(effect), RiskLevel: risk,
		IdempotencyDecision: idempotency, AuditClass: audit, LifecycleStatus: actioncontract.LifecycleActive,
	})
	if err != nil {
		panic("invalid Data Exchange HTTP Action: " + err.Error())
	}
	return HTTPRouteContract{Action: definition, OpenAPIOperation: openAPI}
}

func dataExchangeDownloadOperation() map[string]any {
	operation := dataExchangeJobOperation("downloadDataExchangeJob", "Download an actor-owned completed Data Exchange export")
	operation["responses"] = map[string]any{
		"200": map[string]any{"description": "Export artifact", "content": map[string]any{"application/octet-stream": map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}},
		"404": map[string]any{"description": "Job not found"}, "409": map[string]any{"description": "Artifact unavailable"},
	}
	return operation
}

func dataExchangeJobOperation(operationID, summary string) map[string]any {
	return map[string]any{
		"operationId": operationID, "summary": summary, "tags": []string{"Data Exchange"},
		"security": []map[string]any{{"BearerAuth": []string{}}},
		"parameters": []map[string]any{
			{"name": "jobID", "in": "path", "required": true, "schema": map[string]any{"type": "string", "minLength": 1}},
			{"name": "provider", "in": "query", "required": false, "schema": map[string]any{"type": "string", "minLength": 1}},
			{"name": "operation", "in": "query", "required": false, "schema": map[string]any{"type": "string", "enum": []string{"import", "export"}}},
		},
		"responses": map[string]any{
			"200": map[string]any{"description": "Provider-owned Data Exchange job projection", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "additionalProperties": true}}}},
			"404": map[string]any{"description": "Job not found"},
		},
	}
}
