package dto

import "github.com/volatiletech/null/v8"

type UpdatePushTokenRequest struct {
	User          User
	Token         string
	Provider      string
	ExistingToken null.String
}
