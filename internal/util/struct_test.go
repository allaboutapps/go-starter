package util_test

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readInterface interface {
	Read(p []byte) (n int, err error)
}

type writeInterface interface {
	WriteTo(w io.Writer) (n int64, err error)
}

type testStruct struct {
	// satisfy only readInterface
	LimitedReader *io.LimitedReader
	Reader        io.Reader

	// satisfy both readInterface and writeInterface
	Buffer1   *bytes.Buffer
	Buffer2   *bytes.Buffer
	NetBuffer *net.Buffers
}

func TestGetFieldsImplementingInvalidInput(t *testing.T) {
	// Invalid interfaceObject input param, must be a pointer to an interface
	// Pointer to a struct
	_, err := util.GetFieldsImplementing(&testStructEmpty{}, &testStructEmpty{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")
	// Pointer to a pointer to an interface
	interfaceObjPtr := (*readInterface)(nil)
	_, err = util.GetFieldsImplementing(&testStructEmpty{}, &interfaceObjPtr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")

	// Invalid structPtr input param, must be a pointer to a struct
	_, err = util.GetFieldsImplementing(testStructEmpty{}, (*readInterface)(nil))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing((*readInterface)(nil), (*readInterface)(nil))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing([]*testStructEmpty{}, (*readInterface)(nil))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
}

func TestGetFieldsImplementingNoFields(t *testing.T) {
	// No fields returned from empty structs
	structEmptyFields, err := util.GetFieldsImplementing(&testStructEmpty{}, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Empty(t, structEmptyFields)

	// No fields returned from structs with only private fields
	structPrivate := testStructPrivateFiled{privateMember: bytes.NewBufferString("my content")}
	structPrivateFields, err := util.GetFieldsImplementing(&structPrivate, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Empty(t, structPrivateFields)

	// No fields returned if struct fields are primitive
	structPrimitive := testStructPrimitives{X: 12, Y: "y", XPtr: swag.Int(15), YPtr: swag.String("YPtr")}
	structPrimitiveFields, err := util.GetFieldsImplementing(&structPrimitive, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Empty(t, structPrimitiveFields)

	// No fields returned if struct fields are structs (not pointer to a struct)
	structMemberStruct := testStructMemberStruct{Member: *bytes.NewBufferString("my content")}
	structMemberStructFields, err := util.GetFieldsImplementing(&structMemberStruct, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Empty(t, structMemberStructFields)

	// No fieds returned if an interface is not matching
	type notMatchedInterface interface {
		Read(p []byte) (n int, err error, additional []string)
	}
	testStructObj := testStruct{}
	testStructFields, err := util.GetFieldsImplementing(&testStructObj, (*notMatchedInterface)(nil))
	require.NoError(t, err)
	assert.Empty(t, testStructFields)
}

func TestGetFieldsImplementingMemberStructPointer(t *testing.T) {
	content := "runs all day and never walks"
	testStructObj := testStructMemberStructPtr{
		Member: bytes.NewBufferString(content),
	}
	fields, err := util.GetFieldsImplementing(&testStructObj, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Len(t, fields, 1)

	output := make([]byte, len(content))
	n, err := fields[0].Read(output)
	require.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, string(output))
}

func TestGetFieldsImplementingMemberInterface(t *testing.T) {
	content := "it has a bed and never sleeps"
	testStructObj := testStructMemberInterface{
		Member: bytes.NewBufferString(content),
	}
	fields, err := util.GetFieldsImplementing(&testStructObj, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Len(t, fields, 1)

	output := make([]byte, len(content))
	n, err := fields[0].Read(output)
	require.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, string(output))
}

func TestGetFieldsImplementingSuccess(t *testing.T) {
	// Struct not initialized
	// It's a responsibility of a user to make sure that the fields are not nil before using them.
	structNotInitialized := testStruct{}
	structNotInitializedFields, err := util.GetFieldsImplementing(&structNotInitialized, (*readInterface)(nil))
	require.NoError(t, err)
	// There are 4 pointer members of the testStruct satisfying the interface.
	// Nil interface members are not returned.
	assert.Len(t, structNotInitializedFields, 4)
	for _, f := range structNotInitializedFields {
		assert.Nil(t, f)
		assert.Implements(t, (*readInterface)(nil), f)
	}

	// Struct initialized
	testStructObj := testStruct{
		// satisfy only readInterface
		LimitedReader: &io.LimitedReader{N: 100},
		Reader:        strings.NewReader("did you know that"),
		// satisfy both readInterface and writeInterface
		Buffer1:   bytes.NewBufferString("there are rats with"),
		Buffer2:   bytes.NewBufferString("human BRAIN cells transplanted"),
		NetBuffer: &net.Buffers{[]byte{0x19}},
	}

	// Fields implementing readInterface
	readInterfaceFields, err := util.GetFieldsImplementing(&testStructObj, (*readInterface)(nil))
	require.NoError(t, err)
	assert.Len(t, readInterfaceFields, 5)

	for _, f := range readInterfaceFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*readInterface)(nil), f)
	}

	// Fields implementing writeInterface
	writeInterfaceFields, err := util.GetFieldsImplementing(&testStructObj, (*writeInterface)(nil))
	require.NoError(t, err)
	assert.Len(t, writeInterfaceFields, 3)
	for _, f := range writeInterfaceFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*writeInterface)(nil), f)
	}

	type readWriteInterface interface {
		readInterface
		writeInterface
	}
	readWriteInterfaceFields, err := util.GetFieldsImplementing(&testStructObj, (*readWriteInterface)(nil))
	require.NoError(t, err)
	// All members implementing writeInterface implement readInterface too
	assert.Len(t, readWriteInterfaceFields, 3)
}

func TestIsStructInitialized(t *testing.T) {
	tests := []struct {
		name          string
		testStruct    interface{}
		expectError   bool
		errorContains []string
	}{
		// No error cases
		{
			name:        "Empty struct",
			testStruct:  &testStructEmpty{},
			expectError: false,
		},
		{
			name:        "Struct with initialized private field",
			testStruct:  &testStructPrivateFiled{privateMember: bytes.NewBufferString("my content")},
			expectError: false,
		},
		{
			name: "Struct with initialized primitives",
			testStruct: &testStructPrimitives{
				X:    1,
				Y:    "blabla",
				XPtr: new(int),
				YPtr: new(string),
			},
			expectError: false,
		},
		{
			name: "Struct with initialized maps and slices",
			testStruct: &testStructMapsAndSlices{
				MapMember:   make(map[string]int, 0),
				SliceMember: make([]int, 0),
			},
			expectError: false,
		},
		{
			name: "Struct with initialized struct member",
			testStruct: &testStructMemberStruct{
				Member: *bytes.NewBufferString("my content"),
			},
			expectError: false,
		},
		{
			name: "Struct with initialized struct pointer member",
			testStruct: &testStructMemberStructPtr{
				Member: bytes.NewBufferString("my content"),
			},
			expectError: false,
		},
		{
			name: "Struct with initialized interface member",
			testStruct: &testStructMemberInterface{
				Member: bytes.NewBufferString("my content"),
			},
			expectError: false,
		},

		// Error cases
		{
			name:          "Struct with uninitialized private field",
			testStruct:    &testStructPrivateFiled{},
			expectError:   true,
			errorContains: []string{"privateMember"},
		},
		{
			name:          "Struct with uninitialized primitives",
			testStruct:    &testStructPrimitives{},
			expectError:   true,
			errorContains: []string{"X", "Y", "XPtr", "YPtr"},
		},
		{
			name:          "Struct with uninitialized maps and slices",
			testStruct:    &testStructMapsAndSlices{},
			expectError:   true,
			errorContains: []string{"MapMember", "SliceMember"},
		},
		{
			name:          "Struct with uninitialized struct member",
			testStruct:    &testStructMemberStruct{},
			expectError:   true,
			errorContains: []string{"Member"},
		},
		{
			name:          "Struct with uninitialized struct pointer member",
			testStruct:    &testStructMemberStructPtr{},
			expectError:   true,
			errorContains: []string{"Member"},
		},
		{
			name:          "Struct with uninitialized interface member",
			testStruct:    &testStructMemberInterface{},
			expectError:   true,
			errorContains: []string{"Member"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.IsStructInitialized(tt.testStruct)

			if tt.expectError {
				require.Error(t, err)
				for _, errText := range tt.errorContains {
					assert.Contains(t, err.Error(), errText)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type testStructEmpty struct {
}

type testStructPrivateFiled struct {
	privateMember *bytes.Buffer
}

type testStructPrimitives struct {
	X    int
	Y    string
	XPtr *int
	YPtr *string
}

type testStructMapsAndSlices struct {
	MapMember   map[string]int
	SliceMember []int
}

type testStructMemberStruct struct {
	Member bytes.Buffer
}

type testStructMemberStructPtr struct {
	Member *bytes.Buffer
}

type testStructMemberInterface struct {
	Member io.Reader
}
