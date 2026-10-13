// Noizu extension layer for the noizu-forks fork of
// integrations/terraform-provider-github.
//
// Everything Noizu adds lives in noizu_*.go files in this package (same
// `github` package on purpose: it avoids import cycles with the provider's
// Config/Owner helpers). Upstream files receive exactly ONE touch: the
// return expression of NewProvider in github/provider.go is wrapped with
// registerNoizuExtensions(...). Keeping all Noizu code in noizu_* files
// keeps rebases onto new upstream releases conflict-free. See NOIZU.md.
package github

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// noizuResources holds Noizu-specific managed resources.
//
// Resource names MUST follow the github_noizu_<thing> convention so they
// can never collide with upstream names and are trivially attributable in
// state and configs. Implement new resources in github/noizu_resource_<thing>.go
// and register them here.
var noizuResources = map[string]*schema.Resource{}

// noizuDataSources holds Noizu-specific data sources.
//
// Same convention as noizuResources: names start with github_noizu_,
// implementations live in github/noizu_data_source_<thing>.go.
var noizuDataSources = map[string]*schema.Resource{}

// registerNoizuExtensions merges the Noizu extension maps into the
// provider's resource and data source maps. It panics on any name
// collision with upstream entries so a drift (upstream growing a name we
// also claimed) is caught at provider construction, never at plan time.
// It returns p so NewProvider can wrap its return expression in a single
// statement.
func registerNoizuExtensions(p *schema.Provider) *schema.Provider {
	mergeNoizuMaps("resource", p.ResourcesMap, noizuResources)
	mergeNoizuMaps("data source", p.DataSourcesMap, noizuDataSources)
	return p
}

// mergeNoizuMaps copies src into dst, panicking on collision.
func mergeNoizuMaps(kind string, dst, src map[string]*schema.Resource) {
	for name, r := range src {
		if _, exists := dst[name]; exists {
			panic(fmt.Sprintf("noizu extension layer: %s %q already exists upstream; rename it (names must start with github_noizu_)", kind, name))
		}
		dst[name] = r
	}
}
