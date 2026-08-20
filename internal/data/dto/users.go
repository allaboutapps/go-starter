package dto

import (
	"time"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/aarondl/null/v8"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
)

type User struct {
	ID                  string
	Username            null.String
	PasswordHash        null.String
	IsActive            bool
	Scopes              UserScopes
	LastAuthenticatedAt null.Time
	UpdatedAt           time.Time
	Profile             *AppUserProfile
}

type UserScopes []string

func (s UserScopes) ToTypes() []types.UserScope {
	result := make([]types.UserScope, len(s))

	for i, scope := range s {
		result[i] = types.UserScope(scope)
	}

	return result
}

func (s UserScopes) ToModels() []string {
	return []string(s)
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
		Sub:       conv.Pointer(u.ID),
		UpdatedAt: conv.Pointer(u.LastUpdatedAt().Unix()),
		Email:     strfmt.Email(u.Username.String),
		Scopes:    u.Scopes.ToTypes(),
	}
}

func (u User) ToModels() *models.User {
	return &models.User{
		ID:                  u.ID,
		Username:            u.Username,
		Password:            u.PasswordHash,
		IsActive:            u.IsActive,
		Scopes:              u.Scopes.ToModels(),
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

type RegisterResult struct {
	RequiresConfirmation bool
	ConfirmationToken    null.String
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
		AccessToken:  conv.Pointer(strfmt.UUID4(l.AccessToken)),
		RefreshToken: conv.Pointer(strfmt.UUID4(l.RefreshToken)),
		ExpiresIn:    conv.Pointer(l.ExpiresIn),
		TokenType:    conv.Pointer(l.TokenType),
	}
}

type ResetPasswordRequest struct {
	ResetToken  string
	NewPassword string
}

type Username struct {
	val string
}

func NewUsername(val string) Username {
	return Username{val: val}
}

func (u Username) String() string {
	return util.ToUsernameFormat(u.val)
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

type CompleteRegisterRequest struct {
	ConfirmationToken string
}

type DeleteUserAccountRequest struct {
	User            User
	CurrentPassword string
}

type ConfirmatioNotificationPayload struct {
	ConfirmationLink string
}
