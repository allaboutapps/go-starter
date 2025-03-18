package local

import (
	"context"
	"database/sql"
	"errors"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *Service) UpdatePushToken(ctx context.Context, request dto.UpdatePushTokenRequest) error {
	log := util.LogFromContext(ctx).With().Str("userID", request.User.ID).Logger()

	err := db.WithTransaction(ctx, s.db, func(exec boil.ContextExecutor) error {
		tokenExists, err := models.PushTokens(
			models.PushTokenWhere.Token.EQ(request.Token),
		).Exists(ctx, exec)
		if err != nil {
			log.Err(err).Msg("Failed to check if token exists")
			return err
		}

		if tokenExists {
			log.Debug().Msg("Token already exists")
			return httperrors.ErrConflictPushToken
		}

		newToken := models.PushToken{
			UserID:   request.User.ID,
			Token:    request.Token,
			Provider: request.Provider,
		}

		if err := newToken.Insert(ctx, s.db, boil.Infer()); err != nil {
			log.Err(err).Msg("Failed to insert new token")
			return err
		}

		if request.ExistingToken.IsZero() {
			return nil
		}

		existingToken, err := models.PushTokens(
			models.PushTokenWhere.Token.EQ(request.ExistingToken.String),
			models.PushTokenWhere.UserID.EQ(request.User.ID),
		).One(ctx, exec)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Msg("Existing token not found")
				return httperrors.ErrNotFoundOldPushToken
			}

			log.Err(err).Msg("Failed to find existing token")
			return err
		}

		if _, err := existingToken.Delete(ctx, exec); err != nil {
			log.Err(err).Msg("Failed to delete existing token")
			return err
		}

		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Failed to update push token")
		return err
	}

	return nil
}
