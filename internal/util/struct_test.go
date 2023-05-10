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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")
	// Pointer to a pointer to an interface
	interfaceObjPtr := (*readInterface)(nil)
	_, err = util.GetFieldsImplementing(&testStructEmpty{}, &interfaceObjPtr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")

	// Invalid structPtr input param, must be a pointer to a struct
	_, err = util.GetFieldsImplementing(testStructEmpty{}, (*readInterface)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing((*readInterface)(nil), (*readInterface)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing([]*testStructEmpty{}, (*readInterface)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
}

func TestGetFieldsImplementingNoFields(t *testing.T) {
	// No fields returned from empty structs
	structEmptyFields, err := util.GetFieldsImplementing(&testStructEmpty{}, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structEmptyFields)

	// No fields returned from structs with only private fields
	structPrivate := testStructPrivateFiled{privateMember: bytes.NewBufferString("my content")}
	structPrivateFields, err := util.GetFieldsImplementing(&structPrivate, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structPrivateFields)

	// No fields returned if struct fields are primitive
	structPrimitive := testStructPrimitives{X: 12, Y: "y", XPtr: swag.Int(15), YPtr: swag.String("YPtr")}
	structPrimitiveFields, err := util.GetFieldsImplementing(&structPrimitive, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structPrimitiveFields)

	// No fields returned if struct fields are structs (not pointer to a struct)
	structMemberStruct := testStructMemberStruct{Member: *bytes.NewBufferString("my content")}
	structMemberStructFields, err := util.GetFieldsImplementing(&structMemberStruct, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structMemberStructFields)

	// No fieds returned if an interface is not matching
	type notMatchedInterface interface {
		Read(p []byte) (n int, err error, additional []string)
	}
	testStructObj := testStruct{}
	testStructFields, err := util.GetFieldsImplementing(&testStructObj, (*notMatchedInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, testStructFields)
}

func TestGetFieldsImplementingMemberStructPointer(t *testing.T) {
	content := "runs all day and never walks"
	testStructObj := testStructMemberStructPtr{
		Member: bytes.NewBufferString(content),
	}
	fields, err := util.GetFieldsImplementing(&testStructObj, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Len(t, fields, 1)

	output := make([]byte, len(content))
	n, err := fields[0].Read(output)
	assert.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, string(output))
}

func TestGetFieldsImplementingMemberInterface(t *testing.T) {
	content := "it has a bed and never sleeps"
	testStructObj := testStructMemberInterface{
		Member: bytes.NewBufferString(content),
	}
	fields, err := util.GetFieldsImplementing(&testStructObj, (*readInterface)(nil))
	assert.NoError(t, err)
	assert.Len(t, fields, 1)

	output := make([]byte, len(content))
	n, err := fields[0].Read(output)
	assert.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, string(output))
}

func TestGetFieldsImplementingSuccess(t *testing.T) {
	// Struct not initialized
	// It's a responsibility of a user to make sure that the fields are not nil before using them.
	structNotInitialized := testStruct{}
	structNotInitializedFields, err := util.GetFieldsImplementing(&structNotInitialized, (*readInterface)(nil))
	assert.NoError(t, err)
	// There are 4 pointer members of the testStruct satisfying the interface.
	// Nil interface members are not returned.
	assert.Equal(t, 4, len(structNotInitializedFields))
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
	assert.NoError(t, err)
	assert.Equal(t, 5, len(readInterfaceFields))

	for _, f := range readInterfaceFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*readInterface)(nil), f)
	}

	// Fields implementing writeInterface
	writeInterfaceFields, err := util.GetFieldsImplementing(&testStructObj, (*writeInterface)(nil))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(writeInterfaceFields))
	for _, f := range writeInterfaceFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*writeInterface)(nil), f)
	}

	type readWriteInterface interface {
		readInterface
		writeInterface
	}
	readWriteInterfaceFields, err := util.GetFieldsImplementing(&testStructObj, (*readWriteInterface)(nil))
	assert.NoError(t, err)
	// All members implementing writeInterface implement readInterface too
	assert.Equal(t, 3, len(readWriteInterfaceFields))

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

type testStructMemberStruct struct {
	Member bytes.Buffer
}

type testStructMemberStructPtr struct {
	Member *bytes.Buffer
}

type testStructMemberInterface struct {
	Member io.Reader
}
