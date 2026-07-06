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
// Icons — emoji render in all major email clients (Gmail, iOS, Apple Mail).
// SVGs are stripped by Gmail so we use emoji inside the cyan circle instead.
// ---------------------------------------------------------------------------

const (
	iconCheck  = `&#10003;`  // ✓
	iconKey    = `&#128273;` // 🔑
	iconSend   = `&#128231;` // 📧
	iconRocket = `&#128640;` // 🚀
	iconShield = `&#128274;` // 🔒
)

// ---------------------------------------------------------------------------
// Template engine
// ---------------------------------------------------------------------------

type tpl struct {
	icon     string      // HTML entity / emoji character
	title    string
	body     string      // may contain safe HTML
	rows     [][2]string // optional label-value pairs shown in a table
	ctaHref  string
	ctaLabel string
	footnote string
}

func compose(t tpl) string {
	logoURL := os.Getenv("LOGO_URL")

	var header string
	if logoURL != "" {
		header = fmt.Sprintf(`
		<table cellpadding="0" cellspacing="0"><tr>
		  <td style="vertical-align:middle; padding-right:12px;">
		    <img src="%s" width="32" height="32" alt="" style="display:block; border:0;" />
		  </td>
		  <td style="vertical-align:middle;">
		    <span style="font-family:'Space Grotesk',Georgia,serif; font-size:17px; font-weight:700; letter-spacing:0.08em; text-transform:uppercase; color:#00F0FF;">Consult Prompts</span>
		  </td>
		</tr></table>`, logoURL)
	} else {
		header = `<span style="font-family:'Space Grotesk',Georgia,serif; font-size:17px; font-weight:700; letter-spacing:0.08em; text-transform:uppercase; color:#00F0FF;">Consult Prompts</span>`
	}

	var rowsHTML string
	if len(t.rows) > 0 {
		rowsHTML = `<table width="100%" cellpadding="0" cellspacing="0" style="border-top:1px solid rgba(255,255,255,0.08); margin-bottom:30px;">`
		for _, row := range t.rows {
			rowsHTML += fmt.Sprintf(`
			<tr>
			  <td style="padding:12px 0; color:#A1A1A1; font-size:11px; letter-spacing:0.1em; text-transform:uppercase; width:110px; vertical-align:top; border-bottom:1px solid rgba(255,255,255,0.05);">%s</td>
			  <td style="padding:12px 0; color:#ffffff; font-size:14px; font-weight:300; vertical-align:top; border-bottom:1px solid rgba(255,255,255,0.05);">%s</td>
			</tr>`, row[0], row[1])
		}
		rowsHTML += `</table>`
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <meta name="color-scheme" content="dark">
  <meta name="supported-color-schemes" content="dark light">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;700&family=Inter:wght@300;400;700;900&display=swap" rel="stylesheet">
  <style>
    :root { color-scheme: dark; }
    @media (prefers-color-scheme: dark) {
      body { background:#050505 !important; }
      .email-wrapper { background:#050505 !important; }
      .email-card { background:#0f0f0f !important; }
      .email-footer { background:#0f0f0f !important; }
    }
  </style>
</head>
<body style="margin:0; padding:0; background:#050505; font-family:'Inter',Arial,Helvetica,sans-serif; color:#ffffff;">
  <table class="email-wrapper" width="100%%" cellpadding="0" cellspacing="0" style="background:#050505; padding:48px 24px 80px;">
    <tr><td align="center">
      <table class="email-card" width="560" cellpadding="0" cellspacing="0" style="max-width:560px; width:100%%; background:#0f0f0f; border:1px solid rgba(255,255,255,0.12); border-radius:14px; overflow:hidden;">

        <tr>
          <td style="background:linear-gradient(135deg,#00F0FF22,#7000FF22); padding:28px 40px; border-bottom:1px solid rgba(255,255,255,0.08);">
            %s
          </td>
        </tr>

        <tr><td style="padding:44px 40px 40px;">

          <table cellpadding="0" cellspacing="0" style="margin-bottom:24px;"><tr>
            <td width="56" height="56" style="width:56px; height:56px; background:rgba(0,240,255,0.1); border-radius:28px; text-align:center; vertical-align:middle;">
              <span style="font-size:22px; line-height:56px; display:block;">%s</span>
            </td>
          </tr></table>

          <h2 style="margin:0 0 10px; font-family:'Space Grotesk',Georgia,serif; font-style:italic; font-size:26px; font-weight:700; letter-spacing:-0.02em; color:#ffffff;">%s</h2>

          <p style="margin:0 0 30px; color:#A1A1A1; font-size:14px; font-weight:300; line-height:1.7;">%s</p>

          %s

          <table cellpadding="0" cellspacing="0" style="margin-bottom:28px;"><tr>
            <td style="background:#00F0FF; border-radius:8px;">
              <a href="%s" style="display:inline-block; padding:15px 30px; font-size:12px; font-weight:900; letter-spacing:0.14em; text-transform:uppercase; color:#050505; text-decoration:none; font-family:'Inter',Arial,Helvetica,sans-serif;">%s</a>
            </td>
          </tr></table>

          <p style="margin:0; font-size:12px; color:#555555; line-height:1.6;">%s</p>

        </td></tr>

        <tr>
          <td class="email-footer" style="padding:22px 40px; border-top:1px solid rgba(255,255,255,0.08); text-align:center; background:#0f0f0f;">
            <p style="margin:0; font-size:10px; color:#555555; letter-spacing:0.1em; text-transform:uppercase;">consultprompts.com</p>
          </td>
        </tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`,
		header,
		t.icon,
		t.title,
		t.body,
		rowsHTML,
		t.ctaHref,
		t.ctaLabel,
		t.footnote,
	)
}

// ---------------------------------------------------------------------------
// Auth emails
// ---------------------------------------------------------------------------

func (c *Client) SendVerificationEmail(to, token, frontendURL string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	return c.send(to, "Verify your email — consultprompts.com", compose(tpl{
		icon:     iconCheck,
		title:    "Verify your email",
		body:     `Welcome to consultprompts.com! Click the button below to verify your email address. This link expires in <span style="color:#00F0FF;font-weight:700;">24 hours</span>.`,
		ctaHref:  link,
		ctaLabel: "Verify Email",
		footnote:  "If you didn't create an account, you can safely ignore this email.",
	}))
}

func (c *Client) SendPasswordResetEmail(to, token, frontendURL string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	return c.send(to, "Reset your password — consultprompts.com", compose(tpl{
		icon:     iconKey,
		title:    "Reset your password",
		body:     `We received a request to reset your consultprompts.com password. This link expires in <span style="color:#00F0FF;font-weight:700;">1 hour</span>.`,
		ctaHref:  link,
		ctaLabel: "Reset Password",
		footnote:  "If you didn't request this, you can safely ignore this email.",
	}))
}

func (c *Client) SendLoginNotificationEmail(to, frontendURL string) error {
	resetLink := fmt.Sprintf("%s/reset-password", frontendURL)
	return c.send(to, "New login detected — consultprompts.com", compose(tpl{
		icon:     iconShield,
		title:    "New login detected",
		body:     fmt.Sprintf(`We detected a new login to your consultprompts.com account. If this was you, no action is needed.<br><br>If this wasn't you, <a href="%s" style="color:#00F0FF;font-weight:700;text-decoration:none;">reset your password immediately</a>.`, resetLink),
		ctaHref:  resetLink,
		ctaLabel: "Reset Password",
		footnote:  "This is an automated security notice — replies to this email are not monitored.",
	}))
}

// ---------------------------------------------------------------------------
// Lead emails
// ---------------------------------------------------------------------------

func (c *Client) SendNewLeadNotification(data LeadData) error {
	pkg := "—"
	if data.Package != nil {
		pkg = html.EscapeString(*data.Package)
	}

	rows := [][2]string{
		{"Name", html.EscapeString(data.Name)},
		{"Email", html.EscapeString(data.Email)},
		{"Business", html.EscapeString(data.Business)},
		{"Package", pkg},
	}
	if data.Message != nil {
		rows = append(rows, [2]string{"Message", html.EscapeString(*data.Message)})
	}

	return c.send(
		data.Email,
		fmt.Sprintf("New lead: %s", data.Business),
		compose(tpl{
			icon:     iconSend,
			title:    "New Lead",
			body:     "A new mockup request was submitted on consultprompts.com.",
			rows:     rows,
			ctaHref:  os.Getenv("FRONTEND_URL") + "/admin-console",
			ctaLabel: "View in Dashboard",
			footnote:  fmt.Sprintf("Submitted %s", data.CreatedAt.Format("2006-01-02 15:04 MST")),
		}),
	)
}

func (c *Client) SendLeadConfirmation(data LeadData) error {
	body := fmt.Sprintf(
		`Hey %s — we've got your request for <strong style="color:#ffffff;">%s</strong> and we're already on it. Expect your free mockup within <span style="color:#00F0FF;font-weight:700;">24–48 hours</span>.`,
		html.EscapeString(data.Name),
		html.EscapeString(data.Business),
	)
	if data.Package != nil {
		body += fmt.Sprintf(`<br><br>Package: <span style="color:#00F0FF;">%s</span>`, html.EscapeString(*data.Package))
	}

	return c.send(data.Email, "Transmission received — your mockup is in the queue", compose(tpl{
		icon:     iconRocket,
		title:    "Transmission Received",
		body:     body,
		ctaHref:  "https://wa.me/13026622736",
		ctaLabel: "Chat on WhatsApp",
		footnote:  "Have questions? Just reply to this email — we read every one.",
	}))
}

func (c *Client) SendLeadAccepted(data LeadData, frontendURL string) error {
	body := fmt.Sprintf(
		`Great news, %s — we've accepted your project for <strong style="color:#ffffff;">%s</strong> and we're getting started. Track every milestone in real time from your project dashboard.`,
		html.EscapeString(data.Name),
		html.EscapeString(data.Business),
	)
	if data.Package != nil {
		body += fmt.Sprintf(`<br><br>Package: <span style="color:#00F0FF;">%s</span>`, html.EscapeString(*data.Package))
	}

	return c.send(data.Email, "Your project has been accepted — consultprompts.com", compose(tpl{
		icon:     iconRocket,
		title:    "Project Accepted!",
		body:     body,
		ctaHref:  frontendURL + "/my-projects",
		ctaLabel: "Track My Project",
		footnote:  "Have questions? Just reply to this email — we read every one.",
	}))
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (c *Client) send(to, subject, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	_, err := c.client.Emails.Send(params)
	return err
}
