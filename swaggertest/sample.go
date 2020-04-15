package swaggertest

// SomeSampleType SomeSampleType some sample type
//
// swagger:model SomeSampleType
type SomeSampleType struct {

	// data
	// Required: true
	// Min Length: 5
	Data *string `json:"data"`

	// Id
	// Required: true
	ID *string `json:"id"`

	// moreData
	// Required: true
	MoreData *string `json:"moreData"`
}
