# `/test`

General global test utilities.

#### `/test/testdata` dir

Place any *snapshots/fixtures* for your test scenarios here.

##### Why?

> Go build ignores directory named testdata.
> The Go tool will ignore any directory in your $GOPATH that starts with a period, an underscore, or matches the word testdata
> When go test runs, it sets current directory as package directory

###### Resources

* https://github.com/golang-standards/project-layout/blob/master/test/README.md
* https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859 
* https://dave.cheney.net/2016/05/10/test-fixtures-in-go