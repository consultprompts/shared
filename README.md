# shared

Shared Go modules for the consultprompts platform. Each subdirectory is an independent Go module consumed by the backend services via the `github.com/consultprompts/shared/*` module path.

## Modules

| Module | Path | Description |
|--------|------|-------------|
| [email](./email) | `github.com/consultprompts/shared/email` | Transactional email via Resend |

## Assets

`assets/logo.png` — the logo referenced by `LOGO_URL` in email templates
(`shared/email`'s `compose()` header). Not a Go module; just a shared static file.

## Using a module

Add the dependency to your service's `go.mod`:

```
require github.com/consultprompts/shared/email v0.1.0
```

For local development with the monorepo, add a `replace` directive:

```
replace github.com/consultprompts/shared/email => ../shared/email
```

## Publishing a new version

Modules are versioned with path-prefixed tags. To release `email` at `v0.1.0`:

```sh
git tag email/v0.1.0
git push origin email/v0.1.0
```
