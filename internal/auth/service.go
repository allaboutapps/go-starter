package auth

import (
	"context"
	"database/sql"
	"errors"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/data/mapper"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"allaboutapps.dev/aw/go-starter/internal/util/hashing"
	"github.com/dropbox/godropbox/time2"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Service struct {
	config config.Server
	db     *sql.DB
	clock  time2.Clock
}

func NewService(config config.Server, db *sql.DB, clock time2.Clock) *Service {
	return &Service{
		config: config,
		db:     db,
		clock:  clock,
	}
}

func (s *Service) GetAppUserProfile(ctx context.Context, userID string) (*dto.AppUserProfile, error) {
	log := util.LogFromContext(ctx).With().Str("userID", userID).Logger()

	aup, err := models.AppUserProfiles(
		models.AppUserProfileWhere.UserID.EQ(userID),
	).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("AppUserProfile not found")
			return nil, ErrNotFound
		}

		log.Err(err).Msg("Failed to get AppUserProfile")
		return nil, err
	}

	return mapper.LocalAppUserProfileToDTO(aup).Ptr(), nil
}

func (s *Service) UpdatePassword(ctx context.Context, request dto.UpdatePasswordRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx).With().Str("userID", request.User.ID).Logger()

	if !request.User.IsActive {
		log.Debug().Msg("User is deactivated, rejecting password change")
		return dto.LoginResult{}, httperrors.ErrForbiddenUserDeactivated
	}

	if !request.User.PasswordHash.Valid {
		log.Debug().Msg("Failed to update user password, user is missing password")
		return dto.LoginResult{}, httperrors.ErrForbiddenNotLocalUser
	}

	if !request.SkipCurrentPasswordVerification {
		match, err := hashing.ComparePasswordAndHash(request.CurrentPassword, request.User.PasswordHash.String)
		if err != nil {
			log.Err(err).Msg("Failed to compare password with stored hash")
			return dto.LoginResult{}, err
		}

		if !match {
			log.Debug().Msg("Failed to update user password, provided password does not match stored hash")
			return dto.LoginResult{}, echo.ErrUnauthorized
		}
	}

	hash, err := hashing.HashPassword(request.NewPassword, hashing.DefaultArgon2Params)
	if err != nil {
		log.Err(err).Msg("Failed to hash new password")
		return dto.LoginResult{}, httperrors.ErrBadRequestInvalidPassword
	}

	var result dto.LoginResult
	if err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		request.User.PasswordHash = null.StringFrom(hash)

		user := request.User.ToModels()

		if _, err := user.Update(ctx, exec, boil.Whitelist(models.UserColumns.Password, models.UserColumns.UpdatedAt)); err != nil {
			log.Err(err).Msg("Failed to update user")
			return err
		}

		result, err = s.authenticateUser(ctx, exec, dto.AuthenticateUserRequest{
			User:                     request.User,
			InvalidateExistingTokens: true,
		})
		if err != nil {
			log.Err(err).Msg("Failed to authenticate user after password change")
			return err
		}

		return nil
	}); err != nil {
		log.Debug().Err(err).Msg("Failed to change password")
		return dto.LoginResult{}, err
	}

	return result, nil
}

func (s *Service) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx).With().Logger()

	passwordResetToken, err := models.PasswordResetTokens(
		models.PasswordResetTokenWhere.Token.EQ(request.ResetToken),
		qm.Load(models.PasswordResetTokenRels.User),
	).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("Password reset token not found")
			return dto.LoginResult{}, httperrors.ErrNotFoundTokenNotFound
		}

		log.Err(err).Msg("Failed to load password reset token")
		return dto.LoginResult{}, err
	}

	if s.clock.Now().After(passwordResetToken.ValidUntil) {
		log.Debug().Time("validUntil", passwordResetToken.ValidUntil).Msg("Password reset token is no longer valid, rejecting password reset")
		return dto.LoginResult{}, httperrors.ErrConflictTokenExpired
	}

	return s.UpdatePassword(ctx, dto.UpdatePasswordRequest{
		User:                            mapper.LocalUserToDTO(passwordResetToken.R.User),
		NewPassword:                     request.NewPassword,
		SkipCurrentPasswordVerification: true,
	})
}

func (s *Service) InitPasswordReset(ctx context.Context, request dto.InitPasswordResetRequest) (dto.InitPasswordResetResult, error) {
	log := util.LogFromContext(ctx).With().Str("username", request.Username.String()).Logger()

	user, err := models.Users(
		models.UserWhere.Username.EQ(null.StringFrom(request.Username.String())),
	).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("User not found")
			return dto.InitPasswordResetResult{}, nil
		}

		log.Err(err).Err(err).Msg("Failed to load user")
		return dto.InitPasswordResetResult{}, err
	}

	if !user.IsActive {
		log.Debug().Msg("User is deactivated, skipping password reset")
		return dto.InitPasswordResetResult{}, nil
	}

	if !user.Password.Valid {
		log.Debug().Msg("User is missing password, skipping password reset")
		return dto.InitPasswordResetResult{}, nil
	}

	if s.config.Auth.PasswordResetTokenDebounceDuration > 0 {
		resetTokenInDebounceTimeExists, err := user.PasswordResetTokens(
			models.PasswordResetTokenWhere.CreatedAt.GT(s.clock.Now().Add(-s.config.Auth.PasswordResetTokenDebounceDuration)),
			models.PasswordResetTokenWhere.ValidUntil.GT(s.clock.Now()),
		).Exists(ctx, s.db)
		if err != nil {
			log.Err(err).Msg("Failed to check for existing password reset token")
			return dto.InitPasswordResetResult{}, err
		}

		if resetTokenInDebounceTimeExists {
			log.Debug().Msg("Password reset token exists within debounce time, not sending new one")
			return dto.InitPasswordResetResult{}, nil
		}
	}

	var result dto.InitPasswordResetResult
	if err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		passwordResetToken, err := user.PasswordResetTokens(
			models.PasswordResetTokenWhere.CreatedAt.GT(s.clock.Now().Add(-s.config.Auth.PasswordResetTokenReuseDuration)),
			models.PasswordResetTokenWhere.ValidUntil.GT(s.clock.Now()),
		).One(ctx, exec)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Err(err).Msg("No valid password reset token exists, creating new one")

				passwordResetToken = &models.PasswordResetToken{
					UserID:     user.ID,
					ValidUntil: s.clock.Now().Add(s.config.Auth.PasswordResetTokenValidity),
				}

				if err := passwordResetToken.Insert(ctx, exec, boil.Infer()); err != nil {
					log.Err(err).Msg("Failed to insert password reset token")
					return err
				}
			} else {
				log.Err(err).Msg("Failed to check for existing valid password reset token")
				return err
			}
		}

		result.ResetToken = null.StringFrom(passwordResetToken.Token)

		return nil
	}); err != nil {
		log.Debug().Err(err).Msg("Failed to initiate password reset")
		return dto.InitPasswordResetResult{}, err
	}

	return result, nil
}

func (s *Service) Logout(ctx context.Context, request dto.LogoutRequest) error {
	log := util.LogFromContext(ctx)

	if err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		if _, err := models.AccessTokens(models.AccessTokenWhere.Token.EQ(request.AccessToken)).DeleteAll(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete access token")
			return err
		}

		if request.RefreshToken.IsZero() {
			return nil
		}

		if _, err := models.RefreshTokens(models.RefreshTokenWhere.Token.EQ(request.RefreshToken.String)).DeleteAll(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete refresh token")
			return err
		}

		return nil
	}); err != nil {
		log.Debug().Err(err).Msg("Failed to process logout")
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx)

	user, err := models.Users(
		models.UserWhere.Username.EQ(null.StringFrom(request.Username.String())),
	).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("User not found")
		}

		log.Err(err).Msg("Failed to load user")

		return dto.LoginResult{}, echo.ErrUnauthorized
	}

	if !user.IsActive {
		log.Debug().Msg("User is deactivated, rejecting authentication")
		return dto.LoginResult{}, httperrors.ErrForbiddenUserDeactivated
	}

	if !user.Password.Valid {
		log.Debug().Msg("User is missing password, forbidding authentication")
		return dto.LoginResult{}, echo.ErrUnauthorized
	}

	match, err := hashing.ComparePasswordAndHash(request.Password, user.Password.String)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to compare password with stored hash")
		return dto.LoginResult{}, echo.ErrUnauthorized
	}

	if !match {
		log.Debug().Msg("Provided password does not match stored hash")
		return dto.LoginResult{}, echo.ErrUnauthorized
	}

	var result dto.LoginResult
	err = db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		var err error
		result, err = s.authenticateUser(ctx, exec, dto.AuthenticateUserRequest{
			User: mapper.LocalUserToDTO(user),
		})
		if err != nil {
			log.Err(err).Msg("Failed to authenticate user")
			return err
		}

		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to authenticate user")
		return dto.LoginResult{}, err
	}

	return result, nil
}

func (s *Service) Refresh(ctx context.Context, request dto.RefreshRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx)

	oldRefreshToken, err := models.RefreshTokens(
		models.RefreshTokenWhere.Token.EQ(request.RefreshToken),
		qm.Load(models.RefreshTokenRels.User),
	).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Err(err).Msg("Refresh token not found")
			return dto.LoginResult{}, echo.ErrUnauthorized
		}

		log.Err(err).Msg("Failed to load refresh token")
		return dto.LoginResult{}, err
	}

	user := oldRefreshToken.R.User

	if !user.IsActive {
		log.Debug().Msg("User is deactivated, rejecting token refresh")
		return dto.LoginResult{}, httperrors.ErrForbiddenUserDeactivated
	}

	var result dto.LoginResult
	err = db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		_, err = oldRefreshToken.Delete(ctx, exec)
		if err != nil {
			log.Err(err).Msg("Failed to delete old refresh token")
			return err
		}

		result, err = s.authenticateUser(ctx, exec, dto.AuthenticateUserRequest{
			User:                     mapper.LocalUserToDTO(user),
			InvalidateExistingTokens: false,
		})
		if err != nil {
			log.Err(err).Msg("Failed to authenticate user")
			return err
		}

		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to refresh token")
		return dto.LoginResult{}, err
	}

	return result, nil
}

func (s *Service) Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResult, error) {
	log := util.LogFromContext(ctx).With().Str("username", request.Username.String()).Logger()

	user, err := models.Users(
		models.UserWhere.Username.EQ(null.StringFrom(request.Username.String())),
	).One(ctx, s.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Err(err).Msg("Failed to check whether user exists")
		return dto.RegisterResult{}, err
	}

	if user != nil {
		if !user.RequiresConfirmation {
			log.Debug().Msg("User with given username already exists")
			return dto.RegisterResult{}, httperrors.ErrConflictUserAlreadyExists
		}

		confirmationTokenInDebounceTimeExists, err := user.ConfirmationTokens(
			models.ConfirmationTokenWhere.CreatedAt.GT(s.clock.Now().Add(-s.config.Auth.ConfirmationTokenDebounceDuration)),
			models.ConfirmationTokenWhere.ValidUntil.GT(s.clock.Now()),
		).Exists(ctx, s.db)
		if err != nil {
			log.Err(err).Msg("Failed to check for existing confirmation token")
			return dto.RegisterResult{}, err
		}

		if confirmationTokenInDebounceTimeExists {
			return dto.RegisterResult{
				RequiresConfirmation: user.RequiresConfirmation,
			}, nil
		}

		confirmationToken := models.ConfirmationToken{
			UserID:     user.ID,
			ValidUntil: s.clock.Now().Add(s.config.Auth.ConfirmationTokenValidity),
		}

		if err := confirmationToken.Insert(ctx, s.db, boil.Infer()); err != nil {
			log.Err(err).Msg("Failed to insert confirmation token")
			return dto.RegisterResult{}, err
		}

		return dto.RegisterResult{
			RequiresConfirmation: user.RequiresConfirmation,
			ConfirmationToken:    null.StringFrom(confirmationToken.Token),
		}, nil
	}

	hash, err := hashing.HashPassword(request.Password, hashing.DefaultArgon2Params)
	if err != nil {
		log.Err(err).Msg("Failed to hash user password")
		return dto.RegisterResult{}, httperrors.ErrBadRequestInvalidPassword
	}

	result := dto.RegisterResult{
		RequiresConfirmation: s.config.Auth.RegistrationRequiresConfirmation,
	}

	if err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		user := &models.User{
			Username:             null.StringFrom(request.Username.String()),
			Password:             null.StringFrom(hash),
			LastAuthenticatedAt:  null.TimeFrom(s.clock.Now()),
			IsActive:             !result.RequiresConfirmation,
			RequiresConfirmation: result.RequiresConfirmation,
			Scopes:               s.config.Auth.DefaultUserScopes,
		}

		if err := user.Insert(ctx, exec, boil.Infer()); err != nil {
			log.Err(err).Msg("Failed to insert user")
			return err
		}

		appUserProfile := models.AppUserProfile{
			UserID: user.ID,
		}

		if err := appUserProfile.Insert(ctx, exec, boil.Infer()); err != nil {
			log.Err(err).Msg("Failed to insert app user profile")
			return err
		}

		if result.RequiresConfirmation {
			confirmationToken := models.ConfirmationToken{
				UserID:     user.ID,
				ValidUntil: s.clock.Now().Add(s.config.Auth.ConfirmationTokenValidity),
			}

			if err := confirmationToken.Insert(ctx, exec, boil.Infer()); err != nil {
				log.Err(err).Msg("Failed to insert confirmation token")
				return err
			}

			result.ConfirmationToken = null.StringFrom(confirmationToken.Token)
		}

		return nil
	}); err != nil {
		log.Debug().Err(err).Msg("Failed to register user")
		return dto.RegisterResult{}, err
	}

	return result, nil
}

func (s *Service) DeleteUserAccount(ctx context.Context, request dto.DeleteUserAccountRequest) error {
	log := util.LogFromContext(ctx)

	if !request.User.IsActive {
		log.Debug().Msg("User is deactivated, rejecting deletion")
		return httperrors.ErrForbiddenUserDeactivated
	}

	if !request.User.PasswordHash.Valid {
		log.Debug().Msg("Failed to delete user account, user is missing password")
		return httperrors.ErrForbiddenNotLocalUser
	}

	match, err := hashing.ComparePasswordAndHash(request.CurrentPassword, request.User.PasswordHash.String)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to compare password with stored hash")
		return echo.ErrUnauthorized
	}

	if !match {
		log.Debug().Msg("Provided password does not match stored hash")
		return echo.ErrUnauthorized
	}

	// delete the user and all related data
	err = db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		_, err = models.Users(
			models.UserWhere.ID.EQ(request.User.ID),
		).DeleteAll(ctx, exec)
		if err != nil {
			log.Err(err).Msg("Failed to delete user")
			return err
		}

		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to delete user account")
		return err
	}

	return nil
}

func (s *Service) CompleteRegister(ctx context.Context, request dto.CompleteRegisterRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx)

	var result dto.LoginResult
	err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		confirmationToken, err := models.ConfirmationTokens(
			models.ConfirmationTokenWhere.Token.EQ(request.ConfirmationToken),
			models.ConfirmationTokenWhere.ValidUntil.GT(s.clock.Now()),
			qm.Load(models.ConfirmationTokenRels.User),
		).One(ctx, s.db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Err(err).Msg("Confirmation token not found")
				return httperrors.ErrNotFoundTokenNotFound
			}

			log.Err(err).Msg("Failed to load confirmation token")
			return err
		}

		user := confirmationToken.R.User
		if user.IsActive || !user.RequiresConfirmation {
			log.Debug().Msg("User already active, skipping confirmation")
			return nil
		}

		user.IsActive = true
		user.RequiresConfirmation = false
		if _, err := user.Update(ctx, exec, boil.Whitelist(models.UserColumns.IsActive, models.UserColumns.RequiresConfirmation, models.UserColumns.UpdatedAt)); err != nil {
			log.Err(err).Msg("Failed to update user")
			return err
		}

		if _, err := confirmationToken.Delete(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete confirmation token")
			return err
		}

		result, err = s.authenticateUser(ctx, exec, dto.AuthenticateUserRequest{
			User: mapper.LocalUserToDTO(confirmationToken.R.User),
		})
		if err != nil {
			log.Err(err).Msg("Failed to authenticate user")
			return err
		}

		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to complete registration")
		return dto.LoginResult{}, err
	}

	return result, nil
}

func (s *Service) authenticateUser(ctx context.Context, exec boil.ContextExecutor, request dto.AuthenticateUserRequest) (dto.LoginResult, error) {
	log := util.LogFromContext(ctx)

	result := dto.LoginResult{
		TokenType: TokenTypeBearer,
		ExpiresIn: int64(s.config.Auth.AccessTokenValidity.Seconds()),
	}

	if request.InvalidateExistingTokens {
		if _, err := models.AccessTokens(
			models.AccessTokenWhere.UserID.EQ(request.User.ID),
		).DeleteAll(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete existing access tokens")
			return dto.LoginResult{}, err
		}

		if _, err := models.RefreshTokens(
			models.RefreshTokenWhere.UserID.EQ(request.User.ID),
		).DeleteAll(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete existing refresh tokens")
			return dto.LoginResult{}, err
		}
	}

	accessToken := models.AccessToken{
		ValidUntil: s.clock.Now().Add(s.config.Auth.AccessTokenValidity),
		UserID:     request.User.ID,
	}

	if err := accessToken.Insert(ctx, exec, boil.Infer()); err != nil {
		log.Err(err).Msg("Failed to insert access token")
		return dto.LoginResult{}, err
	}

	refreshToken := models.RefreshToken{
		UserID: request.User.ID,
	}

	if err := refreshToken.Insert(ctx, exec, boil.Infer()); err != nil {
		log.Err(err).Msg("Failed to insert refresh token")
		return dto.LoginResult{}, err
	}

	u := request.User.ToModels()
	u.LastAuthenticatedAt = null.TimeFrom(s.clock.Now())

	if _, err := u.Update(ctx, exec, boil.Whitelist(models.UserColumns.LastAuthenticatedAt, models.UserColumns.UpdatedAt)); err != nil {
		log.Err(err).Msg("Failed to update user last authenticated time")
		return dto.LoginResult{}, err
	}

	result.AccessToken = accessToken.Token
	result.RefreshToken = refreshToken.Token

	return result, nil
}
