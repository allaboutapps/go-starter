package db_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestSearchStringToTSQuery(t *testing.T) {
	expected := "'abcde':* & '12345':* & 'xyz':*"
	search := swag.String("    abcde 12345 xyz   ")
	out := db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = "'abcde':*"
	search = swag.String("abcde")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = "'Hello':* & 'world':* & 'lorem':* & '12345':* & 'ipsum':* & 'abc':* & 'def':*"
	search = swag.String("    Hello  world lorem 12345               ipsum  abc def  ")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = nil
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = swag.String("")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = swag.String("                ''   '    '    '      ")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)
}
