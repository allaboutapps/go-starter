package dto

import (
	"time"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/go-openapi/swag"
	"github.com/volatiletech/null/v8"
)

type User struct {
	ID                  string
	Username            null.String
	PasswordHash        null.String
	IsActive            bool
	Scopes              []string
	LastAuthenticatedAt null.Time
	UpdatedAt           time.Time
	Profile             *AppUserProfile
}

func (u User) LastUpdatedAt() time.Time {
	if u.Profile != nil && u.Profile.UpdatedAt.After(u.UpdatedAt) {
		return u.Profile.UpdatedAt
	}

	return u.UpdatedAt
}

func (u User) Ptr() *User {
	return &u
}

func (u User) ToTypes() *types.GetUserInfoResponse {
	return &types.GetUserInfoResponse{
		Sub:       swag.String(u.ID),
		UpdatedAt: swag.Int64(u.LastUpdatedAt().Unix()),
		Email:     strfmt.Email(u.Username.String),
		Scopes:    u.Scopes,
	}
}

func (u User) ToModels() *models.User {
	return &models.User{
		ID:                  u.ID,
		Username:            u.Username,
		Password:            u.PasswordHash,
		IsActive:            u.IsActive,
		Scopes:              u.Scopes,
		LastAuthenticatedAt: u.LastAuthenticatedAt,
	}
}

type AppUserProfile struct {
	UserID          string
	LegalAcceptedAt null.Time
	UpdatedAt       time.Time
}

func (aup AppUserProfile) Ptr() *AppUserProfile {
	return &aup
}

type UpdatePasswordRequest struct {
	User                            User
	CurrentPassword                 string
	SkipCurrentPasswordVerification bool
	NewPassword                     string
}

type LoginResult struct {
	UserID       string
	AccessToken  string
	ExpiresIn    int64
	RefreshToken string
	TokenType    string
}

func (l LoginResult) ToTypes() *types.PostLoginResponse {
	return &types.PostLoginResponse{
		AccessToken:  conv.UUID4(strfmt.UUID4(l.AccessToken)),
		RefreshToken: conv.UUID4(strfmt.UUID4(l.RefreshToken)),
		ExpiresIn:    swag.Int64(l.ExpiresIn),
		TokenType:    swag.String(l.TokenType),
	}
}

type ResetPasswordRequest struct {
	ResetToken  string
	NewPassword string
}

func NewUsername(val string) Username {
	return Username{val: val}
}

type Username struct {
	val string
}

func (u Username) String() string {
	return util.ToUsernameFormat(string(u.val))
}

type InitPasswordResetRequest struct {
	Username Username
}

type InitPasswordResetResult struct {
	ResetToken null.String
}

type LoginRequest struct {
	Username Username
	Password string
}

type LogoutRequest struct {
	AccessToken  string
	RefreshToken null.String
}

type AuthenticateUserRequest struct {
	User                     User
	InvalidateExistingTokens bool
}

type RefreshRequest struct {
	RefreshToken string
}

type RegisterRequest struct {
	Username Username
	Password string
}
