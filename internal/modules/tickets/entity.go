package tickets

import (
	"database/sql"
	"time"
)

// SupportTicket represents a support ticket from job seekers
type SupportTicket struct {
	ID          uint64       `json:"id" db:"id"`
	UserID      uint64       `json:"user_id" db:"user_id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Category    string       `json:"category" db:"category"`
	Priority    string       `json:"priority" db:"priority"`
	Status      string       `json:"status" db:"status"`
	Email       string       `json:"email" db:"email"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	ResolvedAt  sql.NullTime `json:"resolved_at,omitempty" db:"resolved_at"`
	ClosedAt    sql.NullTime `json:"closed_at,omitempty" db:"closed_at"`
}

// TicketWithDetails includes user info
type TicketWithDetails struct {
	SupportTicket
	UserName      string    `json:"user_name" db:"user_name"`
	UserEmail     string    `json:"user_email" db:"user_email"`
	LastMessage   string    `json:"last_message" db:"last_message"`
	LastMessageAt time.Time `json:"last_message_at" db:"last_message_at"`
	ResponseCount int       `json:"response_count" db:"response_count"`
}

// TicketResponse represents a response/message in a ticket
type TicketResponse struct {
	ID         uint64    `json:"id" db:"id"`
	TicketID   uint64    `json:"ticket_id" db:"ticket_id"`
	SenderID   uint64    `json:"sender_id" db:"sender_id"`
	SenderType string    `json:"sender_type" db:"sender_type"` // "user" or "admin"
	Message    string    `json:"message" db:"message"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TicketResponseWithSender includes sender information
type TicketResponseWithSender struct {
	TicketResponse
	SenderName  string `json:"sender_name" db:"sender_name"`
	SenderEmail string `json:"sender_email" db:"sender_email"`
}

// CreateTicketRequest represents request to create new ticket
type CreateTicketRequest struct {
	Title       string `json:"title" validate:"required,min=5,max=100"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
	Category    string `json:"category" validate:"required,oneof=account cv-builder job-applications profile payment search-filter notification technical feature-request other"`
	Priority    string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	Email       string `json:"email" validate:"required,email"`
}

// AddResponseRequest represents request to add response to ticket
type AddResponseRequest struct {
	Message string `json:"message" validate:"required,min=1,max=500"`
}

// UpdateTicketStatusRequest represents request to update ticket status
type UpdateTicketStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress pending_response resolved closed"`
}

// TicketCooldown represents cooldown information for user
type TicketCooldown struct {
	UserID            uint64    `db:"user_id"`
	LastTicketCreated time.Time `db:"last_ticket_created"`
}

// Ticket status constants
const (
	TicketStatusOpen            = "open"
	TicketStatusInProgress      = "in_progress"
	TicketStatusPendingResponse = "pending_response"
	TicketStatusResolved        = "resolved"
	TicketStatusClosed          = "closed"
)

// Ticket priority constants
const (
	TicketPriorityLow    = "low"
	TicketPriorityMedium = "medium"
	TicketPriorityHigh   = "high"
	TicketPriorityUrgent = "urgent"
)

// Ticket category constants
const (
	TicketCategoryAccount         = "account"
	TicketCategoryCVBuilder       = "cv-builder"
	TicketCategoryJobApplications = "job-applications"
	TicketCategoryProfile         = "profile"
	TicketCategoryPayment         = "payment"
	TicketCategorySearchFilter    = "search-filter"
	TicketCategoryNotification    = "notification"
	TicketCategoryTechnical       = "technical"
	TicketCategoryFeatureRequest  = "feature-request"
	TicketCategoryOther           = "other"
)

// CooldownDuration is the minimum time between ticket creations (1 hour)
const CooldownDuration = 1 * time.Hour

// GetCategoryLabel returns human-readable category label
func GetCategoryLabel(category string) string {
	labels := map[string]string{
		"account":          "Akun & Login",
		"cv-builder":       "CV Builder",
		"job-applications": "Lamaran Kerja",
		"profile":          "Profil",
		"payment":          "Pembayaran",
		"search-filter":    "Pencarian & Filter",
		"notification":     "Notifikasi",
		"technical":        "Masalah Teknis",
		"feature-request":  "Permintaan Fitur",
		"other":            "Lainnya",
	}
	if label, ok := labels[category]; ok {
		return label
	}
	return category
}

// GetPriorityLabel returns human-readable priority label
func GetPriorityLabel(priority string) string {
	labels := map[string]string{
		"low":    "Rendah",
		"medium": "Sedang",
		"high":   "Tinggi",
		"urgent": "Mendesak",
	}
	if label, ok := labels[priority]; ok {
		return label
	}
	return priority
}

// GetStatusLabel returns human-readable status label
func GetStatusLabel(status string) string {
	labels := map[string]string{
		"open":             "Terbuka",
		"in_progress":      "Sedang Diproses",
		"pending_response": "Menunggu Respons",
		"resolved":         "Terselesaikan",
		"closed":           "Ditutup",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}
