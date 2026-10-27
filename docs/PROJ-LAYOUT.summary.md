# Project Layout (summary)

```text
terraform-provider-github/          # noizu-forks fork of integrations/terraform-provider-github (base v6.13.0)
├── main.go                         # plugin.Serve entry point
├── github/                         # provider package: provider.go, config.go, transport.go, noizu_provider.go
│   ├── resource_github_*.go        # 106 impl files (88 registered resources)
│   ├── data_source_github_*.go     # 75 registered data sources
│   ├── util*.go                    # 15 utility files
│   └── noizu_*.go                  # Noizu extension layer (scaffold, no resources yet)
├── internal/ghclient/              # new client: Source (token/app/anonymous), REST+GraphQL, cache, throttle
├── docs/                           # generated provider docs + PROJ-*/THREAT-MODEL docs
├── templates/                      # tfplugindocs templates
├── examples/                       # HCL examples (resources 87, data-sources 74)
├── scripts/noizu/                  # build-dev-override.sh (OpenTofu dev_overrides)
├── tools/                          # pinned tfplugindocs module
├── .github/                        # 8 workflows (DISABLED on the fork), instructions, templates
├── GNUmakefile                     # build/lint/test/testacc/sweep/generatedocs
├── ARCHITECTURE.md                 # contributor implementation guide
├── DECISIONS.md                    # decision log
├── RESOURCES.md                    # per-resource standards tracker
├── NOIZU.md                        # fork contract (conflict-avoidance rules, branch model)
├── NOIZU-CHANGELOG.md              # Noizu-only change log
├── VERIFY_ATTESTATIONS.md          # release attestation verification
└── .meta/pointers.yaml             # doc-pointer store
```

Details: [PROJ-LAYOUT.md](PROJ-LAYOUT.md), [layout/github.md](layout/github.md), [layout/internal.md](layout/internal.md), [layout/github-dir.md](layout/github-dir.md).

No relational database, images, or deploy manifests: the artifact is one Go plugin binary.
