package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/karirnusantara/api/internal/config"
	apperrors "github.com/karirnusantara/api/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the auth service interface
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, userID uint64, refreshToken string) error
	GetCurrentUser(ctx context.Context, userID uint64) (*UserWithCompanyResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) (*UserResponse, error)
	ValidateAccessToken(tokenString string) (*TokenClaims, error)
}

// service implements Service
type service struct {
	repo   Repository
	config *config.JWTConfig
}

// NewService creates a new auth service
func NewService(repo Repository, cfg *config.JWTConfig) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

// Register creates a new user account
func (s *service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check email", err)
	}
	if exists {
		return nil, apperrors.NewDuplicateEntryError("Email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to hash password", err)
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		IsActive:     true,
		IsVerified:   false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("Failed to create user", err)
	}

	// Generate tokens
	return s.generateAuthResponse(ctx, user)
}

// Login authenticates a user
func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil {
		return nil, apperrors.NewInvalidCredentialsError()
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("Account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewInvalidCredentialsError()
	}

	// Generate tokens
	return s.generateAuthResponse(ctx, user)
}

// RefreshToken generates new tokens using refresh token
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Hash the token to find in database
	tokenHash := hashToken(refreshToken)

	// Get refresh token from database
	storedToken, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get refresh token", err)
	}
	if storedToken == nil {
		return nil, apperrors.NewTokenInvalidError()
	}

	// Check if token is revoked
	if storedToken.RevokedAt.Valid {
		return nil, apperrors.NewTokenInvalidError()
	}

	// Check if token is expired
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, apperrors.NewTokenExpiredError()
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil || !user.IsActive {
		return nil, apperrors.NewUnauthorizedError("User not found or inactive")
	}

	// Revoke old refresh token
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, apperrors.NewInternalError("Failed to revoke old token", err)
	}

	// Generate new tokens
	return s.generateAuthResponse(ctx, user)
}

// Logout revokes the refresh token
func (s *service) Logout(ctx context.Context, userID uint64, refreshToken string) error {
	if refreshToken != "" {
		tokenHash := hashToken(refreshToken)
		if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			return apperrors.NewInternalError("Failed to revoke token", err)
		}
	}
	return nil
}

// GetCurrentUser retrieves the current user's information with company data if available
func (s *service) GetCurrentUser(ctx context.Context, userID uint64) (*UserWithCompanyResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("User")
	}

	// If user is a company, get company data
	var companyData *CompanyData
	if user.Role == RoleCompany {
		companyData, err = s.repo.GetCompanyByUserID(ctx, userID)
		if err != nil {
			// Log error but don't fail - return user data without company info
			return user.ToResponseWithCompany(nil), nil
		}
	}

	return user.ToResponseWithCompany(companyData), nil
}

// UpdateProfile updates the current user's profile
func (s *service) UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) (*UserResponse, error) {
	// Get current user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("User")
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	// Note: Company information is now managed through the companies table
	// These fields are kept in UpdateProfileRequest for API compatibility
	// but are not stored in the users table anymore
	
	user.UpdatedAt = time.Now()

	// Update user in database
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("Failed to update user", err)
	}

	return user.ToResponse(), nil
}

// ValidateAccessToken validates and parses an access token
func (s *service) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, apperrors.NewTokenInvalidError()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, apperrors.NewTokenExpiredError()
			}
		}

		// Check token type
		tokenType, _ := claims["token_type"].(string)
		if tokenType != "access" {
			return nil, apperrors.NewTokenInvalidError()
		}

		return &TokenClaims{
			UserID:    uint64(claims["user_id"].(float64)),
			Email:     claims["email"].(string),
			Role:      claims["role"].(string),
			TokenType: tokenType,
		}, nil
	}

	return nil, apperrors.NewTokenInvalidError()
}

// generateAuthResponse generates tokens and auth response
func (s *service) generateAuthResponse(ctx context.Context, user *User) (*AuthResponse, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to generate access token", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to generate refresh token", err)
	}

	// Get company data if user is a company
	var userResp interface{} = user.ToResponse()
	if user.Role == RoleCompany {
		companyData, err := s.repo.GetCompanyByUserID(ctx, user.ID)
		if err == nil && companyData != nil {
			userResp = user.ToResponseWithCompany(companyData)
		}
	}

	return &AuthResponse{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessExpiry.Seconds()),
	}, nil
}

// generateAccessToken generates a new access token
func (s *service) generateAccessToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"token_type": "access",
		"exp":        time.Now().Add(s.config.AccessExpiry).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// generateRefreshToken generates a new refresh token and stores it
func (s *service) generateRefreshToken(ctx context.Context, userID uint64) (string, error) {
	// Generate random token
	tokenValue := uuid.New().String()
	tokenHash := hashToken(tokenValue)

	// Store in database
	refreshToken := &RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.config.RefreshExpiry),
	}

	if err := s.repo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err
	}

	return tokenValue, nil
}

// hashToken creates a SHA256 hash of a token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
