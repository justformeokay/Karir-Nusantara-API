package chat

import (
	"context"
	"fmt"
)

// Service defines chat business logic
type Service interface {
	// For Companies
	CreateConversation(ctx context.Context, companyID uint64, req *CreateConversationRequest) (*ConversationWithDetails, error)
	GetMyConversations(ctx context.Context, companyID uint64) ([]*ConversationWithDetails, error)
	GetConversation(ctx context.Context, conversationID uint64, userID uint64) (*ConversationWithDetails, []*ChatMessageWithSender, error)
	SendMessage(ctx context.Context, conversationID uint64, userID uint64, senderType string, req *SendMessageRequest) (*ChatMessageWithSender, error)
	
	// For Admin
	GetAllConversations(ctx context.Context) ([]*ConversationWithDetails, error)
	UpdateConversationStatus(ctx context.Context, conversationID uint64, req *UpdateConversationStatusRequest) error
}

type service struct {
	repo Repository
}

// NewService creates a new chat service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateConversation creates a new conversation
func (s *service) CreateConversation(ctx context.Context, companyID uint64, req *CreateConversationRequest) (*ConversationWithDetails, error) {
	// Check if company has active conversation (ticketing mode)
	hasActive, err := s.repo.HasActiveConversation(ctx, companyID)
	if err != nil {
		return nil, err
	}
	
	if hasActive {
		return nil, fmt.Errorf("Anda masih memiliki percakapan aktif. Silakan tutup percakapan tersebut terlebih dahulu sebelum membuat percakapan baru")
	}
	
	// Generate title from category and subject
	categoryLabels := map[string]string{
		"complaint": "Komplain",
		"helpdesk":  "Help Desk",
		"general":   "Pertanyaan",
		"urgent":    "Urgent",
	}
	
	label := categoryLabels[req.Category]
	if label == "" {
		label = "Pertanyaan"
	}
	
	title := fmt.Sprintf("%s: %s", label, req.Subject)
	if len(title) > 255 {
		title = title[:252] + "..."
	}
	
	conv := &Conversation{
		CompanyID: companyID,
		Title:     title,
		Subject:   req.Subject,
		Category:  req.Category,
		Status:    "open",
	}
	
	err = s.repo.CreateConversation(ctx, conv)
	if err != nil {
		return nil, err
	}
	
	// Get conversation with details
	return s.repo.GetConversationByID(ctx, conv.ID)
}

// GetMyConversations gets all conversations for a company
func (s *service) GetMyConversations(ctx context.Context, companyID uint64) ([]*ConversationWithDetails, error) {
	return s.repo.ListConversationsByCompany(ctx, companyID)
}

// GetConversation gets a conversation with messages
func (s *service) GetConversation(ctx context.Context, conversationID uint64, userID uint64) (*ConversationWithDetails, []*ChatMessageWithSender, error) {
	// Get conversation
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, nil, err
	}
	
	if conv == nil {
		return nil, nil, fmt.Errorf("conversation not found")
	}
	
	// Get messages
	messages, err := s.repo.ListMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, nil, err
	}
	
	// Mark messages as read
	_ = s.repo.MarkMessagesAsRead(ctx, conversationID, userID)
	
	return conv, messages, nil
}

// SendMessage sends a message in a conversation
func (s *service) SendMessage(ctx context.Context, conversationID uint64, userID uint64, senderType string, req *SendMessageRequest) (*ChatMessageWithSender, error) {
	// Verify conversation exists
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	
	if conv == nil {
		return nil, fmt.Errorf("conversation not found")
	}
	
	// Check if conversation is closed
	if conv.Status == "closed" {
		return nil, fmt.Errorf("conversation is closed")
	}
	
	// Validate: message or attachment required
	if req.Message == "" && req.AttachmentURL == "" {
		return nil, fmt.Errorf("message or attachment required")
	}
	
	// Create message
	msg := &ChatMessage{
		ConversationID: conversationID,
		SenderID:       userID,
		SenderType:     senderType,
		Message:        req.Message,
		IsRead:         false,
	}
	
	// Add attachment if provided
	if req.AttachmentURL != "" {
		msg.AttachmentURL.String = req.AttachmentURL
		msg.AttachmentURL.Valid = true
		msg.AttachmentType.String = req.AttachmentType
		msg.AttachmentType.Valid = true
		msg.AttachmentFilename.String = req.AttachmentFilename
		msg.AttachmentFilename.Valid = true
	}
	
	err = s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	
	// Get message with sender info
	messages, err := s.repo.ListMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	
	// Find the message we just created
	for _, m := range messages {
		if m.ID == msg.ID {
			return m, nil
		}
	}
	
	return nil, fmt.Errorf("failed to retrieve sent message")
}

// GetAllConversations gets all conversations (for admin)
func (s *service) GetAllConversations(ctx context.Context) ([]*ConversationWithDetails, error) {
	return s.repo.ListAllConversations(ctx)
}

// UpdateConversationStatus updates conversation status
func (s *service) UpdateConversationStatus(ctx context.Context, conversationID uint64, req *UpdateConversationStatusRequest) error {
	return s.repo.UpdateConversationStatus(ctx, conversationID, req.Status)
}
