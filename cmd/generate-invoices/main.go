package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/karirnusantara/api/internal/config"
	"github.com/karirnusantara/api/internal/database"
	"github.com/karirnusantara/api/internal/shared/invoice"
)

type Payment struct {
	ID          uint64
	CompanyID   uint64
	Amount      int64
	Status      string
	ConfirmedAt sql.NullTime
	SubmittedAt time.Time
}

type Company struct {
	CompanyName  string
	CompanyEmail string
}

func main() {
	log.Println("Starting invoice generation for confirmed payments...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize invoice service
	invoiceSvc := invoice.NewService("./docs/invoices")

	// Get all confirmed payments
	query := `
		SELECT p.id, p.company_id, p.amount, p.status, p.confirmed_at, p.submitted_at
		FROM payments p
		WHERE p.status = 'confirmed'
		ORDER BY p.id
	`

	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to query payments: %v", err)
	}
	defer rows.Close()

	generated := 0
	skipped := 0

	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.CompanyID,
			&payment.Amount,
			&payment.Status,
			&payment.ConfirmedAt,
			&payment.SubmittedAt,
		)
		if err != nil {
			log.Printf("Failed to scan payment: %v", err)
			continue
		}

		// Check if invoice already exists
		var dateStr string
		if payment.ConfirmedAt.Valid {
			dateStr = payment.ConfirmedAt.Time.Format("20060102")
		} else {
			dateStr = payment.SubmittedAt.Format("20060102")
		}
		
		invoicePath := fmt.Sprintf("./docs/invoices/invoice_%s_%d.pdf", dateStr, payment.ID)
		if _, err := os.Stat(invoicePath); err == nil {
			log.Printf("Invoice already exists for payment #%d, skipping", payment.ID)
			skipped++
			continue
		}

		// Get company info
		var company Company
		companyQuery := `
			SELECT 
				COALESCE(c.company_name, u.full_name) as company_name,
				u.email as company_email
			FROM users u
			LEFT JOIN companies c ON c.user_id = u.id
			WHERE u.id = ?
		`
		err = db.QueryRowContext(context.Background(), companyQuery, payment.CompanyID).Scan(
			&company.CompanyName,
			&company.CompanyEmail,
		)
		if err != nil {
			log.Printf("Failed to get company info for payment #%d: %v", payment.ID, err)
			continue
		}

		// Generate invoice number
		var invoiceDate time.Time
		if payment.ConfirmedAt.Valid {
			invoiceDate = payment.ConfirmedAt.Time
		} else {
			invoiceDate = payment.SubmittedAt
		}
		
		invoiceNumber := fmt.Sprintf("INV/%s/%05d", invoiceDate.Format("2006/01"), payment.ID)

		// Create invoice data
		invoiceData := &invoice.PaymentInvoiceData{
			InvoiceNumber:  invoiceNumber,
			PaymentID:      payment.ID,
			CompanyName:    company.CompanyName,
			CompanyEmail:   company.CompanyEmail,
			Amount:         payment.Amount,
			PaymentDate:    payment.SubmittedAt,
			ConfirmedDate:  invoiceDate,
			Description:    "Pembayaran Kuota Lowongan Kerja",
			AdminNote:      "",
		}

		// Generate invoice
		pdfPath, err := invoiceSvc.GeneratePaymentInvoice(invoiceData)
		if err != nil {
			log.Printf("Failed to generate invoice for payment #%d: %v", payment.ID, err)
			continue
		}

		log.Printf("✓ Generated invoice for payment #%d: %s", payment.ID, pdfPath)
		generated++
	}

	log.Printf("\n=== Summary ===")
	log.Printf("Generated: %d invoices", generated)
	log.Printf("Skipped: %d invoices (already exist)", skipped)
	log.Printf("Total: %d", generated+skipped)
}
