package tickets

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository defines ticket data access methods
type Repository interface {
	// Tickets
	CreateTicket(ctx context.Context, ticket *SupportTicket) error
	GetTicketByID(ctx context.Context, id uint64) (*TicketWithDetails, error)
	ListTicketsByUser(ctx context.Context, userID uint64) ([]*TicketWithDetails, error)
	ListAllTickets(ctx context.Context, status, priority string) ([]*TicketWithDetails, error)
	UpdateTicketStatus(ctx context.Context, id uint64, status string) error
	CloseTicket(ctx context.Context, id uint64) error
	ResolveTicket(ctx context.Context, id uint64) error

	// Responses
	CreateResponse(ctx context.Context, resp *TicketResponse) error
	ListResponsesByTicket(ctx context.Context, ticketID uint64) ([]*TicketResponseWithSender, error)

	// Cooldown
	GetLastTicketTime(ctx context.Context, userID uint64) (*time.Time, error)
	CanCreateTicket(ctx context.Context, userID uint64) (bool, time.Duration, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new ticket repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// CreateTicket creates a new support ticket
func (r *repository) CreateTicket(ctx context.Context, ticket *SupportTicket) error {
	query := `
		INSERT INTO support_tickets (user_id, title, description, category, priority, status, email)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		ticket.UserID,
		ticket.Title,
		ticket.Description,
		ticket.Category,
		ticket.Priority,
		ticket.Status,
		ticket.Email,
	)

	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ticket id: %w", err)
	}

	ticket.ID = uint64(id)
	return nil
}

// GetTicketByID gets a ticket by ID with details
func (r *repository) GetTicketByID(ctx context.Context, id uint64) (*TicketWithDetails, error) {
	query := `
		SELECT 
			t.*,
			u.full_name as user_name,
			u.email as user_email,
			COALESCE(
				(SELECT message FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1), 
				t.description
			) as last_message,
			COALESCE(
				(SELECT created_at FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1),
				t.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM ticket_responses WHERE ticket_id = t.id) as response_count
		FROM support_tickets t
		JOIN users u ON u.id = t.user_id
		WHERE t.id = ?
	`

	var ticket TicketWithDetails
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}

// ListTicketsByUser lists all tickets for a user
func (r *repository) ListTicketsByUser(ctx context.Context, userID uint64) ([]*TicketWithDetails, error) {
	query := `
		SELECT 
			t.*,
			u.full_name as user_name,
			u.email as user_email,
			COALESCE(
				(SELECT message FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1), 
				t.description
			) as last_message,
			COALESCE(
				(SELECT created_at FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1),
				t.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM ticket_responses WHERE ticket_id = t.id) as response_count
		FROM support_tickets t
		JOIN users u ON u.id = t.user_id
		WHERE t.user_id = ?
		ORDER BY t.updated_at DESC
	`

	var tickets []*TicketWithDetails
	err := r.db.SelectContext(ctx, &tickets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	return tickets, nil
}

// ListAllTickets lists all tickets (for admin)
func (r *repository) ListAllTickets(ctx context.Context, status, priority string) ([]*TicketWithDetails, error) {
	query := `
		SELECT 
			t.*,
			u.full_name as user_name,
			u.email as user_email,
			COALESCE(
				(SELECT message FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1), 
				t.description
			) as last_message,
			COALESCE(
				(SELECT created_at FROM ticket_responses 
				 WHERE ticket_id = t.id 
				 ORDER BY created_at DESC LIMIT 1),
				t.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM ticket_responses WHERE ticket_id = t.id) as response_count
		FROM support_tickets t
		JOIN users u ON u.id = t.user_id
		WHERE 1=1
	`

	args := []interface{}{}

	if status != "" && status != "all" {
		query += " AND t.status = ?"
		args = append(args, status)
	}

	if priority != "" && priority != "all" {
		query += " AND t.priority = ?"
		args = append(args, priority)
	}

	query += " ORDER BY FIELD(t.priority, 'urgent', 'high', 'medium', 'low'), t.updated_at DESC"

	var tickets []*TicketWithDetails
	err := r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list all tickets: %w", err)
	}

	return tickets, nil
}

// UpdateTicketStatus updates ticket status
func (r *repository) UpdateTicketStatus(ctx context.Context, id uint64, status string) error {
	query := `UPDATE support_tickets SET status = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	return nil
}

// CloseTicket closes a ticket
func (r *repository) CloseTicket(ctx context.Context, id uint64) error {
	query := `UPDATE support_tickets SET status = 'closed', closed_at = NOW(), updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to close ticket: %w", err)
	}

	return nil
}

// ResolveTicket marks ticket as resolved
func (r *repository) ResolveTicket(ctx context.Context, id uint64) error {
	query := `UPDATE support_tickets SET status = 'resolved', resolved_at = NOW(), updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to resolve ticket: %w", err)
	}

	return nil
}

// CreateResponse creates a response to a ticket
func (r *repository) CreateResponse(ctx context.Context, resp *TicketResponse) error {
	query := `
		INSERT INTO ticket_responses (ticket_id, sender_id, sender_type, message)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		resp.TicketID,
		resp.SenderID,
		resp.SenderType,
		resp.Message,
	)

	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get response id: %w", err)
	}

	resp.ID = uint64(id)

	// Update ticket updated_at and status
	updateQuery := `UPDATE support_tickets SET updated_at = NOW(), status = ? WHERE id = ?`
	newStatus := TicketStatusPendingResponse
	if resp.SenderType == "admin" {
		newStatus = TicketStatusInProgress
	}

	_, err = r.db.ExecContext(ctx, updateQuery, newStatus, resp.TicketID)
	if err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	return nil
}

// ListResponsesByTicket lists all responses for a ticket
func (r *repository) ListResponsesByTicket(ctx context.Context, ticketID uint64) ([]*TicketResponseWithSender, error) {
	query := `
		SELECT 
			tr.*,
			COALESCE(u.full_name, 'Admin Support') as sender_name,
			COALESCE(u.email, 'support@karirnusantara.id') as sender_email
		FROM ticket_responses tr
		LEFT JOIN users u ON u.id = tr.sender_id
		WHERE tr.ticket_id = ?
		ORDER BY tr.created_at ASC
	`

	var responses []*TicketResponseWithSender
	err := r.db.SelectContext(ctx, &responses, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to list responses: %w", err)
	}

	return responses, nil
}

// GetLastTicketTime gets the time of user's last created ticket
func (r *repository) GetLastTicketTime(ctx context.Context, userID uint64) (*time.Time, error) {
	query := `
		SELECT created_at FROM support_tickets 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`

	var lastTime time.Time
	err := r.db.GetContext(ctx, &lastTime, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last ticket time: %w", err)
	}

	return &lastTime, nil
}

// CanCreateTicket checks if user can create a new ticket (cooldown check)
func (r *repository) CanCreateTicket(ctx context.Context, userID uint64) (bool, time.Duration, error) {
	lastTime, err := r.GetLastTicketTime(ctx, userID)
	if err != nil {
		return false, 0, err
	}

	if lastTime == nil {
		return true, 0, nil
	}

	elapsed := time.Since(*lastTime)
	if elapsed >= CooldownDuration {
		return true, 0, nil
	}

	remaining := CooldownDuration - elapsed
	return false, remaining, nil
}
