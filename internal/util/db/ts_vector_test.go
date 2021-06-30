package db_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestSearchStringToTSQuery(t *testing.T) {
	expected := "'abcde':* & '12345':* & 'xyz':*"
	in := swag.String("    abcde 12345 xyz   ")
	out := db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = "'abcde':*"
	in = swag.String("abcde")
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = "'Hello':* & 'world':* & 'lorem':* & '12345':* & 'ipsum':* & 'abc':* & 'def':*"
	in = swag.String("    Hello  world lorem 12345               ipsum  abc def  ")
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = ""
	in = nil
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = ""
	in = swag.String("")
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = ""
	in = swag.String("                ''   '    '    '      ")
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)
}
