package chat

import (
	"database/sql"
	"time"
)

// Conversation represents a chat conversation between company and admin
type Conversation struct {
	ID         uint64         `json:"id" db:"id"`
	CompanyID  uint64         `json:"company_id" db:"company_id"`
	Title      string         `json:"title" db:"title"`
	Subject    string         `json:"subject" db:"subject"`
	Category   string         `json:"category" db:"category"`
	Status     string         `json:"status" db:"status"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	ClosedAt   sql.NullTime   `json:"closed_at,omitempty" db:"closed_at"`
}

// ConversationWithDetails includes additional details
type ConversationWithDetails struct {
	Conversation
	CompanyName    string    `json:"company_name" db:"company_name"`
	LastMessage    string    `json:"last_message" db:"last_message"`
	LastMessageAt  time.Time `json:"last_message_at" db:"last_message_at"`
	UnreadCount    int       `json:"unread_count" db:"unread_count"`
}

// ChatMessage represents a message in a conversation
type ChatMessage struct {
	ID                 uint64         `json:"id" db:"id"`
	ConversationID     uint64         `json:"conversation_id" db:"conversation_id"`
	SenderID           uint64         `json:"sender_id" db:"sender_id"`
	SenderType         string         `json:"sender_type" db:"sender_type"`
	Message            string         `json:"message" db:"message"`
	AttachmentURL      sql.NullString `json:"attachment_url,omitempty" db:"attachment_url"`
	AttachmentType     sql.NullString `json:"attachment_type,omitempty" db:"attachment_type"`
	AttachmentFilename sql.NullString `json:"attachment_filename,omitempty" db:"attachment_filename"`
	IsRead             bool           `json:"is_read" db:"is_read"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}

// ChatMessageWithSender includes sender information
type ChatMessageWithSender struct {
	ChatMessage
	SenderName  string `json:"sender_name" db:"sender_name"`
	SenderEmail string `json:"sender_email" db:"sender_email"`
}

// CreateConversationRequest represents request to create new conversation
type CreateConversationRequest struct {
	Subject  string `json:"subject" validate:"required,min=5,max=255"`
	Category string `json:"category" validate:"required,oneof=complaint helpdesk general urgent"`
}

// SendMessageRequest represents request to send a message
type SendMessageRequest struct {
	Message            string `json:"message"`
	AttachmentURL      string `json:"attachment_url,omitempty"`
	AttachmentType     string `json:"attachment_type,omitempty"`
	AttachmentFilename string `json:"attachment_filename,omitempty"`
}

// UpdateConversationStatusRequest represents request to update conversation status
type UpdateConversationStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress resolved closed"`
}
