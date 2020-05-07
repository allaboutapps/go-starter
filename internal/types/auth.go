package types

import "github.com/go-openapi/strfmt"

// basic swagger model properties: https://goswagger.io/use/spec/model.html#properties
// special types: https://goswagger.io/use/spec/strfmt.html

// swagger:parameters PostLoginRoute
type PostLoginPayloadParam struct {
	// in:body
	Payload PostLoginPayload
}

// swagger:model
type PostLoginPayload struct {
	// Username of user to authenticate as
	// required: true
	// min length: 1
	// max length: 255
	// example: user@example.com
	Username *strfmt.Email `json:"username"`
	// Password of user to authenticate as
	// required: true
	// min length: 1
	// example: correct horse battery staple
	Password *string `json:"password"`
}

// swagger:model
type PostLoginResponse struct {
	// Access token required for accessing protected API endpoints
	// required: true
	// example: c1247d8d-0d65-41c4-bc86-ec041d2ac437
	AccessToken strfmt.UUID4 `json:"access_token"`
	// Type of access token, will always be `bearer`
	// required: true
	// example: bearer
	TokenType string `json:"token_type"`
	// Access token expiry in seconds
	// required: true
	// example: 86400
	ExpiresIn int `json:"expires_in"`
	// Refresh token for refreshing the access token once it expires
	// required: true
	// example: 1dadb3bd-50d8-485d-83a3-6111392568f0
	RefreshToken strfmt.UUID4 `json:"refresh_token"`
}

// swagger:parameters PostRefreshRoute
type PostRefreshPayloadParam struct {
	// in:body
	Payload PostRefreshPayload
}

// swagger:model
type PostRefreshPayload struct {
	// Refresh token to use for retrieving new token set
	// required: true
	// example: 7503cd8a-c921-4368-a32d-6c1d01d86da9
	RefreshToken strfmt.UUID4 `json:"refresh_token"`
}

// swagger:parameters PostLogoutRoute
type PostLogoutPayloadParam struct {
	// in:body
	Payload PostLogoutPayload
}

// swagger:model
type PostLogoutPayload struct {
	// Optional refresh token to delete while logging out
	// example: 700ebed3-40f7-4211-bc83-a89b22b9875e
	RefreshToken strfmt.UUID4 `json:"refresh_token"`
}

// swagger:parameters PostRegisterRoute
type PostRegisterPayloadParam struct {
	// in:body
	Payload PostRegisterPayload
}

// swagger:model
type PostRegisterPayload struct {
	// Username to register with
	// required: true
	// min length: 1
	// max length: 255
	// example: user@example.com
	Username *strfmt.Email `json:"username"`
	// Password to register with
	// required: true
	// min length: 1
	// example: correct horse battery staple
	Password *string `json:"password"`
}

// swagger:parameters PostChangePasswordRoute
type PostChangePasswordPayloadParam struct {
	// in:body
	Payload PostChangePasswordPayload
}

// swagger:model
type PostChangePasswordPayload struct {
	// Current password of user
	// required: true
	// min length: 1
	// example: correct horse battery staple
	CurrentPassword *string `json:"currentPassword"`
	// New password to set for user
	// required: true
	// min length: 1
	// example: correct horse battery staple
	NewPassword *string `json:"newPassword"`
}
