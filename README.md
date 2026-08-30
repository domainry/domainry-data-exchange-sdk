# Domainry Data Exchange SDK

This repository is the deployment-neutral contract between a Domainry Runtime and Data Exchange.

- `Binding` is consumed by Runtime application/HTTP layers.
- `modulehost` defines database, migration, workspace-context, import, and export Provider ports for the embedded Module.
- `saashost` defines the remote transport and explicit Runtime Provider bridge connection for SaaS deployment.

The SDK contains no Runtime domain models and no database/file implementation. Import Providers validate/apply authorized domain batches; Export Providers return authorized projected pages. File storage, parsing, checkpoints, leases, and artifacts remain Data Exchange implementation concerns.
