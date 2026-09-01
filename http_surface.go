package dataexchange

const DataExchangeHTTPSurfaceContractVersion = "domainry-data-exchange-http-surface-v1"

type HTTPRouteContract struct {
	Pattern             string   `json:"pattern"`
	Exposures           []string `json:"exposures"`
	Authentication      string   `json:"authentication"`
	Permission          string   `json:"permission,omitempty"`
	AnyPermissions      []string `json:"any_permissions,omitempty"`
	PrincipalOnly       bool     `json:"principal_only,omitempty"`
	EffectClass         string   `json:"effect_class"`
	HighRiskPolicy      string   `json:"high_risk_policy"`
	IdempotencyDecision string   `json:"idempotency_decision"`
	AuditClass          string   `json:"audit_class"`
}

type HTTPSurfaceContract struct {
	ContractVersion string                    `json:"contract_version"`
	Owner           string                    `json:"owner"`
	Name            string                    `json:"name"`
	Routes          []HTTPRouteContract       `json:"routes"`
	OpenAPI         map[string]map[string]any `json:"openapi_operations"`
}

func DataExchangeHTTPSurfaceContract() HTTPSurfaceContract {
	routes := []HTTPRouteContract{
		jobHTTPRoute("GET /data-exchange/jobs/{jobID}", "read", "not_applicable", "owner_read_audit_policy"),
		jobHTTPRoute("POST /data-exchange/jobs/{jobID}/cancel", "write", "natural_key", "mutation_audit_required"),
		jobHTTPRoute("GET /data-exchange/jobs/{jobID}/download", "read", "not_applicable", "business_export_download_audit"),
	}
	return HTTPSurfaceContract{
		ContractVersion: DataExchangeHTTPSurfaceContractVersion, Owner: "data_exchange", Name: "job_management", Routes: routes,
		OpenAPI: map[string]map[string]any{
			"GET /data-exchange/jobs/{jobID}":          dataExchangeJobOperation("getDataExchangeJob", "Get an actor-owned Data Exchange job"),
			"POST /data-exchange/jobs/{jobID}/cancel":  dataExchangeJobOperation("cancelDataExchangeJob", "Cancel an actor-owned Data Exchange job"),
			"GET /data-exchange/jobs/{jobID}/download": dataExchangeDownloadOperation(),
		},
	}
}

func jobHTTPRoute(pattern, effect, idempotency, audit string) HTTPRouteContract {
	return HTTPRouteContract{
		Pattern: pattern, Exposures: []string{"public"}, Authentication: "authenticated", PrincipalOnly: true,
		EffectClass: effect, HighRiskPolicy: "none", IdempotencyDecision: idempotency, AuditClass: audit,
	}
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
