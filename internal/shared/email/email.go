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

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
