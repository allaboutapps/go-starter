# `/api`

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

Use `make swagger` do generate `/internal/types/**/*.go` from these specs. You may use `make watch-swagger` while writing these schemas to automatically execute `make swagger` every time there are changes (especially useful if you have `http://localhost:8081/` opened in your browser while mutating these specs).

* [Swagger v2 specification](https://swagger.io/specification/v2/)
* [go-swagger](https://github.com/go-swagger/go-swagger)
* [go-swagger: Schema Generation Rules](https://github.com/go-swagger/go-swagger/blob/master/docs/use/models/schemas.md)


Related examples regarding this project layout:

* https://github.com/golang-standards/project-layout/blob/master/api/README.md
* https://github.com/kubernetes/kubernetes/tree/master/api
* https://github.com/moby/moby/tree/master/api