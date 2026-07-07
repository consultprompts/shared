# shared/email

Shared email module for consultprompts services. Wraps the [Resend](https://resend.com) SDK and provides typed methods for all transactional emails sent by the platform.

## Installation

```
go get github.com/consultprompts/shared/email
```

For local development with the monorepo, add a `replace` directive to your `go.mod`:

```
replace github.com/consultprompts/shared/email => ../shared/email
```

## Configuration

`RESEND_API_KEY` and `RESEND_FROM` are required. If either is missing, `NewClient()` returns `nil` and callers should treat email as disabled.

| Variable         | Description                          |
|------------------|--------------------------------------|
| `RESEND_API_KEY` | API key from your Resend dashboard   |
| `RESEND_FROM`    | Sender address (e.g. `no-reply@consultprompts.com`) |
| `LOGO_URL`       | Optional. Logo image shown in the email header; falls back to a text wordmark when unset |
| `FRONTEND_URL`   | Used by `SendNewLeadNotification` to link to `/admin-console` |

## Usage

```go
client := email.NewClient()
if client == nil {
    // email not configured — handle gracefully
}

// Auth emails
client.SendVerificationEmail(to, token, frontendURL)
client.SendPasswordResetEmail(to, token, frontendURL)
client.SendLoginNotificationEmail(to, frontendURL)

// Lead emails — recipient comes from LeadData.Email, not a separate param
client.SendNewLeadNotification(email.LeadData{...})
client.SendLeadConfirmation(email.LeadData{...})
client.SendLeadAccepted(email.LeadData{...}, frontendURL)
```

### LeadData fields

```go
type LeadData struct {
    Name      string
    Email     string
    Business  string
    Package   *string   // optional
    Message   *string   // optional
    CreatedAt time.Time
}
```

## Emails

| Method                      | Recipient      | Description                                      |
|-----------------------------|----------------|--------------------------------------------------|
| `SendVerificationEmail`     | New user       | Email verification link (expires 24h)            |
| `SendPasswordResetEmail`    | User           | Password reset link (expires 1h)                 |
| `SendLoginNotificationEmail`| User           | Alert on new login                               |
| `SendNewLeadNotification`   | Lead (`data.Email`) | New mockup request details, links to `/admin-console` |
| `SendLeadConfirmation`      | Lead           | Confirmation to the person who submitted the form|
| `SendLeadAccepted`          | Lead           | Notifies the client their project was accepted, links to `/my-projects` |

Everything past mockup delivery — mockup-ready, revision requests, payment
requests/receipts, and launch — is templated directly in agency-service's own
`internal/email/email.go`, not in this shared module.
