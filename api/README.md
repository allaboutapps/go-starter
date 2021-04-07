# `/api`

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

Use `make swagger` do generate `/internal/types/**/*.go` from these specs. You may use `make watch-swagger` while writing these schemas to automatically execute `make swagger` every time there are changes (especially useful if you have `http://localhost:8081/` opened in your browser while mutating these specs).

The final `/api/swagger.yml` is auto-generated from the following specs:
1. The skeleton spec at `/api/config/main.yml`,
2. all path specs living within `/api/paths/*.yml` and
3. any `/api/definitions/*.yml` referenced by in step 1 and 2.

Further reading:
* [Swagger v2 specification](https://swagger.io/specification/v2/)
* [go-swagger](https://github.com/go-swagger/go-swagger)
* [go-swagger: Schema Generation Rules](https://github.com/go-swagger/go-swagger/blob/master/docs/use/models/schemas.md)


Related examples regarding this project layout:

* https://github.com/golang-standards/project-layout/blob/master/api/README.md
* https://github.com/kubernetes/kubernetes/tree/master/api
* https://github.com/moby/moby/tree/master/api