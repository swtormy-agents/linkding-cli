# Agent guide: linkding-cli

This document complements the bundled skill (`linkding skill print`). The CLI is **unofficial** and not affiliated with linkding.

## Setup

1. Install `linkding` (release binary or `go install`).
2. User sets `LINKDING_URL` and `LINKDING_TOKEN`, or a config file under the OS user config directory with `chmod 600`.
3. Run `linkding status` — exit **0** even when `"has_token": false`; exit **3** only when a token is set but `GET /api/user/profile/` fails auth.

## Command map

| Area | Commands |
|------|----------|
| Bookmarks | `list`, `list-archived`, `search QUERY`, `get ID`, `check URL`, `create`, `update ID`, `archive ID`, `unarchive ID`, `delete ID` |
| Tags | `list`, `get ID`, `create` |
| Bundles | `list` |
| User | `user profile` |
| Policy | `access show`, `access write`, `access read` |
| Meta | `status`, `skill print`, `skill install` |

Global flags: `--dry-run`, `--allow-write`, `--jq`, `--pretty`, `--url`, `--http-timeout`.

Mutating bookmark/tag commands accept `--json` or `--body` (file; `-` for stdin).

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | API or runtime error |
| 2 | Usage error |
| 3 | Auth failure (missing/invalid token) |
| 4 | Blocked by `allow_write` policy |

## Writes and agents

- Never run `linkding access write` or `linkding access read` for the user.
- Never ask for the token in chat.
- Before a write, show the exact command and get confirmation.
- For automated sandboxes, the human may set `LINKDING_ALLOW_WRITE=true` or `"allow_write": true` in config — not something the agent should set silently.

## Troubleshooting

| Symptom | Action |
|---------|--------|
| stderr `"type":"policy"` | User enables writes (see README) |
| stderr `"type":"auth"` | Check token in linkding Settings |
| `LINKDING_URL is not set` | Export URL or set `url` in config |
| Empty list | Try `bookmarks search` with UI-style query syntax |

API reference: https://linkding.link/api/
