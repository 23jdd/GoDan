package email

import (
	"crypto/tls"
	"net/smtp"

	"github.com/jordan-wright/email"

	"godan/internal/config"
)

type Client struct {
	From     string
	Password string
	Host     string
	Port     string
}

func NewEmail(c *config.Config) *Client {
	return &Client{
		From:     c.Stmp.From,
		Password: c.Stmp.Password,
		Host:     c.Stmp.Host,
		Port:     c.Stmp.Port,
	}
}

// SendPlain 发送纯文本邮件（UTF-8）。适用于验证码、通知等。
// to：收件人列表；subject、body：主题与正文。
func (c *Client) SendPlain(to []string, subject, body string) error {
	e := email.NewEmail()
	e.From = c.From
	e.To = to
	e.Subject = subject
	e.Text = []byte(body)
	auth := smtp.PlainAuth("", e.From, c.Password, c.Host)
	tlsConfig := &tls.Config{ServerName: c.Host}
	return e.SendWithTLS(c.Host+":"+c.Port, auth, tlsConfig)
}

// SendHtml 发送 HTML 邮件。
func (c *Client) SendHtml(to []string, subject, htmlBody string) error {
	e := email.NewEmail()
	e.From = c.From
	e.To = to
	e.Subject = subject
	e.HTML = []byte(htmlBody)
	auth := smtp.PlainAuth("", e.From, c.Password, c.Host)
	tlsConfig := &tls.Config{ServerName: c.Host}
	return e.SendWithTLS(c.Host+":"+c.Port, auth, tlsConfig)
}

// SendCode 发送验证码 HTML 邮件。
func (c *Client) SendCode(to, verify string) error {
	return c.SendVerificationCode(to, verify)
}

// SendVerificationCode 发送验证码 HTML 邮件。
func (c *Client) SendVerificationCode(to, code string) error {
	html := CreateMessageHtml(code)
	e := email.NewEmail()
	e.From = c.From
	e.To = []string{to}
	e.Subject = "GoDan 验证码"
	e.HTML = []byte(html)
	auth := smtp.PlainAuth("", e.From, c.Password, c.Host)
	tlsConfig := &tls.Config{ServerName: c.Host}
	return e.SendWithTLS(c.Host+":"+c.Port, auth, tlsConfig)
}

func CreateMessageHtml(verify string) string {
	return `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { margin:0; padding:0; background:#f4f4f4; font-family:Arial,sans-serif; }
  .container { max-width:600px; margin:40px auto; background:#fff; border-radius:8px; overflow:hidden; box-shadow:0 2px 12px rgba(0,0,0,.1); }
  .header { background:#00a1d6; padding:30px; text-align:center; }
  .header h1 { color:#fff; margin:0; font-size:24px; }
  .body { padding:40px 30px; }
  .body p { color:#666; font-size:14px; line-height:1.8; margin:0 0 16px; }
  .code-box { background:#f0f9ff; border:2px dashed #00a1d6; border-radius:6px; padding:20px; text-align:center; margin:24px 0; }
  .code { font-size:36px; font-weight:bold; color:#00a1d6; letter-spacing:6px; }
  .expire { color:#999; font-size:12px; text-align:center; }
  .footer { background:#fafafa; padding:20px; text-align:center; border-top:1px solid #eee; }
  .footer p { color:#bbb; font-size:12px; margin:0; }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <h1>GoDan</h1>
  </div>
  <div class="body">
    <p>您好，</p>
    <p>您的验证码如下，请在 5 分钟内完成验证：</p>
    <div class="code-box">
      <span class="code">` + verify + `</span>
    </div>
    <p class="expire">验证码 5 分钟内有效，请勿泄露给他人。</p>
    <p>如非本人操作，请忽略此邮件。</p>
  </div>
  <div class="footer">
    <p>GoDan - 视频分享平台</p>
  </div>
</div>
</body>
</html>`
}
