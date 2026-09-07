package dataexchange

import (
	"errors"
	"fmt"
	"testing"
)

func TestIdempotencyKeyReusedErrorSupportsErrorsIs(t *testing.T) {
	err := fmt.Errorf("submit Data Exchange export: %w", ErrIdempotencyKeyReused)
	if !errors.Is(err, ErrIdempotencyKeyReused) {
		t.Fatalf("wrapped error %v does not match ErrIdempotencyKeyReused", err)
	}
}

func TestIdentityAndDescriptorContractsFailClosed(t *testing.T) {
	if (ApplicationRef{}).Validate() == nil || (Scope{}).Validate() == nil || (Descriptor{}).Validate() == nil || (JobRequest{}).Validate() == nil {
		t.Fatal("empty Data Exchange contract was accepted")
	}
	if err := (ApplicationRef{ApplicationID: "app", RuntimeID: "runtime"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Scope{WorkspaceID: "workspace", ActorID: "actor"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Descriptor{ProtocolVersion: ProtocolVersionV1, Mode: DeploymentModeModule}).Validate(); err != nil {
		t.Fatal(err)
	}
	request := JobRequest{Scope: Scope{WorkspaceID: "workspace", ActorID: "actor"}, JobID: "job", Provider: "records", Operation: "export"}
	if err := request.Validate(); err != nil {
		t.Fatal(err)
	}
	request.Operation = "unknown"
	if request.Validate() == nil {
		t.Fatal("unknown Data Exchange operation was accepted")
	}
}
