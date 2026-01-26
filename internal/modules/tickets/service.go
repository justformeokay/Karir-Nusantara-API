package tickets

import (
	"context"
	"fmt"
	"time"
)

// Service defines ticket business logic
type Service interface {
	// For Job Seekers
	CreateTicket(ctx context.Context, userID uint64, req *CreateTicketRequest) (*TicketWithDetails, error)
	GetMyTickets(ctx context.Context, userID uint64) ([]*TicketWithDetails, error)
	GetTicket(ctx context.Context, ticketID uint64, userID uint64) (*TicketWithDetails, []*TicketResponseWithSender, error)
	AddResponse(ctx context.Context, ticketID uint64, userID uint64, senderType string, req *AddResponseRequest) (*TicketResponseWithSender, error)
	CheckCooldown(ctx context.Context, userID uint64) (bool, time.Duration, error)

	// For Admin
	GetAllTickets(ctx context.Context, status, priority string) ([]*TicketWithDetails, error)
	UpdateTicketStatus(ctx context.Context, ticketID uint64, req *UpdateTicketStatusRequest) error
	CloseTicket(ctx context.Context, ticketID uint64) error
	ResolveTicket(ctx context.Context, ticketID uint64) error
}

type service struct {
	repo Repository
}

// NewService creates a new ticket service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateTicket creates a new support ticket
func (s *service) CreateTicket(ctx context.Context, userID uint64, req *CreateTicketRequest) (*TicketWithDetails, error) {
	// Check cooldown
	canCreate, remaining, err := s.repo.CanCreateTicket(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !canCreate {
		minutes := int(remaining.Minutes())
		seconds := int(remaining.Seconds()) % 60
		return nil, fmt.Errorf("Anda hanya bisa membuat 1 ticket per jam. Coba lagi dalam %d menit %d detik", minutes, seconds)
	}

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = TicketPriorityMedium
	}

	ticket := &SupportTicket{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Priority:    priority,
		Status:      TicketStatusOpen,
		Email:       req.Email,
	}

	err = s.repo.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, err
	}

	// Create initial response from user's description
	initialResp := &TicketResponse{
		TicketID:   ticket.ID,
		SenderID:   userID,
		SenderType: "user",
		Message:    req.Description,
	}

	err = s.repo.CreateResponse(ctx, initialResp)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to create initial response: %v\n", err)
	}

	// Reset status to open after initial response
	_ = s.repo.UpdateTicketStatus(ctx, ticket.ID, TicketStatusOpen)

	// Get ticket with details
	return s.repo.GetTicketByID(ctx, ticket.ID)
}

// GetMyTickets gets all tickets for a user
func (s *service) GetMyTickets(ctx context.Context, userID uint64) ([]*TicketWithDetails, error) {
	return s.repo.ListTicketsByUser(ctx, userID)
}

// GetTicket gets a ticket with responses
func (s *service) GetTicket(ctx context.Context, ticketID uint64, userID uint64) (*TicketWithDetails, []*TicketResponseWithSender, error) {
	// Get ticket
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, nil, err
	}

	if ticket == nil {
		return nil, nil, fmt.Errorf("ticket not found")
	}

	// Verify user owns the ticket or is admin
	// For now, we check ownership - admin check should be in handler
	if ticket.UserID != userID {
		return nil, nil, fmt.Errorf("unauthorized to view this ticket")
	}

	// Get responses
	responses, err := s.repo.ListResponsesByTicket(ctx, ticketID)
	if err != nil {
		return nil, nil, err
	}

	return ticket, responses, nil
}

// GetTicketForAdmin gets a ticket with responses for admin
func (s *service) GetTicketForAdmin(ctx context.Context, ticketID uint64) (*TicketWithDetails, []*TicketResponseWithSender, error) {
	// Get ticket
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, nil, err
	}

	if ticket == nil {
		return nil, nil, fmt.Errorf("ticket not found")
	}

	// Get responses
	responses, err := s.repo.ListResponsesByTicket(ctx, ticketID)
	if err != nil {
		return nil, nil, err
	}

	return ticket, responses, nil
}

// AddResponse adds a response to a ticket
func (s *service) AddResponse(ctx context.Context, ticketID uint64, userID uint64, senderType string, req *AddResponseRequest) (*TicketResponseWithSender, error) {
	// Verify ticket exists
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	// Check if ticket is closed
	if ticket.Status == TicketStatusClosed {
		return nil, fmt.Errorf("ticket sudah ditutup dan tidak bisa menerima respons baru")
	}

	// For user, verify ownership
	if senderType == "user" && ticket.UserID != userID {
		return nil, fmt.Errorf("unauthorized to respond to this ticket")
	}

	resp := &TicketResponse{
		TicketID:   ticketID,
		SenderID:   userID,
		SenderType: senderType,
		Message:    req.Message,
	}

	err = s.repo.CreateResponse(ctx, resp)
	if err != nil {
		return nil, err
	}

	// Return response with sender info
	return &TicketResponseWithSender{
		TicketResponse: *resp,
		SenderName:     "Anda",
		SenderEmail:    ticket.Email,
	}, nil
}

// CheckCooldown checks if user can create a new ticket
func (s *service) CheckCooldown(ctx context.Context, userID uint64) (bool, time.Duration, error) {
	return s.repo.CanCreateTicket(ctx, userID)
}

// GetAllTickets gets all tickets for admin
func (s *service) GetAllTickets(ctx context.Context, status, priority string) ([]*TicketWithDetails, error) {
	return s.repo.ListAllTickets(ctx, status, priority)
}

// UpdateTicketStatus updates ticket status
func (s *service) UpdateTicketStatus(ctx context.Context, ticketID uint64, req *UpdateTicketStatusRequest) error {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket == nil {
		return fmt.Errorf("ticket not found")
	}

	if req.Status == TicketStatusResolved {
		return s.repo.ResolveTicket(ctx, ticketID)
	}

	if req.Status == TicketStatusClosed {
		return s.repo.CloseTicket(ctx, ticketID)
	}

	return s.repo.UpdateTicketStatus(ctx, ticketID, req.Status)
}

// CloseTicket closes a ticket
func (s *service) CloseTicket(ctx context.Context, ticketID uint64) error {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket == nil {
		return fmt.Errorf("ticket not found")
	}

	return s.repo.CloseTicket(ctx, ticketID)
}

// ResolveTicket marks ticket as resolved
func (s *service) ResolveTicket(ctx context.Context, ticketID uint64) error {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket == nil {
		return fmt.Errorf("ticket not found")
	}

	return s.repo.ResolveTicket(ctx, ticketID)
}
