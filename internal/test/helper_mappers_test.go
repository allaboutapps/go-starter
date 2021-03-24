package test_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestGetMapFromStruct(t *testing.T) {
	type tmp struct {
		A string
		B int
		C interface{}
		D float32
		E bool
		F *string
		G *int
	}

	f := "stringPtr"
	g := 5
	x := tmp{
		A: "string",
		B: 1,
		C: tmp{
			A: "string2",
		},
		D: 2.3,
		E: true,
		F: &f,
		G: &g,
	}

	xMap := test.GetMapFromStruct(x)
	assert.Len(t, xMap, 7)
	assert.Equal(t, "string", xMap["A"])
	assert.Equal(t, "1", xMap["B"])
	assert.Contains(t, xMap["C"], "string2")
	assert.Equal(t, "2.3", xMap["D"])
	assert.Equal(t, "true", xMap["E"])
	assert.Equal(t, "stringPtr", xMap["F"])
	assert.Equal(t, "5", xMap["G"])
}

func TestGetMapFromStructByTag(t *testing.T) {
	type tmp struct {
		A string      `x:"1,omitempty" y:"2"`
		B int         `x:"3"`
		C interface{} `x:"12"`
		D float32     `x:"2"`
		E bool        `x:"5"`
		F *string     `x:"4"`
		G *int        `x:"6"`
	}

	f := "stringPtr"
	g := 5
	x := tmp{
		A: "string",
		B: 1,
		C: tmp{
			A: "string2",
		},
		D: 2.3,
		E: true,
		F: &f,
		G: &g,
	}

	xMap := test.GetMapFromStructByTag("x", x)
	assert.Len(t, xMap, 7)
	assert.Equal(t, "string", xMap["1"])
	assert.Equal(t, "1", xMap["3"])
	assert.Contains(t, xMap["12"], "string2")
	assert.Equal(t, "2.3", xMap["2"])
	assert.Equal(t, "true", xMap["5"])
	assert.Equal(t, "stringPtr", xMap["4"])
	assert.Equal(t, "5", xMap["6"])

	yMap := test.GetMapFromStructByTag("y", x)
	assert.Len(t, yMap, 1)
	assert.Equal(t, "string", yMap["2"])
}
