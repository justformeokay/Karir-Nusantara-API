-- Add attachment support to chat messages
-- Supports image files and voice recordings

ALTER TABLE chat_messages 
ADD COLUMN attachment_url VARCHAR(500) NULL AFTER message,
ADD COLUMN attachment_type ENUM('image', 'audio') NULL AFTER attachment_url,
ADD COLUMN attachment_filename VARCHAR(255) NULL AFTER attachment_type;

-- Add index for searching messages with attachments
ALTER TABLE chat_messages 
ADD INDEX idx_attachment_type (attachment_type);
