// Package classification awesome.
//
// Documentation of our awesome API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: some-url.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package types

// swagger:parameters initializeTemplateHandler
type InitializeTemplatePayload struct {
	// Hash used in initializing a template
	//
	// min items: 1
	// max items: 1
	// unique: true
	// in: query
	Hash string `json:"hash"`
}
