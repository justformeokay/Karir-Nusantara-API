package invoice

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PaymentInvoiceData holds data for generating payment invoice
type PaymentInvoiceData struct {
	InvoiceNumber   string
	PaymentID       uint64
	CompanyName     string
	CompanyEmail    string
	CompanyAddress  string
	Amount          int64
	PaymentDate     time.Time
	ConfirmedDate   time.Time
	Description     string
	AdminNote       string
}

// Service handles invoice generation
type Service struct {
	invoiceDir string
}

// NewService creates a new invoice service
func NewService(invoiceDir string) *Service {
	// Ensure invoice directory exists
	if err := os.MkdirAll(invoiceDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create invoice directory: %v\n", err)
	}
	
	return &Service{
		invoiceDir: invoiceDir,
	}
}

// GeneratePaymentInvoice generates a PDF invoice for payment confirmation
func (s *Service) GeneratePaymentInvoice(data *PaymentInvoiceData) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 24)

	// Header - Company Logo/Name
	pdf.SetTextColor(37, 99, 235) // Primary blue color
	pdf.Cell(0, 15, "Karir Nusantara")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 5, "Platform Rekrutmen Terpercaya")
	pdf.Ln(3)
	pdf.Cell(0, 5, "Email: info@karirnusantara.com | Website: www.karirnusantara.com")
	pdf.Ln(15)

	// Invoice Title
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "INVOICE PEMBAYARAN")
	pdf.Ln(15)

	// Invoice Info Box
	pdf.SetFillColor(240, 248, 255)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Rect(10, pdf.GetY(), 190, 30, "FD")
	
	currentY := pdf.GetY() + 5
	pdf.SetY(currentY)
	
	// Left column
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 5, "No. Invoice:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(45, 5, data.InvoiceNumber)
	
	// Right column
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 5, "Tanggal Invoice:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, data.ConfirmedDate.Format("02 January 2006"))
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 5, "ID Pembayaran:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(45, 5, fmt.Sprintf("#%d", data.PaymentID))
	
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 5, "Status:")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(34, 197, 94) // Green
	pdf.Cell(0, 5, "LUNAS")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 5, "Tanggal Bayar:")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, data.PaymentDate.Format("02 January 2006 15:04"))
	pdf.Ln(15)

	// Company Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Ditagihkan Kepada:")
	pdf.Ln(6)
	
	pdf.SetFillColor(250, 250, 250)
	pdf.Rect(10, pdf.GetY(), 190, 25, "FD")
	
	currentY = pdf.GetY() + 4
	pdf.SetY(currentY)
	
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 6, data.CompanyName)
	pdf.Ln(5)
	
	pdf.SetFont("Arial", "", 10)
	if data.CompanyEmail != "" {
		pdf.Cell(0, 5, "Email: "+data.CompanyEmail)
		pdf.Ln(5)
	}
	if data.CompanyAddress != "" {
		pdf.MultiCell(0, 5, "Alamat: "+data.CompanyAddress, "", "L", false)
	}
	pdf.Ln(10)

	// Payment Details Table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Rincian Pembayaran:")
	pdf.Ln(8)

	// Table Header
	pdf.SetFillColor(37, 99, 235)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)
	
	pdf.CellFormat(100, 10, "Deskripsi", "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 10, "Tanggal", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Jumlah", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Table Row
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	description := data.Description
	if description == "" {
		description = "Pembayaran Kuota Job Posting"
	}
	
	pdf.CellFormat(100, 10, description, "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 10, data.PaymentDate.Format("02 Jan 2006"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, formatRupiah(data.Amount), "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Subtotal
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(145, 8, "Subtotal", "1", 0, "R", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 8, formatRupiah(data.Amount), "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Tax (0%)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(145, 8, "PPN (0%)", "1", 0, "R", true, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 8, "Rp 0", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Total
	pdf.SetFillColor(37, 99, 235)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(145, 12, "TOTAL", "1", 0, "R", true, 0, "")
	pdf.CellFormat(45, 12, formatRupiah(data.Amount), "1", 0, "R", true, 0, "")
	pdf.Ln(15)

	// Admin Note (if any)
	if data.AdminNote != "" {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "Catatan:")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 9)
		pdf.SetFillColor(255, 251, 235)
		pdf.Rect(10, pdf.GetY(), 190, 15, "FD")
		currentY = pdf.GetY() + 4
		pdf.SetY(currentY)
		pdf.MultiCell(0, 5, data.AdminNote, "", "L", false)
		pdf.Ln(10)
	}

	// Payment Info
	pdf.SetTextColor(100, 100, 100)
	pdf.SetFont("Arial", "I", 9)
	pdf.MultiCell(0, 5, "Pembayaran ini telah dikonfirmasi oleh tim Karir Nusantara. Kuota job posting Anda telah ditambahkan dan siap digunakan.", "", "L", false)
	pdf.Ln(5)

	// Footer
	pdf.SetY(-30)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(3)
	
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.Cell(0, 4, "Invoice ini digenerate secara otomatis oleh sistem Karir Nusantara")
	pdf.Ln(3)
	pdf.Cell(0, 4, fmt.Sprintf("Dicetak pada: %s", time.Now().Format("02 January 2006 15:04:05")))
	pdf.Ln(3)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 4, "Terima kasih atas kepercayaan Anda menggunakan Karir Nusantara")

	// Generate filename
	filename := fmt.Sprintf("invoice_%s_%d.pdf", 
		data.ConfirmedDate.Format("20060102"), 
		data.PaymentID)
	filepath := filepath.Join(s.invoiceDir, filename)

	// Save PDF
	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	return filepath, nil
}

// formatRupiah formats an amount to Indonesian Rupiah format
func formatRupiah(amount int64) string {
	// Simple formatting (you can enhance this)
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

// GetInvoicePath returns the full path to an invoice file
func (s *Service) GetInvoicePath(filename string) string {
	return filepath.Join(s.invoiceDir, filename)
}
