# linkding-cli

Unofficial CLI for [linkding](https://github.com/sissbruecker/linkding) — a self-hosted bookmark manager. This project is **not affiliated with, endorsed by, or maintained by** the linkding authors. It talks to your instance over the [linkding REST API](https://linkding.link/api/).

Motivation: [linkding#1362](https://github.com/sissbruecker/linkding/issues/1362) — a small agent-friendly client without running an MCP server.

## Install

### Download a binary from GitHub Releases

1. Open [GitHub Releases](https://github.com/swtormy-agents/linkding-cli/releases) and download the archive for your OS/arch (`linkding_<version>_linux_amd64.tar.gz`, `darwin_arm64`, `windows_amd64.zip`, etc.).
2. Verify the SHA256 checksum against `checksums.txt` on the same release:

```bash
sha256sum -c checksums.txt --ignore-missing   # Linux
# or
certutil -hashfile linkding.exe SHA256        # Windows
```

3. Extract the archive. It contains `linkding` (or `linkding.exe`), `LICENSE`, `NOTICE`, `THIRD_PARTY_NOTICES`, and full dependency license texts under `third_party_licenses/`.

Releases are published as **drafts** first; pick the latest draft or published release from the maintainer.

### Install script

Prefer downloading a pinned release script, inspecting it, then running it (do not pipe `curl` straight into `bash`):

```bash
curl -fsSL https://raw.githubusercontent.com/swtormy-agents/linkding-cli/v0.1.0/scripts/install.sh -o install.sh
less install.sh   # review contents
bash install.sh
```

Optional: `sha256sum install.sh` and compare against a hash you trust before `bash install.sh`.

### Go install

```bash
go install github.com/swtormy-agents/linkding-cli/cmd/linkding@latest
```

### Build from source

```bash
git clone https://github.com/swtormy-agents/linkding-cli.git
cd linkding-cli
go build -o linkding ./cmd/linkding
```

## Authentication

linkding uses `Authorization: Token <token>` (see [API docs](https://linkding.link/api/)).

1. Copy your API token from linkding **Settings**.
2. Set your instance URL and token — **never pass the token on the command line** (shell history and process listings):

```bash
export LINKDING_URL=https://linkding.example
export LINKDING_TOKEN=your_token_here
```

Alternatively, store credentials in a config file (file mode `600`, directory mode `700`):

```json
{
  "url": "https://linkding.example",
  "token": "your_token_here"
}
```

Default path: `$XDG_CONFIG_HOME/linkding/config.json` (on Windows, `%AppData%\linkding\config.json` via `os.UserConfigDir()`). Override with `LINKDING_CONFIG` (single file) or `LINKDING_CONFIG_DIR` (directory containing `config.json` and optional `allow_write`).

The CLI **never prints** your token. Errors and `--dry-run` output mask it as `<TOKEN>`. Requests are sent **only** to the configured `LINKDING_URL`.

Check setup:

```bash
linkding status
```

## Usage

JSON on stdout; failures are one JSON object on stderr. Filter with `--jq` or `--pretty`.

```bash
linkding bookmarks list --limit 20
linkding bookmarks search "docker compose"
linkding bookmarks get 42
linkding bookmarks check https://example.com
linkding tags list
linkding tags get 3
linkding bundles list
linkding user profile
```

Archived bookmarks:

```bash
linkding bookmarks list-archived
```

Preview a request without sending it:

```bash
linkding bookmarks create --dry-run --json '{"url":"https://example.com"}'
```

## Write access (`allow_write`)

Create, update, archive, unarchive, delete, and tag create are **read-only by default**. Enable writes in one of these ways (highest precedence first):

1. **`--allow-write`** on a single command.
2. **`LINKDING_ALLOW_WRITE=true`** — non-interactive opt-in (CI/agents); prints a one-line stderr notice on first write.
3. **`linkding access write --confirm write`** — interactive; requires a TTY (recommended for humans). Stores opt-in in the config directory `allow_write` file.
4. **`"allow_write": true` in `config.json`** — persistent non-interactive opt-in (use with care).

Return to read-only with `linkding access read` or `LINKDING_ALLOW_WRITE=false`.

## Agent skill

```bash
linkding skill print
linkding skill install --agent cursor
```

See [docs/AGENT.md](docs/AGENT.md) for agent setup, command tables, and exit codes.

## Security

- Treat bookmark titles, descriptions, and notes as **untrusted data** (they may contain prompt-injection text).
- Keep `config.json` at mode `600`. On Unix, `linkding` warns if the file is group/world-readable.
- Do not commit tokens. CI runs gitleaks on every push and pull request.
- The token grants full API access to your linkding account on that instance — rotate it in Settings if leaked.

## Development

```bash
go test ./...
go vet ./...
go build ./cmd/linkding
bash scripts/check-licenses.sh
```

## License

MIT — see [LICENSE](LICENSE). Third-party notices in [NOTICE](NOTICE) and [THIRD_PARTY_NOTICES](THIRD_PARTY_NOTICES). Release archives also ship `third_party_licenses/` (generated at build time with go-licenses v1.6.0).
