package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"os"
)

// Config holds email configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromName     string
	FromEmail    string
}

// Service handles email operations
type Service struct {
	config *Config
}

// NewService creates a new email service
func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// LoadConfigFromEnv loads email config from environment variables
func LoadConfigFromEnv() *Config {
	return &Config{
		SMTPHost:     getEnv("SMTP_HOST", "mail.karyadeveloperindonesia.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", "no-reply@karyadeveloperindonesia.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "Justformeokay23"),
		FromName:     getEnv("MAIL_FROM_NAME", "Karir Nusantara"),
		FromEmail:    getEnv("MAIL_FROM_EMAIL", "no-reply@karyadeveloperindonesia.com"),
	}
}

// SendEmail sends an email with STARTTLS support and custom TLS config
// Equivalent to PHP's SMTPOptions with verify_peer=false, allow_self_signed=true
func (s *Service) SendEmail(to string, subject string, body string) error {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	
	// Setup message headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to SMTP server
	host := s.config.SMTPHost
	addr := fmt.Sprintf("%s:%s", host, s.config.SMTPPort)
	
	// Establish connection
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// TLS config - equivalent to PHP's SMTPOptions
	// verify_peer => false, verify_peer_name => false, allow_self_signed => true
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Allow self-signed certificates
		ServerName:         host,
	}

	// Start TLS (STARTTLS for port 587)
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	client.Quit()

	return nil
}

// SendWelcomeEmail sends welcome email to new company
func (s *Service) SendWelcomeEmail(to string, companyName string, fullName string) error {
	subject := "Selamat Datang di Karir Nusantara"
	
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9fafb; }
		.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Selamat Datang di Karir Nusantara!</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.FullName}}</strong>,</p>
			<p>Terima kasih telah mendaftar di Karir Nusantara sebagai <strong>{{.CompanyName}}</strong>.</p>
			<p>Akun Anda telah berhasil dibuat. Silakan lengkapi profil perusahaan Anda dan unggah dokumen-dokumen yang diperlukan untuk proses verifikasi.</p>
			<p>Setelah verifikasi disetujui, Anda dapat mulai memposting lowongan pekerjaan dan menemukan talenta terbaik untuk perusahaan Anda.</p>
			<a href="https://company.karirnusantara.com" class="button">Login ke Dashboard</a>
			<p>Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi tim dukungan kami.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		FullName    string
		CompanyName string
	}{
		FullName:    fullName,
		CompanyName: companyName,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}

// SendPasswordResetEmail sends password reset email
func (s *Service) SendPasswordResetEmail(to string, resetToken string, fullName string) error {
	subject := "Reset Password - Karir Nusantara"
	
	// URL untuk reset password di frontend
	resetURL := fmt.Sprintf("https://company.karirnusantara.com/reset-password?token=%s", resetToken)
	
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9fafb; }
		.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.warning { background-color: #fef3c7; padding: 15px; border-left: 4px solid #f59e0b; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Reset Password</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.FullName}}</strong>,</p>
			<p>Kami menerima permintaan untuk mereset password akun Anda di Karir Nusantara.</p>
			<p>Klik tombol di bawah ini untuk mereset password Anda:</p>
			<a href="{{.ResetURL}}" class="button">Reset Password</a>
			<p>Atau salin dan tempel URL berikut ke browser Anda:</p>
			<p style="word-break: break-all; background-color: #e5e7eb; padding: 10px; border-radius: 5px;">{{.ResetURL}}</p>
			<div class="warning">
				<strong>Perhatian:</strong>
				<ul>
					<li>Link ini hanya berlaku selama 1 jam</li>
					<li>Jika Anda tidak meminta reset password, abaikan email ini</li>
					<li>Pastikan untuk tidak membagikan link ini kepada siapapun</li>
				</ul>
			</div>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("reset").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		FullName string
		ResetURL string
	}{
		FullName: fullName,
		ResetURL: resetURL,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
