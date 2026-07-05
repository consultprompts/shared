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

Two environment variables are required. If either is missing, `NewClient()` returns `nil` and callers should treat email as disabled.

| Variable         | Description                          |
|------------------|--------------------------------------|
| `RESEND_API_KEY` | API key from your Resend dashboard   |
| `RESEND_FROM`    | Sender address (e.g. `no-reply@consultprompts.com`) |

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

// Lead emails
client.SendNewLeadNotification(notifyTo, email.LeadData{...})
client.SendLeadConfirmation(email.LeadData{...})
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
| `SendNewLeadNotification`   | Internal (ops) | New mockup request details                       |
| `SendLeadConfirmation`      | Lead           | Confirmation to the person who submitted the form|
