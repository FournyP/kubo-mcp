# CLAUDE.md

Stateless HTTP MCP server storing files on IPFS through Kubo's RPC API. One Go module, one binary (`src/kubo_mcp_host`). Mirror the closest sibling package when unsure.

## Layout (hexagonal)

| Module (`src/`) | Role |
| --- | --- |
| `kubo_mcp` | domain: use cases (+ requests/responses), ports, models, settings |
| `kubo_mcp_http_kubo` | Kubo RPC adapter, generated client |
| `kubo_mcp_http_downloader` | URL download adapter + SSRF guard |
| `kubo_mcp_mcp` | MCP tools, `/mcp` controller, bearer auth middleware |
| `kubo_mcp_host` | composition root, routes, `main.go`, Dockerfile |
| `kubo_mcp_common` | errors, koanf, logging |

- Business logic lives in `kubo_mcp` only; adapters translate.
- A tool decodes args, calls one use case, maps the result.
- Use cases validate their request first and depend on ports only.
- New behavior: port → adapter → use case → tool, each wired in its `configurations/`.

## Style

- Comments: one line, only the non-obvious why. Contracts go on the port.
- `snake_case.go`, one type per file, `NewX`, ports `XInterface`, `Add<Module>Configuration`.
- Fields named after their type: `ipfsStorage`, `kuboSettings`.
- Private helpers are methods on the type that uses them.
- Every setting works from `appsettings.json` and `SECTION__FIELD` env. Bad required value → `logger.Fatal`.
- Wrap errors with `%w`; `ValidationError` for bad input, `FileTooLargeError` for size caps.

## Tests

`tests/*_unit_tests`, mirroring `src`:
- `when_<x>_test.go`, `TestWhen<X>`, `When<X>TestingSuite` (`sut` + mocks), `When<X>BeforeEach`.
- Nested `Given…`/`Should…`, `t.Parallel()` except with `t.Setenv`, no tables.
- Expectations in a per-Given `initSuite`, exact args, `.Times(n)`.
- Kubo adapter: gomock on the generated client, plus integration tests against Kubo. Downloader: httptest.
- CIDs from `common_tests.RandomCid()`. Mocks are generated: `make mocks`.

`tests/*_integration_tests` (`make test-integration`, needs Docker) hit a real Kubo shared through `common_tests/container_builders`. Use `common_tests.RandomContent` so parallel tests never share a CID. Not run in CI.

## Generated Kubo client

`clients/oapi_client.gen.go` comes from the hand-written `clients/kubo_rpc.v3.yaml` (oapi-codegen v2.5.1). Never edit it. New command: update the YAML, run `clients/generate_client.sh`, `make mocks`, then smoke-test against `ipfs/kubo`.

## Railway

`.railway/railway.ts` (TypeScript, `npm run plan|apply`) declares only `kubo-mcp`. Keep its env in sync when adding a setting.

## Commands

`make fmt`, `make lint`, `make test` before handing back. `make mocks` after changing an interface.
