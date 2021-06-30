package test

import "testing"

// TestingT is used to generate a mock of testing.T to enable testing
// of helper methods which are using assert/require
// Inspired by: https://github.com/uber-go/zap/blob/master/zaptest/testingt_test.go, commit 5b0fd114dcc089875ee61dfad3617c3a43c2e93e
type TestingT interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

// used to ensure compatibility between this interface and testing.TB
var _ TestingT = (testing.TB)(nil)
