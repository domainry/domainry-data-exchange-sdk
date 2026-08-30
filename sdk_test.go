package dataexchange

import "testing"

func TestIdentityAndDescriptorContractsFailClosed(t *testing.T) {
	if (ApplicationRef{}).Validate() == nil || (Scope{}).Validate() == nil || (Descriptor{}).Validate() == nil {
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
}
