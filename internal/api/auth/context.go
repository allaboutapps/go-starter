package auth

import (
	"context"

	"allaboutapps.at/aw/go-mranftl-sample/internal/models"
	"allaboutapps.at/aw/go-mranftl-sample/internal/util"
	"github.com/labstack/echo/v4"
)

// EnrichContextWithCredentials stores the provided credentials in the form of user and access token used for authentication
// in the give context and updates the logger associated with ctx to include the user's ID.
func EnrichContextWithCredentials(ctx context.Context, user *models.User, accessToken *models.AccessToken) context.Context {
	// Retrieve current logger associated with context and extend it ID of authenticated user
	l := util.LogFromContext(ctx).With().Str("user_id", user.ID).Logger()
	c := l.WithContext(ctx)

	// Store authenticated user's instance in context
	c = context.WithValue(c, util.CTXKeyUser, user)
	// Store access token used for authentication in context
	c = context.WithValue(c, util.CTXKeyAccessToken, accessToken)

	return c
}

// EnrichEchoContextWithCredentials stores the provided credentials in the form of user and access token user for authentication
// in the given echo context's request and updates the logger associated with c to include the user's ID.
func EnrichEchoContextWithCredentials(c echo.Context, user *models.User, accessToken *models.AccessToken) echo.Context {
	// Get current context and enrich it with credentials
	req := c.Request()
	ctx := EnrichContextWithCredentials(req.Context(), user, accessToken)

	// Set updated request with enriched context in echo context
	c.SetRequest(req.WithContext(ctx))

	return c
}

// UserFromContext returns the user model of the currently authenticated user from a context. If no authentication was provided
// or the current context does not carry any user information, nil will be returned instead.
func UserFromContext(ctx context.Context) *models.User {
	u := ctx.Value(util.CTXKeyUser)
	if u == nil {
		return nil
	}

	user, ok := u.(*models.User)
	if !ok {
		return nil
	}

	return user
}

// UserFromEchoContext returns the user model of the currently authenticated user from an echo context. If no authentication was
// provided or the current echo context does not carry any user information, nil will be returned instead.
func UserFromEchoContext(c echo.Context) *models.User {
	return UserFromContext(c.Request().Context())
}

// AccessTokenFromContext returns the access token model of the token used to authentication from a context. If no authentication was
// provided or the current context does not carry any access token information, nil will be returned instead.
func AccessTokenFromContext(ctx context.Context) *models.AccessToken {
	t := ctx.Value(util.CTXKeyAccessToken)
	if t == nil {
		return nil
	}

	token, ok := t.(*models.AccessToken)
	if !ok {
		return nil
	}

	return token
}

// AccessTokenFromEchoContext returns the access token model of the token used to authentication from an echo context. If no authentication
// was provided or the current context does not carry any access token information, nil will be returned instead.
func AccessTokenFromEchoContext(c echo.Context) *models.AccessToken {
	return AccessTokenFromContext(c.Request().Context())
}
