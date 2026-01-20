package company

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FileService handles file operations for company documents
type FileService struct {
	basePath string
}

// NewFileService creates a new file service
func NewFileService(basePath string) *FileService {
	return &FileService{
		basePath: basePath,
	}
}

// SaveCompanyDocument saves a company document and returns the relative file path
func (fs *FileService) SaveCompanyDocument(companyID uint64, docType string, file io.Reader, originalFilename string) (string, error) {
	// Create company directory
	companyDir := filepath.Join(fs.basePath, fmt.Sprintf("%d", companyID))
	if err := os.MkdirAll(companyDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create company directory: %w", err)
	}

	// Generate unique filename with timestamp
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", docType, timestamp, ext)
	filePath := filepath.Join(companyDir, filename)

	// Create file
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Copy file content
	if _, err := io.Copy(f, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative path from docs/companies
	relativePath := filepath.Join(fmt.Sprintf("%d", companyID), filename)
	return relativePath, nil
}

// DeleteCompanyDocument deletes a company document
func (fs *FileService) DeleteCompanyDocument(filePath string) error {
	fullPath := filepath.Join(fs.basePath, filePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetCompanyDocumentPath returns the full path to a document
func (fs *FileService) GetCompanyDocumentPath(filePath string) string {
	return filepath.Join(fs.basePath, filePath)
}

// ValidateImageFile validates if the file is a valid image
func ValidateImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	validExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// ValidateImageFileSize validates if the file size is acceptable
func ValidateImageFileSize(size int64, maxSizeBytes int64) bool {
	return size <= maxSizeBytes
}
