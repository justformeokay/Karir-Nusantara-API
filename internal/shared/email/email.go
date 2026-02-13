package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
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
	headers["Reply-To"] = s.config.FromEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""
	headers["X-Mailer"] = "Karir Nusantara Mailer"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to SMTP server
	host := s.config.SMTPHost
	addr := fmt.Sprintf("%s:%s", host, s.config.SMTPPort)
	log.Printf("[EMAIL] Connecting to SMTP server: %s", addr)

	// Establish connection
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to connect to SMTP server %s: %v", addr, err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()
	log.Printf("[EMAIL] Successfully connected to SMTP server")

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
	log.Printf("[EMAIL] Sending email with subject: %s", subject)
	w, err := client.Data()
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to get data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to write message: %v", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to close data writer: %v", err)
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	log.Printf("[EMAIL] Email sent successfully to: %s", to)
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

// SendJobSeekerWelcomeEmail sends welcome email to new job seeker
func (s *Service) SendJobSeekerWelcomeEmail(to string, fullName string) error {
	subject := "Selamat Datang di Karir Nusantara!"

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%); color: white; padding: 30px 20px; text-align: center; border-radius: 10px 10px 0 0; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 30px 20px; background-color: #f9fafb; }
		.welcome-box { background-color: #dbeafe; padding: 20px; border-radius: 10px; text-align: center; margin: 20px 0; }
		.welcome-icon { font-size: 48px; margin-bottom: 10px; }
		.feature-list { background-color: #fff; padding: 20px; border-radius: 10px; margin: 20px 0; }
		.feature-item { display: flex; align-items: center; padding: 12px 0; border-bottom: 1px solid #e5e7eb; }
		.feature-item:last-child { border-bottom: none; }
		.feature-icon { font-size: 24px; margin-right: 15px; }
		.button { display: inline-block; padding: 14px 28px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 8px; margin: 20px 0; font-weight: bold; }
		.button:hover { background-color: #1d4ed8; }
		.tips { background-color: #fef3c7; padding: 15px; border-left: 4px solid #f59e0b; margin: 20px 0; border-radius: 0 8px 8px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f3f4f6; border-radius: 0 0 10px 10px; }
		.social-links { margin: 15px 0; }
		.social-links a { margin: 0 10px; color: #2563eb; text-decoration: none; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Selamat Datang di Karir Nusantara!</h1>
		</div>
		<div class="content">
			<div class="welcome-box">
				<div class="welcome-icon">👋</div>
				<h2>Halo, {{.FullName}}!</h2>
				<p>Akun Anda telah berhasil dibuat. Selamat bergabung di platform pencari kerja terpercaya di Indonesia!</p>
			</div>
			
			<h3>🚀 Yang Bisa Anda Lakukan:</h3>
			<div class="feature-list">
				<div class="feature-item">
					<span class="feature-icon">📄</span>
					<div>
						<strong>Buat CV Profesional</strong>
						<p style="margin: 5px 0 0 0; font-size: 14px; color: #666;">Buat CV yang menarik dan profesional dengan mudah</p>
					</div>
				</div>
				<div class="feature-item">
					<span class="feature-icon">🔍</span>
					<div>
						<strong>Cari Lowongan Kerja</strong>
						<p style="margin: 5px 0 0 0; font-size: 14px; color: #666;">Temukan ribuan lowongan dari perusahaan terkemuka</p>
					</div>
				</div>
				<div class="feature-item">
					<span class="feature-icon">📨</span>
					<div>
						<strong>Lamar Pekerjaan</strong>
						<p style="margin: 5px 0 0 0; font-size: 14px; color: #666;">Kirim lamaran dengan sekali klik dan pantau statusnya</p>
					</div>
				</div>
				<div class="feature-item">
					<span class="feature-icon">🎯</span>
					<div>
						<strong>Rekomendasi Personal</strong>
						<p style="margin: 5px 0 0 0; font-size: 14px; color: #666;">Dapatkan rekomendasi pekerjaan sesuai profil Anda</p>
					</div>
				</div>
			</div>
			
			<div class="tips">
				<strong>💡 Tips:</strong> Lengkapi profil dan CV Anda untuk meningkatkan peluang dilihat oleh recruiter!
			</div>
			
			<div style="text-align: center;">
				<a href="https://karirnusantara.com" class="button">Mulai Cari Kerja Sekarang →</a>
			</div>
			
			<p style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb;">
				Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi tim support kami di <a href="mailto:support@karirnusantara.com">support@karirnusantara.com</a>
			</p>
		</div>
		<div class="footer">
			<p>Ikuti kami di social media:</p>
			<div class="social-links">
				<a href="#">Instagram</a> | <a href="#">LinkedIn</a> | <a href="#">Twitter</a>
			</div>
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis karena Anda mendaftar di Karir Nusantara.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("jobseekerwelcome").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		FullName string
	}{
		FullName: fullName,
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

	// URL untuk reset password di frontend job seeker
	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", resetToken)

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

// SendPasswordChangeConfirmationEmail sends confirmation email after password change
func (s *Service) SendPasswordChangeConfirmationEmail(to string, fullName string) error {
	subject := "Password Berhasil Diubah - Karir Nusantara"

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #10b981; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9fafb; }
		.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.alert { background-color: #fef3c7; padding: 15px; border-left: 4px solid #f59e0b; margin: 20px 0; }
		.success { background-color: #d1fae5; padding: 15px; border-left: 4px solid #10b981; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>✓ Password Berhasil Diubah</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.FullName}}</strong>,</p>
			<div class="success">
				<strong>Password akun Anda telah berhasil diubah.</strong>
			</div>
			<p>Perubahan ini dilakukan pada: <strong>{{.ChangeTime}}</strong></p>
			<p>Untuk keamanan akun Anda:</p>
			<ul>
				<li>Semua sesi login aktif telah diakhiri</li>
				<li>Anda perlu login kembali dengan password baru</li>
				<li>Pastikan password Anda tersimpan dengan aman</li>
			</ul>
			<div class="alert">
				<strong>Perhatian:</strong><br>
				Jika Anda tidak melakukan perubahan password ini, segera hubungi tim dukungan kami dan reset password Anda.
			</div>
			<a href="https://company.karirnusantara.com/login" class="button">Login Sekarang</a>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("passwordchange").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		FullName   string
		ChangeTime string
	}{
		FullName:   fullName,
		ChangeTime: "Baru saja",
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}

// SendEmailWithAttachment sends an email with PDF attachment
func (s *Service) SendEmailWithAttachment(to string, subject string, htmlBody string, attachmentPath string) error {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)

	// Read attachment file
	fileData, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %w", err)
	}

	// Get filename
	filename := filepath.Base(attachmentPath)

	// Generate boundary
	boundary := "boundary_karir_nusantara_" + fmt.Sprintf("%d", len(fileData))

	// Build email with attachment
	var message bytes.Buffer

	// Headers
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	message.WriteString("\r\n")

	// HTML Body
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	message.WriteString("\r\n")
	message.WriteString(htmlBody)
	message.WriteString("\r\n\r\n")

	// PDF Attachment
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: application/pdf\r\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("\r\n")

	// Encode file to base64
	encoded := base64.StdEncoding.EncodeToString(fileData)
	// Split into 76-character lines (RFC 2045)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		message.WriteString(encoded[i:end] + "\r\n")
	}

	message.WriteString("\r\n")
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Send via SMTP
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

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Start TLS
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

	_, err = io.Copy(w, &message)
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

// SendPaymentConfirmationEmail sends payment confirmation email with invoice PDF
func (s *Service) SendPaymentConfirmationEmail(to string, companyName string, invoiceNumber string, amount int64, invoicePDFPath string) error {
	subject := "Konfirmasi Pembayaran & Invoice - Karir Nusantara"

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #10b981; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9fafb; }
		.invoice-box { background-color: #fff; padding: 20px; border: 2px solid #10b981; border-radius: 8px; margin: 20px 0; }
		.amount { font-size: 32px; font-weight: bold; color: #10b981; text-align: center; margin: 15px 0; }
		.info-row { display: flex; justify-content: space-between; margin: 10px 0; padding: 10px; background-color: #f3f4f6; border-radius: 4px; }
		.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
		.success-badge { background-color: #d1fae5; color: #065f46; padding: 8px 16px; border-radius: 20px; display: inline-block; font-weight: bold; }
		.attachment-note { background-color: #dbeafe; padding: 15px; border-left: 4px solid: #2563eb; margin: 15px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>✓ Pembayaran Berhasil Dikonfirmasi</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.CompanyName}}</strong>,</p>
			
			<div style="text-align: center; margin: 20px 0;">
				<span class="success-badge">PEMBAYARAN LUNAS</span>
			</div>
			
			<p>Kami dengan senang hati menginformasikan bahwa pembayaran Anda telah <strong>berhasil dikonfirmasi</strong> oleh tim kami.</p>
			
			<div class="invoice-box">
				<h3 style="margin-top: 0; color: #10b981;">Detail Pembayaran</h3>
				
				<div class="info-row">
					<span><strong>No. Invoice:</strong></span>
					<span>{{.InvoiceNumber}}</span>
				</div>
				
				<div class="amount">
					{{.Amount}}
				</div>
				
				<div style="text-align: center; color: #666; font-size: 14px;">
					Status: <strong style="color: #10b981;">LUNAS</strong>
				</div>
			</div>
			
			<h3>Apa Selanjutnya?</h3>
			<ul>
				<li><strong>Kuota job posting</strong> Anda telah ditambahkan dan siap digunakan</li>
				<li>Anda dapat langsung memposting lowongan kerja di dashboard</li>
				<li>Invoice PDF terlampir pada email ini untuk arsip keuangan Anda</li>
			</ul>
			
			<div class="attachment-note">
				<strong>📎 Lampiran:</strong> Invoice pembayaran dalam format PDF sudah terlampir pada email ini. Silakan simpan untuk keperluan pelaporan keuangan perusahaan Anda.
			</div>
			
			<div style="text-align: center;">
				<a href="https://company.karirnusantara.com/dashboard" class="button">Buka Dashboard</a>
			</div>
			
			<p style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb;">
				Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi tim support kami.
			</p>
			
			<p><strong>Terima kasih telah menggunakan Karir Nusantara!</strong></p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("paymentconfirmation").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Format amount to Rupiah
	amountStr := formatRupiahHTML(amount)

	data := struct {
		CompanyName   string
		InvoiceNumber string
		Amount        string
	}{
		CompanyName:   companyName,
		InvoiceNumber: invoiceNumber,
		Amount:        amountStr,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Send email with PDF attachment
	return s.SendEmailWithAttachment(to, subject, body.String(), invoicePDFPath)
}

// formatRupiahHTML formats an amount to Indonesian Rupiah format for HTML
func formatRupiahHTML(amount int64) string {
	amountStr := fmt.Sprintf("%d", amount)

	// Add thousands separator
	result := ""
	for i, digit := range amountStr {
		if i > 0 && (len(amountStr)-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	return "Rp " + result
}

// SendCompanyVerificationEmail sends verification status email to company
func (s *Service) SendCompanyVerificationEmail(to string, companyName string, fullName string, isApproved bool, reason string) error {
	var subject string
	var statusText string
	var statusColor string
	var message string
	var nextSteps string

	if isApproved {
		subject = "Selamat! Akun Perusahaan Anda Telah Diverifikasi - Karir Nusantara"
		statusText = "DISETUJUI"
		statusColor = "#10b981"
		message = "Kami dengan senang hati memberitahukan bahwa akun perusahaan Anda telah berhasil diverifikasi oleh tim kami."
		nextSteps = `
			<ul style="margin: 10px 0; padding-left: 20px;">
				<li>Anda sekarang dapat memposting lowongan pekerjaan</li>
				<li>Akses fitur pencarian kandidat</li>
				<li>Kelola lamaran yang masuk</li>
				<li>Gunakan fitur chat untuk berkomunikasi dengan kandidat</li>
			</ul>
		`
	} else {
		subject = "Informasi Status Verifikasi Akun - Karir Nusantara"
		statusText = "DITOLAK"
		statusColor = "#ef4444"
		message = "Mohon maaf, setelah tim kami meninjau dokumen dan informasi yang Anda berikan, kami belum dapat menyetujui verifikasi akun perusahaan Anda saat ini."
		nextSteps = `
			<ul style="margin: 10px 0; padding-left: 20px;">
				<li>Periksa kembali kelengkapan dokumen perusahaan</li>
				<li>Pastikan informasi yang diberikan akurat</li>
				<li>Upload ulang dokumen yang diperlukan</li>
				<li>Ajukan verifikasi ulang setelah melengkapi persyaratan</li>
			</ul>
		`
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2563eb; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9fafb; }
		.status-badge { display: inline-block; padding: 8px 16px; border-radius: 5px; font-weight: bold; color: white; margin: 15px 0; }
		.button { display: inline-block; padding: 12px 24px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		.reason-box { background-color: #fff3cd; border: 1px solid #ffc107; padding: 15px; border-radius: 5px; margin: 15px 0; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Status Verifikasi Perusahaan</h1>
		</div>
		<div class="content">
			<p>Halo <strong>{{.FullName}}</strong>,</p>
			<p>{{.Message}}</p>
			
			<div style="text-align: center;">
				<span class="status-badge" style="background-color: {{.StatusColor}};">{{.StatusText}}</span>
			</div>
			
			<p><strong>Detail Perusahaan:</strong></p>
			<table style="width: 100%; border-collapse: collapse; margin: 10px 0;">
				<tr>
					<td style="padding: 8px; border-bottom: 1px solid #ddd; width: 40%;">Nama Perusahaan</td>
					<td style="padding: 8px; border-bottom: 1px solid #ddd;"><strong>{{.CompanyName}}</strong></td>
				</tr>
				<tr>
					<td style="padding: 8px; border-bottom: 1px solid #ddd;">Status Verifikasi</td>
					<td style="padding: 8px; border-bottom: 1px solid #ddd;"><strong style="color: {{.StatusColor}};">{{.StatusText}}</strong></td>
				</tr>
			</table>
			
			{{if .Reason}}
			<div class="reason-box">
				<strong>Catatan dari Admin:</strong>
				<p>{{.Reason}}</p>
			</div>
			{{end}}
			
			<p><strong>Langkah Selanjutnya:</strong></p>
			{{.NextSteps}}
			
			<div style="text-align: center;">
				<a href="https://company.karirnusantara.com" class="button">Masuk ke Dashboard</a>
			</div>
			
			<p>Jika Anda memiliki pertanyaan atau membutuhkan bantuan lebih lanjut, silakan hubungi tim dukungan kami.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		FullName    string
		CompanyName string
		Message     string
		StatusText  string
		StatusColor string
		Reason      string
		NextSteps   template.HTML
	}{
		FullName:    fullName,
		CompanyName: companyName,
		Message:     message,
		StatusText:  statusText,
		StatusColor: statusColor,
		Reason:      reason,
		NextSteps:   template.HTML(nextSteps),
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

// SendPartnerWelcomeEmail sends welcome email to new partner
func (s *Service) SendPartnerWelcomeEmail(to, partnerName, referralCode string) error {
	subject := "Selamat Bergabung Sebagai Partner Karir Nusantara!"

	tmpl := `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Selamat Bergabung</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; background-color: #f4f4f4; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background: linear-gradient(135deg, #059669 0%, #10B981 100%); color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .content { padding: 30px; }
        .welcome-box { background-color: #ECFDF5; border-left: 4px solid #10B981; padding: 20px; margin: 20px 0; border-radius: 0 8px 8px 0; }
        .referral-box { background-color: #F0FDF4; border: 2px dashed #059669; padding: 20px; margin: 20px 0; text-align: center; border-radius: 8px; }
        .referral-code { font-size: 32px; font-weight: bold; color: #059669; letter-spacing: 2px; margin: 10px 0; }
        .button { display: inline-block; background: linear-gradient(135deg, #059669 0%, #10B981 100%); color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; font-weight: bold; margin: 20px 0; }
        .footer { background-color: #374151; color: #9CA3AF; padding: 20px; text-align: center; font-size: 12px; }
        .benefits { margin: 20px 0; }
        .benefits li { margin: 10px 0; padding-left: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Selamat Bergabung!</h1>
        </div>
        <div class="content">
            <p>Halo <strong>{{.PartnerName}}</strong>,</p>
            
            <div class="welcome-box">
                <p>Terima kasih telah mendaftar sebagai Partner Karir Nusantara! Akun Anda sedang dalam proses verifikasi oleh tim kami.</p>
            </div>

            <p>Berikut adalah kode referral unik Anda:</p>
            
            <div class="referral-box">
                <p style="margin: 0; color: #6B7280;">Kode Referral Anda</p>
                <div class="referral-code">{{.ReferralCode}}</div>
                <p style="margin: 0; color: #6B7280; font-size: 14px;">Bagikan kode ini untuk mendapatkan komisi!</p>
            </div>

            <h3>Keuntungan Menjadi Partner:</h3>
            <ul class="benefits">
                <li>💰 Komisi hingga 40% dari setiap transaksi perusahaan yang Anda referensikan</li>
                <li>📊 Dashboard lengkap untuk tracking performa</li>
                <li>💳 Pencairan dana mudah dan cepat</li>
                <li>🤝 Dukungan penuh dari tim Karir Nusantara</li>
            </ul>

            <p><strong>Status Akun:</strong> Menunggu Verifikasi</p>
            <p>Tim kami akan memverifikasi akun Anda dalam 1-2 hari kerja. Anda akan menerima email konfirmasi setelah akun diaktifkan.</p>

            <div style="text-align: center;">
                <a href="https://partner.karirnusantara.com" class="button">Kunjungi Dashboard Partner</a>
            </div>

            <p>Jika ada pertanyaan, jangan ragu untuk menghubungi tim support kami.</p>
            
            <p>Salam sukses,<br><strong>Tim Karir Nusantara</strong></p>
        </div>
        <div class="footer">
            <p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
            <p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("partner_welcome").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		PartnerName  string
		ReferralCode string
	}{
		PartnerName:  partnerName,
		ReferralCode: referralCode,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}

// SendPartnerPasswordResetEmail sends password reset email to partner
func (s *Service) SendPartnerPasswordResetEmail(to, partnerName, resetLink string) error {
	subject := "Reset Password Akun Partner Karir Nusantara"

	tmpl := `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #333; background-color: #f4f4f4; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background: linear-gradient(135deg, #059669 0%, #10B981 100%); color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .content { padding: 30px; }
        .warning-box { background-color: #FEF3C7; border-left: 4px solid #F59E0B; padding: 15px; margin: 20px 0; border-radius: 0 8px 8px 0; }
        .button { display: inline-block; background: linear-gradient(135deg, #059669 0%, #10B981 100%); color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; font-weight: bold; margin: 20px 0; }
        .footer { background-color: #374151; color: #9CA3AF; padding: 20px; text-align: center; font-size: 12px; }
        .link-box { background-color: #F3F4F6; padding: 15px; border-radius: 8px; word-break: break-all; font-size: 12px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Reset Password</h1>
        </div>
        <div class="content">
            <p>Halo <strong>{{.PartnerName}}</strong>,</p>
            
            <p>Kami menerima permintaan untuk reset password akun Partner Karir Nusantara Anda.</p>

            <div style="text-align: center;">
                <a href="{{.ResetLink}}" class="button">Reset Password Sekarang</a>
            </div>

            <p>Atau salin link berikut ke browser Anda:</p>
            <div class="link-box">{{.ResetLink}}</div>

            <div class="warning-box">
                <p style="margin: 0;"><strong>⚠️ Penting:</strong></p>
                <ul style="margin: 5px 0;">
                    <li>Link ini akan kadaluarsa dalam <strong>1 jam</strong></li>
                    <li>Jika Anda tidak meminta reset password, abaikan email ini</li>
                    <li>Jangan bagikan link ini kepada siapapun</li>
                </ul>
            </div>

            <p>Jika Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini atau hubungi tim support kami jika Anda khawatir tentang keamanan akun Anda.</p>
            
            <p>Salam,<br><strong>Tim Karir Nusantara</strong></p>
        </div>
        <div class="footer">
            <p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
            <p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("partner_reset").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		PartnerName string
		ResetLink   string
	}{
		PartnerName: partnerName,
		ResetLink:   resetLink,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}

// InterviewScheduleData holds data for interview schedule email
type InterviewScheduleData struct {
	ApplicantName   string
	JobTitle        string
	CompanyName     string
	InterviewType   string
	ScheduledAt     string
	Location        string
	MeetingLink     string
	MeetingPlatform string
	ContactPerson   string
	ContactPhone    string
	Notes           string
}

// SendInterviewScheduleEmail sends interview schedule notification to candidate
func (s *Service) SendInterviewScheduleEmail(to string, data InterviewScheduleData) error {
	subject := fmt.Sprintf("Jadwal Interview - %s di %s", data.JobTitle, data.CompanyName)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%); color: white; padding: 30px 20px; text-align: center; border-radius: 10px 10px 0 0; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 30px 20px; background-color: #f9fafb; }
		.success-box { background-color: #d1fae5; padding: 20px; border-radius: 10px; text-align: center; margin: 20px 0; border-left: 4px solid #10b981; }
		.success-icon { font-size: 48px; margin-bottom: 10px; }
		.interview-details { background-color: #fff; padding: 20px; border-radius: 10px; margin: 20px 0; border: 1px solid #e5e7eb; }
		.detail-row { padding: 12px 0; border-bottom: 1px solid #e5e7eb; display: flex; }
		.detail-row:last-child { border-bottom: none; }
		.detail-label { font-weight: bold; min-width: 150px; color: #6b7280; }
		.detail-value { color: #111827; }
		.important-note { background-color: #fef3c7; padding: 15px; border-left: 4px solid #f59e0b; margin: 20px 0; border-radius: 0 8px 8px 0; }
		.button { display: inline-block; padding: 14px 28px; background-color: #2563eb; color: white; text-decoration: none; border-radius: 8px; margin: 20px 0; font-weight: bold; }
		.footer { padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f3f4f6; border-radius: 0 0 10px 10px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎯 Jadwal Interview</h1>
		</div>
		<div class="content">
			<div class="success-box">
				<div class="success-icon">✅</div>
				<h2>Selamat, {{.ApplicantName}}!</h2>
				<p>Lamaran Anda untuk posisi <strong>{{.JobTitle}}</strong> di <strong>{{.CompanyName}}</strong> telah dipilih untuk tahap interview.</p>
			</div>
			
			<h3>📋 Detail Interview:</h3>
			<div class="interview-details">
				{{if .ScheduledAt}}
				<div class="detail-row">
					<div class="detail-label">📅 Tanggal & Waktu:</div>
					<div class="detail-value">{{.ScheduledAt}}</div>
				</div>
				{{end}}
				{{if .InterviewType}}
				<div class="detail-row">
					<div class="detail-label">💼 Tipe Interview:</div>
					<div class="detail-value">{{.InterviewType}}</div>
				</div>
				{{end}}
				{{if .Location}}
				<div class="detail-row">
					<div class="detail-label">📍 Lokasi:</div>
					<div class="detail-value">{{.Location}}</div>
				</div>
				{{end}}
				{{if .MeetingPlatform}}
				<div class="detail-row">
					<div class="detail-label">💻 Platform:</div>
					<div class="detail-value">{{.MeetingPlatform}}</div>
				</div>
				{{end}}
				{{if .MeetingLink}}
				<div class="detail-row">
					<div class="detail-label">🔗 Link Meeting:</div>
					<div class="detail-value"><a href="{{.MeetingLink}}" style="color: #2563eb;">{{.MeetingLink}}</a></div>
				</div>
				{{end}}
				{{if .ContactPerson}}
				<div class="detail-row">
					<div class="detail-label">👤 Contact Person:</div>
					<div class="detail-value">{{.ContactPerson}}</div>
				</div>
				{{end}}
				{{if .ContactPhone}}
				<div class="detail-row">
					<div class="detail-label">📞 Telepon:</div>
					<div class="detail-value">{{.ContactPhone}}</div>
				</div>
				{{end}}
			</div>

			{{if .Notes}}
			<div class="important-note">
				<strong>📝 Catatan Penting:</strong>
				<p style="margin: 10px 0 0 0;">{{.Notes}}</p>
			</div>
			{{end}}

			<div class="important-note">
				<strong>⚠️ Persiapan Interview:</strong>
				<ul style="margin: 10px 0 0 0; padding-left: 20px;">
					<li>Pastikan Anda hadir tepat waktu</li>
					<li>Siapkan dokumen dan portofolio yang relevan</li>
					<li>Pelajari lebih lanjut tentang perusahaan</li>
					<li>Siapkan pertanyaan untuk interviewer</li>
				</ul>
			</div>

			<center>
				<a href="https://karirnusantara.com/dashboard/applications" class="button">Lihat Detail Lamaran</a>
			</center>

			<p style="text-align: center; color: #6b7280; margin-top: 20px;">
				Semoga sukses dengan interview Anda! 💪
			</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Karir Nusantara. All rights reserved.</p>
			<p>Email ini dikirim secara otomatis, mohon untuk tidak membalas.</p>
			<p>Jika Anda memiliki pertanyaan, silakan hubungi perusahaan melalui kontak yang tertera di atas.</p>
		</div>
	</div>
</body>
</html>
`

	t, err := template.New("interview-schedule").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendEmail(to, subject, body.String())
}
