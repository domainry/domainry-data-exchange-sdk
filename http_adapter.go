package dataexchange

import (
	"strings"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const DataExchangeHTTPAdapterContractVersion = "domainry-data-exchange-http-adapter-v1"

const (
	ActionDataExchangeJobList     = "data_exchange.jobs.list"
	ActionDataExchangeJobGet      = "data_exchange.jobs.get"
	ActionDataExchangeJobCancel   = "data_exchange.jobs.cancel"
	ActionDataExchangeJobDownload = "data_exchange.jobs.download"
)

type HTTPRouteContract struct {
	Action actioncontract.ActionDefinition `json:"action"`
}

func (route HTTPRouteContract) Pattern() string {
	if route.Action.HTTP == nil {
		return ""
	}
	return route.Action.HTTP.Method + " " + route.Action.HTTP.RouteTemplate
}

type HTTPAdapterContract struct {
	ContractVersion string              `json:"contract_version"`
	Owner           string              `json:"owner"`
	Name            string              `json:"name"`
	Routes          []HTTPRouteContract `json:"routes"`
}

func DataExchangeHTTPAdapterContract() HTTPAdapterContract {
	routes := []HTTPRouteContract{
		jobHTTPRoute(ActionDataExchangeJobList, "list", "GET /data-exchange/jobs", "List Data Exchange jobs", "read", "not_applicable", "owner_read_audit_policy"),
		jobHTTPRoute(ActionDataExchangeJobGet, "get", "GET /data-exchange/jobs/{jobID}", "Get Data Exchange job", "read", "not_applicable", "owner_read_audit_policy"),
		jobHTTPRoute(ActionDataExchangeJobCancel, "cancel", "POST /data-exchange/jobs/{jobID}/cancel", "Cancel Data Exchange job", "write", "natural_key", "mutation_audit_required"),
		jobHTTPRoute(ActionDataExchangeJobDownload, "download", "GET /data-exchange/jobs/{jobID}/download", "Download Data Exchange job", "read", "not_applicable", "business_export_download_audit"),
	}
	return HTTPAdapterContract{
		ContractVersion: DataExchangeHTTPAdapterContractVersion, Owner: "data_exchange", Name: "job_management", Routes: routes,
	}
}

func jobHTTPRoute(key, operation, pattern, label, effect, idempotency, audit string) HTTPRouteContract {
	method, route, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	separator := strings.LastIndexByte(key, '.')
	resourceKey := key[:separator]
	risk := actioncontract.RiskLow
	if actioncontract.EffectClass(effect) == actioncontract.EffectWrite {
		risk = actioncontract.RiskMedium
	}
	definition, err := actioncontract.NormalizeDefinition(actioncontract.ActionDefinition{
		Key: key, Owner: "module:data_exchange", SourceKind: "module_http",
		CapabilityKey: "data_exchange.jobs", CapabilityLabel: "Data Exchange jobs",
		OperationKey: operation, OperationLabel: label, Label: label,
		Exposures:     []actioncontract.Exposure{actioncontract.ExposurePublic},
		Authorization: actioncontract.Authorization{Strategy: actioncontract.AuthorizationAuthenticated},
		HTTP:          &actioncontract.HTTPBinding{Method: method, RouteTemplate: route},
		Permission: &actioncontract.PermissionDefinition{
			Key: key, Owner: "module:data_exchange", ResourceKey: resourceKey, OperationKey: operation,
			Label: label, Category: "Data Exchange jobs", LifecycleStatus: actioncontract.LifecycleActive,
		},
		EffectClass: actioncontract.EffectClass(effect), RiskLevel: risk,
		IdempotencyDecision: idempotency, AuditClass: audit, LifecycleStatus: actioncontract.LifecycleActive,
	})
	if err != nil {
		panic("invalid Data Exchange HTTP Action: " + err.Error())
	}
	return HTTPRouteContract{Action: definition}
}
