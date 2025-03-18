package mapper

import (
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/models"
)

func LocalAppUserProfileToDTOPtr(appUserProfile *models.AppUserProfile) dto.AppUserProfile {
	return dto.AppUserProfile{
		UserID:          appUserProfile.UserID,
		LegalAcceptedAt: appUserProfile.LegalAcceptedAt,
		UpdatedAt:       appUserProfile.UpdatedAt,
	}
}

func LocalUserToDTO(user *models.User) dto.User {
	return dto.User{
		ID:                  user.ID,
		Username:            user.Username,
		IsActive:            user.IsActive,
		Scopes:              user.Scopes,
		LastAuthenticatedAt: user.LastAuthenticatedAt,
		UpdatedAt:           user.UpdatedAt,
		PasswordHash:        user.Password,
	}
}
