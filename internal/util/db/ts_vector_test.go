package db_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/go-openapi/swag/conv"
	"github.com/stretchr/testify/assert"
)

func TestSearchStringToTSQuery(t *testing.T) {
	expected := "'abcde':* & '12345':* & 'xyz':*"
	search := conv.Pointer("    abcde 12345 xyz   ")
	out := db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = "'abcde':*"
	search = conv.Pointer("abcde")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = "'Hello':* & 'world':* & 'lorem':* & '12345':* & 'ipsum':* & 'abc':* & 'def':*"
	search = conv.Pointer("    Hello  world lorem 12345               ipsum  abc def  ")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = nil
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = conv.Pointer("")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)

	expected = ""
	search = conv.Pointer("                ''   '    '    '      ")
	out = db.SearchStringToTSQuery(search)
	assert.Equal(t, expected, out)
}
