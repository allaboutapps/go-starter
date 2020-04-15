package swaggertest

import "github.com/go-openapi/strfmt"

// SomeSampleType SomeSampleType some sample type
//
// swagger:model SomeSampleType
type SomeSampleType struct {

	// Min Length: 5
	// Max Length: 10
	Data string `json:"data"`

	ID strfmt.UUID `json:"id"`

	// Min: 1
	// Max: 100
	Num int `json:"num"`

	// Default: false
	// Required: true
	IsActive bool `json:"isActive"`

	Mail strfmt.Email `json:"mail"`

	MoreData string `json:"moreData"`
}
