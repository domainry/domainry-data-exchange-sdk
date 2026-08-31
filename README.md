# Domainry Data Exchange SDK

This repository is the deployment-neutral contract between a Domainry Runtime and Data Exchange.

- `Binding` is consumed by Runtime application/HTTP layers.
- `modulehost` defines database, migration, workspace-context, import, and export Provider ports for the embedded Module.
- `saashost` defines the remote transport and explicit Runtime Provider bridge connection for SaaS deployment.

The SDK contains no Runtime domain models and no database/file implementation. Import Providers validate/apply authorized domain batches; Export Providers return authorized projected pages. File storage, parsing, checkpoints, leases, and artifacts remain Data Exchange implementation concerns.

## Package layout

- The root package is the stable Data Exchange `Factory`, `Binding`, request, and result contract.
- `modulehost` describes infrastructure and providers supplied to the embedded module.
- `saashost` describes authenticated remote transport and Runtime provider bridging.

The SDK intentionally has no public `persistence` package. Data Exchange stores, leases, checkpoints, and file DML remain source-owned implementation details.

Run `go test ./...` before publishing an immutable SDK version.
