package github

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestNoizuProviderValid asserts the provider builds with the Noizu
// extension layer merged in and passes InternalValidate. Purely
// in-memory: no network and no credentials (InternalValidate does not
// invoke ConfigureContextFunc).
func TestNoizuProviderValid(t *testing.T) {
	p := NewProvider("test", "test")()
	if p == nil {
		t.Fatal("NewProvider returned nil provider")
	}
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("provider fails InternalValidate: %s", err)
	}
}

// TestNoizuExtensionNaming enforces the github_noizu_ prefix convention
// for every Noizu-registered resource and data source.
func TestNoizuExtensionNaming(t *testing.T) {
	for name := range noizuResources {
		if !strings.HasPrefix(name, "github_noizu_") {
			t.Errorf("noizu resource %q violates naming convention: must start with github_noizu_", name)
		}
	}
	for name := range noizuDataSources {
		if !strings.HasPrefix(name, "github_noizu_") {
			t.Errorf("noizu data source %q violates naming convention: must start with github_noizu_", name)
		}
	}
}

// TestNoizuMergePanicsOnCollision guards the collision detector: merging
// a name that already exists upstream must panic.
func TestNoizuMergePanicsOnCollision(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on upstream name collision")
		}
	}()
	dst := map[string]*schema.Resource{"github_upstream_thing": {Schema: map[string]*schema.Schema{}}}
	src := map[string]*schema.Resource{"github_upstream_thing": {Schema: map[string]*schema.Schema{}}}
	mergeNoizuMaps("resource", dst, src)
}

// TestNoizuMergeNoCollisionPassthrough verifies non-colliding merges land
// in the destination map.
func TestNoizuMergeNoCollisionPassthrough(t *testing.T) {
	dst := map[string]*schema.Resource{"github_upstream_thing": {Schema: map[string]*schema.Schema{}}}
	src := map[string]*schema.Resource{"github_noizu_thing": {Schema: map[string]*schema.Schema{}}}
	mergeNoizuMaps("resource", dst, src)
	if _, ok := dst["github_noizu_thing"]; !ok {
		t.Fatal("noizu resource was not merged into the provider map")
	}
}
