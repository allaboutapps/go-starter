package pgservetypes

import (
	"testing"

	"github.com/go-openapi/strfmt"
)

func TestHelloWorldType(t *testing.T) {

	x := HelloWorld{
		Hello: "World",
	}

	if err := x.Validate(strfmt.Default); err != nil {
		t.Fatalf("failed: %v", err)
	}

	y := HelloWorld{}

	want := []byte("{\"hello\":\"World\"}")

	if err := y.UnmarshalBinary(want); err != nil {
		t.Fatalf("failed: %v", err)
	}

	if err := y.Validate(strfmt.Default); err != nil {
		t.Fatalf("failed: %v", err)
	}

	got, err := y.MarshalBinary()

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	t.Log(string(got))

	if string(got) != string(want) {
		t.Errorf("got %v, want %v", string(got), string(want))
	}

}
