package chat

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository defines chat data access methods
type Repository interface {
	// Conversations
	CreateConversation(ctx context.Context, conv *Conversation) error
	GetConversationByID(ctx context.Context, id uint64) (*ConversationWithDetails, error)
	ListConversationsByCompany(ctx context.Context, companyID uint64) ([]*ConversationWithDetails, error)
	ListAllConversations(ctx context.Context) ([]*ConversationWithDetails, error)
	UpdateConversationStatus(ctx context.Context, id uint64, status string) error
	HasActiveConversation(ctx context.Context, companyID uint64) (bool, error)
	
	// Messages
	CreateMessage(ctx context.Context, msg *ChatMessage) error
	ListMessagesByConversation(ctx context.Context, conversationID uint64) ([]*ChatMessageWithSender, error)
	MarkMessagesAsRead(ctx context.Context, conversationID uint64, userID uint64) error
	GetUnreadCount(ctx context.Context, conversationID uint64, userID uint64) (int, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new chat repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// CreateConversation creates a new conversation
func (r *repository) CreateConversation(ctx context.Context, conv *Conversation) error {
	query := `
		INSERT INTO conversations (company_id, title, subject, category, status)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		conv.CompanyID,
		conv.Title,
		conv.Subject,
		conv.Category,
		conv.Status,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get conversation id: %w", err)
	}
	
	conv.ID = uint64(id)
	return nil
}

// GetConversationByID gets a conversation by ID with details
func (r *repository) GetConversationByID(ctx context.Context, id uint64) (*ConversationWithDetails, error)  {
	query := `
		SELECT 
			c.*,
			COALESCE(co.company_name, u.full_name) as company_name,
			COALESCE(
				(SELECT message FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1), 
				''
			) as last_message,
			COALESCE(
				(SELECT created_at FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1),
				c.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM chat_messages 
			 WHERE conversation_id = c.id AND is_read = FALSE
			) as unread_count
		FROM conversations c
		JOIN users u ON u.id = c.company_id
		LEFT JOIN companies co ON co.user_id = u.id
		WHERE c.id = ?
	`
	
	var conv ConversationWithDetails
	err := r.db.GetContext(ctx, &conv, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	
	return &conv, nil
}

// ListConversationsByCompany lists all conversations for a company
func (r *repository) ListConversationsByCompany(ctx context.Context, companyID uint64) ([]*ConversationWithDetails, error) {
	query := `
		SELECT 
			c.*,
			COALESCE(co.company_name, u.full_name) as company_name,
			COALESCE(
				(SELECT message FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1), 
				''
			) as last_message,
			COALESCE(
				(SELECT created_at FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1),
				c.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM chat_messages 
			 WHERE conversation_id = c.id AND is_read = FALSE AND sender_type = 'admin'
			) as unread_count
		FROM conversations c
		JOIN users u ON u.id = c.company_id
		LEFT JOIN companies co ON co.user_id = u.id
		WHERE c.company_id = ?
		ORDER BY last_message_at DESC
	`
	
	var convs []*ConversationWithDetails
	err := r.db.SelectContext(ctx, &convs, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	
	return convs, nil
}

// ListAllConversations lists all conversations (for admin)
func (r *repository) ListAllConversations(ctx context.Context) ([]*ConversationWithDetails, error) {
	query := `
		SELECT 
			c.*,
			COALESCE(co.company_name, u.full_name) as company_name,
			COALESCE(
				(SELECT message FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1), 
				''
			) as last_message,
			COALESCE(
				(SELECT created_at FROM chat_messages 
				 WHERE conversation_id = c.id 
				 ORDER BY created_at DESC LIMIT 1),
				c.created_at
			) as last_message_at,
			(SELECT COUNT(*) FROM chat_messages 
			 WHERE conversation_id = c.id AND is_read = FALSE AND sender_type = 'company'
			) as unread_count
		FROM conversations c
		JOIN users u ON u.id = c.company_id
		LEFT JOIN companies co ON co.user_id = u.id
		ORDER BY last_message_at DESC
	`
	
	var convs []*ConversationWithDetails
	err := r.db.SelectContext(ctx, &convs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all conversations: %w", err)
	}
	
	return convs, nil
}

// UpdateConversationStatus updates conversation status
func (r *repository) UpdateConversationStatus(ctx context.Context, id uint64, status string) error {
	query := `
		UPDATE conversations 
		SET status = ?, 
		    closed_at = CASE WHEN ? = 'closed' THEN NOW() ELSE closed_at END
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(ctx, query, status, status, id)
	if err != nil {
		return fmt.Errorf("failed to update conversation status: %w", err)
	}
	
	return nil
}

// CreateMessage creates a new message
func (r *repository) CreateMessage(ctx context.Context, msg *ChatMessage) error {
	query := `
		INSERT INTO chat_messages (conversation_id, sender_id, sender_type, message, attachment_url, attachment_type, attachment_filename, is_read)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query,
		msg.ConversationID,
		msg.SenderID,
		msg.SenderType,
		msg.Message,
		msg.AttachmentURL,
		msg.AttachmentType,
		msg.AttachmentFilename,
		msg.IsRead,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	msg.ID = uint64(id)
	
	// Update conversation updated_at
	_, _ = r.db.ExecContext(ctx, "UPDATE conversations SET updated_at = NOW() WHERE id = ?", msg.ConversationID)
	
	return nil
}

// ListMessagesByConversation lists all messages in a conversation
func (r *repository) ListMessagesByConversation(ctx context.Context, conversationID uint64) ([]*ChatMessageWithSender, error) {
	query := `
		SELECT 
			m.*,
			u.full_name as sender_name,
			u.email as sender_email
		FROM chat_messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
		ORDER BY m.created_at ASC
	`
	
	var messages []*ChatMessageWithSender
	err := r.db.SelectContext(ctx, &messages, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	
	return messages, nil
}

// MarkMessagesAsRead marks all messages as read for a user
func (r *repository) MarkMessagesAsRead(ctx context.Context, conversationID uint64, userID uint64) error {
	query := `
		UPDATE chat_messages 
		SET is_read = TRUE 
		WHERE conversation_id = ? 
		  AND sender_id != ? 
		  AND is_read = FALSE
	`
	
	_, err := r.db.ExecContext(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}
	
	return nil
}

// GetUnreadCount gets unread message count for a user in a conversation
func (r *repository) GetUnreadCount(ctx context.Context, conversationID uint64, userID uint64) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM chat_messages 
		WHERE conversation_id = ? 
		  AND sender_id != ? 
		  AND is_read = FALSE
	`
	
	var count int
	err := r.db.GetContext(ctx, &count, query, conversationID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	
	return count, nil
}

// HasActiveConversation checks if company has any active (open/in_progress) conversation
func (r *repository) HasActiveConversation(ctx context.Context, companyID uint64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM conversations 
		WHERE company_id = ? 
		  AND status IN ('open', 'in_progress')
	`
	
	var count int
	err := r.db.GetContext(ctx, &count, query, companyID)
	if err != nil {
		return false, fmt.Errorf("failed to check active conversation: %w", err)
	}
	
	return count > 0, nil
}
