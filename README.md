# kubo-mcp

Stateless HTTP MCP server (`POST /mcp`) that stores files on IPFS through a [Kubo](https://github.com/ipfs/kubo) node's RPC API. Uploads are streamed, pinned and addressed with CIDv1.

## Tools

| Tool | Does |
| --- | --- |
| `add_file` | Upload base64 content or a `data:` URI |
| `add_file_from_url` | Download a public URL and upload it |
| `get_file` | Read `<cid>[/sub/path]`; images come back as image content |
| `get_file_info` | CID, type, size, gateway URL |
| `pin_file` / `unpin_file` | Pin or unpin a CID |
| `list_pins` | List recursive pins |

## Run

```bash
docker run --rm -p 3000:3000 -e KUBO__API_URL=http://host.docker.internal:5001 ghcr.io/fournyp/kubo-mcp:latest
claude mcp add kubo --transport http http://localhost:3000/mcp --header "Authorization: Bearer <key>"
```

`GET /health` answers `ok` without auth.

### Railway

`.railway/railway.ts` declares the `kubo-mcp` service, built from this repo. It points at an existing Kubo: `http://<kubo>.railway.internal:5001` in the same project, its public URL plus `KUBO__AUTHORIZATION` otherwise.

```bash
npm install && railway link
export AUTH__API_KEYS=$(openssl rand -hex 32) KUBO__API_URL=... KUBO__AUTHORIZATION="Bearer ..."
npm run plan && npm run apply
```

Then add a public domain to `kubo-mcp`.

## Configuration

Read from an optional `appsettings.json`, overridden by `SECTION__FIELD` env vars. See [`appsettings.Example.json`](./src/kubo_mcp_host/appsettings.Example.json) and [`.env.example`](./src/kubo_mcp_host/.env.example).

| Key | Default |
| --- | --- |
| `host.port` | `3000` |
| `auth.api_keys` | none (open); bearer tokens, comma-separated |
| `kubo.api_url` | `http://127.0.0.1:5001` |
| `kubo.authorization` | sent verbatim to Kubo |
| `kubo.request_timeout` / `kubo.upload_timeout` | `60s` / `10m` |
| `ipfs.max_upload_size_bytes` / `ipfs.max_read_size_bytes` | 100 MiB / 5 MiB |
| `ipfs.public_gateway_url` | `https://ipfs.io` |
| `downloader.request_timeout` | `5m` |
| `downloader.allow_private_networks` | `false` |
| `mcp.enabled_tools` | all tools; list to restrict, e.g. `add_file,get_file` |
| `logging.level` | `info` |

## Security

- Set `auth.api_keys`: clients send `Authorization: Bearer <key>`. Empty keeps `/mcp` open, for local use only.
- Never expose Kubo's RPC port (`:5001`).
- `add_file_from_url` refuses non-public IPs unless `downloader.allow_private_networks` is set.

## Development

`make build`, `make test`, `make test-integration` (Docker), `make lint`, `make fmt`, `make mocks`, `make docker`. Conventions in [`CLAUDE.md`](./CLAUDE.md).
