package email

import (
	"fmt"
	"html"
	"os"
	"time"

	"github.com/resend/resend-go/v2"
)

// Client wraps the Resend SDK. NewClient returns nil when RESEND_API_KEY or
// RESEND_FROM are not set, so callers can treat email as optional.
type Client struct {
	client *resend.Client
	from   string
}

func NewClient() *Client {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("RESEND_FROM")
	if apiKey == "" || from == "" {
		return nil
	}
	return &Client{client: resend.NewClient(apiKey), from: from}
}

// LeadData carries the fields needed for lead-related emails.
type LeadData struct {
	Name      string
	Email     string
	Business  string
	Package   *string
	Message   *string
	CreatedAt time.Time
}

// ---------------------------------------------------------------------------
// Auth emails
// ---------------------------------------------------------------------------

func (c *Client) SendVerificationEmail(to, token, frontendURL string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">Verify your email</h2>
	<p style="margin:0 0 32px;color:#A1A1A1;font-size:14px;line-height:1.6;">
		Welcome to consultprompts.com! Click the button below to verify your email address.<br>
		This link expires in <span style="color:#00F0FF;font-weight:700;">24 hours</span>.
	</p>
	<table cellpadding="0" cellspacing="0" style="margin:0 0 32px;">
		<tr>
			<td style="background:#00F0FF;border-radius:6px;">
				<a href="%s" style="display:inline-block;padding:14px 28px;font-size:13px;font-weight:900;letter-spacing:0.12em;text-transform:uppercase;color:#050505;text-decoration:none;">
					Verify Email
				</a>
			</td>
		</tr>
	</table>
	<p style="margin:0;font-size:12px;color:#555555;">If you didn't create an account, you can safely ignore this email.</p>`,
		link,
	)
	return c.send(to, "Verify your email — consultprompts.com", body)
}

func (c *Client) SendPasswordResetEmail(to, token, frontendURL string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">Reset your password</h2>
	<p style="margin:0 0 32px;color:#A1A1A1;font-size:14px;line-height:1.6;">
		We received a request to reset your consultprompts.com password.<br>
		This link expires in <span style="color:#00F0FF;font-weight:700;">1 hour</span>.
	</p>
	<table cellpadding="0" cellspacing="0" style="margin:0 0 32px;">
		<tr>
			<td style="background:#00F0FF;border-radius:6px;">
				<a href="%s" style="display:inline-block;padding:14px 28px;font-size:13px;font-weight:900;letter-spacing:0.12em;text-transform:uppercase;color:#050505;text-decoration:none;">
					Reset Password
				</a>
			</td>
		</tr>
	</table>
	<p style="margin:0;font-size:12px;color:#555555;">If you didn't request this, you can safely ignore this email.</p>`,
		link,
	)
	return c.send(to, "Reset your password — consultprompts.com", body)
}

func (c *Client) SendLoginNotificationEmail(to, frontendURL string) error {
	resetLink := fmt.Sprintf("%s/reset-password", frontendURL)
	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">New login detected</h2>
	<p style="margin:0 0 24px;color:#A1A1A1;font-size:14px;line-height:1.6;">
		We detected a new login to your consultprompts.com account.<br>
		If this was you, no action is needed.
	</p>
	<p style="margin:0;font-size:14px;color:#A1A1A1;line-height:1.6;">
		If this wasn't you, <a href="%s" style="color:#00F0FF;text-decoration:none;font-weight:700;">reset your password immediately</a>.
	</p>`,
		resetLink,
	)
	return c.send(to, "New login detected — consultprompts.com", body)
}

// ---------------------------------------------------------------------------
// Lead emails
// ---------------------------------------------------------------------------

func (c *Client) SendLeadAccepted(data LeadData, frontendURL string) error {
	link := frontendURL + "/my-projects"
	pkgRow := ""
	if data.Package != nil {
		pkgRow = fmt.Sprintf(
			`<p style="margin:0 0 8px;font-size:13px;color:#A1A1A1;">Package: <span style="color:#00F0FF;">%s</span></p>`,
			html.EscapeString(*data.Package),
		)
	}

	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">Project Accepted!</h2>
	<p style="margin:0 0 24px;color:#A1A1A1;font-size:14px;line-height:1.6;">
		Great news, %s — we've accepted your project for <strong style="color:#ffffff;">%s</strong> and we're getting started.<br>
		Track every milestone in real time from your project dashboard.
	</p>
	%s
	<table cellpadding="0" cellspacing="0" style="margin:32px 0;">
		<tr>
			<td style="background:#00F0FF;border-radius:6px;">
				<a href="%s" style="display:inline-block;padding:14px 28px;font-size:13px;font-weight:900;letter-spacing:0.12em;text-transform:uppercase;color:#050505;text-decoration:none;">
					Track My Project
				</a>
			</td>
		</tr>
	</table>
	<p style="margin:0;font-size:12px;color:#555555;line-height:1.6;">Have questions? Just reply to this email — we read every one.</p>`,
		html.EscapeString(data.Name),
		html.EscapeString(data.Business),
		pkgRow,
		link,
	)
	return c.send(data.Email, "Your project has been accepted — consultprompts.com", body)
}

func (c *Client) SendNewLeadNotification(data LeadData) error {
	to := data.Email
	pkg := "—"
	if data.Package != nil {
		pkg = html.EscapeString(*data.Package)
	}
	message := "—"
	if data.Message != nil {
		message = html.EscapeString(*data.Message)
	}

	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">New Lead</h2>
	<p style="margin:0 0 32px;color:#A1A1A1;font-size:14px;">A new mockup request was submitted on consultprompts.com</p>
	<table width="100%%" cellpadding="0" cellspacing="0" style="border-top:1px solid rgba(255,255,255,0.08);">
		%s%s%s%s%s
	</table>
	<p style="margin:32px 0 0;font-size:12px;color:#555555;">Submitted %s</p>`,
		labelRow("Name", html.EscapeString(data.Name)),
		labelRow("Email", html.EscapeString(data.Email)),
		labelRow("Business", html.EscapeString(data.Business)),
		labelRow("Package", pkg),
		labelRow("Message", message),
		data.CreatedAt.Format("2006-01-02 15:04 MST"),
	)
	return c.send(to, fmt.Sprintf("New lead: %s", data.Business), body)
}

func (c *Client) SendLeadConfirmation(data LeadData) error {
	pkgRow := ""
	if data.Package != nil {
		pkgRow = fmt.Sprintf(
			`<p style="margin:0 0 8px;font-size:13px;color:#A1A1A1;">Package: <span style="color:#00F0FF;">%s</span></p>`,
			html.EscapeString(*data.Package),
		)
	}

	body := fmt.Sprintf(`
	<h2 style="margin:0 0 8px;font-size:24px;font-weight:900;letter-spacing:-0.02em;color:#ffffff;">Transmission Received</h2>
	<p style="margin:0 0 24px;color:#A1A1A1;font-size:14px;line-height:1.6;">
		Hey %s — we've got your request for <strong style="color:#ffffff;">%s</strong> and we're already on it.<br>
		Expect your free mockup within <span style="color:#00F0FF;font-weight:700;">24–48 hours</span>.
	</p>
	%s
	<table cellpadding="0" cellspacing="0" style="margin:32px 0;">
		<tr>
			<td style="background:#00F0FF;border-radius:6px;">
				<a href="https://wa.me/13026622736" style="display:inline-block;padding:14px 28px;font-size:13px;font-weight:900;letter-spacing:0.12em;text-transform:uppercase;color:#050505;text-decoration:none;">
					Chat on WhatsApp
				</a>
			</td>
		</tr>
	</table>
	<p style="margin:0;font-size:12px;color:#555555;line-height:1.6;">Have questions? Just reply to this email — we read every one.</p>`,
		html.EscapeString(data.Name),
		html.EscapeString(data.Business),
		pkgRow,
	)
	return c.send(data.Email, "Transmission received — your mockup is in the queue", body)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (c *Client) send(to, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		Html:    wrap(body),
	}
	_, err := c.client.Emails.Send(params)
	return err
}

func wrap(body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#050505;font-family:Arial,Helvetica,sans-serif;color:#ffffff;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#050505;padding:48px 16px;">
    <tr><td align="center">
      <table width="560" cellpadding="0" cellspacing="0" style="max-width:560px;width:100%%;background:#0f0f0f;border:1px solid rgba(255,255,255,0.12);border-radius:12px;overflow:hidden;">
        <tr>
          <td style="background:linear-gradient(135deg,#00F0FF22,#7000FF22);padding:32px 40px;border-bottom:1px solid rgba(255,255,255,0.08);">
            <span style="font-size:18px;font-weight:900;letter-spacing:0.15em;text-transform:uppercase;color:#00F0FF;">CONSULTPROMPTS</span>
          </td>
        </tr>
        <tr><td style="padding:40px;">%s</td></tr>
        <tr>
          <td style="padding:24px 40px;border-top:1px solid rgba(255,255,255,0.08);text-align:center;">
            <p style="margin:0;font-size:11px;color:#555555;letter-spacing:0.08em;text-transform:uppercase;">
              consultprompts.com &nbsp;·&nbsp; Helping local businesses look world-class
            </p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, body)
}

func labelRow(label, value string) string {
	return fmt.Sprintf(`
	<tr>
		<td style="padding:10px 0;color:#A1A1A1;font-size:12px;letter-spacing:0.1em;text-transform:uppercase;width:120px;vertical-align:top;">%s</td>
		<td style="padding:10px 0;color:#ffffff;font-size:14px;vertical-align:top;">%s</td>
	</tr>`, label, value)
}
