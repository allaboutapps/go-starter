# `/test`

General global test utilities.

### Regarding `test/test_*.go` and its `test.With*` functions 

Test helpers like `test.WithTestDatabase` and `test.WithTestServer` require usage of a closure, as these functions automatically manage the setup **and** teardown (e.g. server shutdown, db connection drop) for your testcase.

Other pkgs don't have this requirement (e.g. the initialization code for `test.NewTestMailer` which covers the setup for the `mailer` mock), thus, please use this `With*` convention incl. closure **only** when it makes sense.

### Regarding `test/fixtures.go`

This are your global db test fixtures, that are only available while testing. However, feel free to setup specialized fixtures per package if required (e.g. just initialize an additional IntegreSQL template).

### Regarding `test/helper_*.go`

Please use this convention to specify test only utility functions.

### Regarding snapshot testing

Golden files can be created by using the `Snapshoter.Save(t TestingT, data ...interface{})` method. The snapshot can be configured to force an update, use a different replacer function or to set a different file location and suffix for the snaphot.

A snapshot can be updated by either calling the `Update(true)` method, or by using the global override by setting the environment variable `TEST_UPDATE_GOLDEN` to `true`. To update all snapshots the call might look like this: `TEST_UPDATE_GOLDEN=true make test`.

### `testdata` and `.` or `_` prefixed files

Note that Go will ignore directories or files that begin with "." or "_", so you have more flexibility in terms of how you name your test data directory.

> Go build ignores directory named testdata.
> The Go tool will ignore any directory in your $GOPATH that starts with a period, an underscore, or matches the word testdata
> When go test runs, it sets current directory as package directory

* https://github.com/golang-standards/project-layout/blob/master/test/README.md
* https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859 
* https://dave.cheney.net/2016/05/10/test-fixtures-in-go

Examples:
* https://github.com/openshift/origin/tree/master/test (test data is in the `/testdata` subdirectory)

https://github.com/golang-standards/project-layout/tree/master/test