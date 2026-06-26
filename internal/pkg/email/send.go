package email

import (
	"crypto/tls"
	"godan/internal/config"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// Configured 是否已配置可发送邮件的最小项（host、from）。
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
	e.HTML = []byte(body)
	auth := smtp.PlainAuth("", e.From, c.Password, c.Host)
	tlsConfig := &tls.Config{ServerName: c.Host}
	return e.SendWithTLS(c.Host+":"+c.Port, auth, tlsConfig)
}
