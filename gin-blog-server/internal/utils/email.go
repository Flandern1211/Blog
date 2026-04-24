package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	g "gin-blog/internal/global"
	"html/template"
	"log/slog"

	"github.com/k3a/html2text"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	UserName string
	Subject  string
	Code     string // 6位验证码
	URL      string // 验证链接
}

// GetEmailData 获取邮件数据
func GetEmailData(username, info string) *EmailData {
	// TODO: 从配置中获取基础 URL
	return &EmailData{
		UserName: username,
		Subject:  "注册验证",
		URL:      "http://localhost:8080/api/verify?info=" + info,
	}
}

// SendEmail 发送验证邮件（链接形式）
func SendEmail(username string, data *EmailData) error {
	conf := g.GetConfig().Email
	host := conf.Host
	port := conf.Port
	user := conf.SmtpUser
	pass := conf.SmtpPass
	from := conf.From

	slog.Info(fmt.Sprintf("发送验证邮件 to=%s", username))

	// 这里简单起见，直接使用 buildCodeEmailHTML 的逻辑，或者你可以加载 assets/templates/email-verify.tpl
	// 由于 assets/templates/email-verify.tpl 依赖 base.tpl，这里先用简单的 HTML
	htmlBody := fmt.Sprintf(`<h3>你好，%s</h3><p>请点击以下链接激活账户：</p><a href="%s">%s</a>`, data.UserName, data.URL, data.URL)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", username) // 注意：这里 handle_auth.go 传进来的是 username，但实际上应该是 email。
	// 查阅 handle_auth.go:214: err = utils.SendEmail(regreq.Username,EmailData)
	// 确实传的是 username。这可能是原代码的 bug，或者 username 就是 email。
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}

// SendCodeEmail 发送含验证码的邮件
func SendCodeEmail(email string, data *EmailData) error {
	conf := g.GetConfig().Email
	host := conf.Host
	port := conf.Port
	user := conf.SmtpUser
	pass := conf.SmtpPass
	from := conf.From

	slog.Info(fmt.Sprintf("发送验证码邮件 to=%s code=%s", email, data.Code))

	htmlBody := buildCodeEmailHTML(data.UserName, data.Code)

	prem, _ := premailer.NewPremailerFromString(htmlBody, nil)
	htmlInline, _ := prem.Transform()

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", htmlInline)
	m.AddAlternative("text/plain", html2text.HTML2Text(htmlBody))

	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}

func buildCodeEmailHTML(username, code string) string {
	const tpl = `<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="UTF-8"><title>注册验证码</title></head>
<body style="font-family:Arial,sans-serif;background:#f4f4f4;padding:20px;">
  <div style="max-width:480px;margin:0 auto;background:#fff;border-radius:8px;padding:32px;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
    <h2 style="color:#333;">您好，{{.Username}}</h2>
    <p style="color:#555;">您正在注册账号，验证码为：</p>
    <div style="font-size:36px;font-weight:bold;letter-spacing:8px;color:#4caf50;text-align:center;margin:24px 0;">{{.Code}}</div>
    <p style="color:#999;font-size:13px;">验证码10分钟内有效，请勿泄露给他人。</p>
  </div>
</body>
</html>`

	t, _ := template.New("code").Parse(tpl)
	var buf bytes.Buffer
	t.Execute(&buf, map[string]string{"Username": username, "Code": code})
	return buf.String()
}
