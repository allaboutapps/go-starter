package util

// FalseIfNil returns false if the passed pointer is nil. Passing a pointer to a bool will return the value of the bool.
func FalseIfNil(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}
