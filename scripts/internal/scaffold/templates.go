//go:build scripts

package scaffold

const (
	swaggerDefinitionsTemplate = `swagger: "2.0"
info:
  title: ""
  version: 0.1.0
paths: {}
definitions:
  {{ .Name }}:
    type: object
    required:
    {{- range .Properties}}
      {{- if .Required }}
      - {{ .Name }}
      {{- end }}
    {{- end }}
    properties:
    {{- range .Properties}}
      {{ .Name }}:
        type: {{ .Type }}
        {{- if not .Required }}
        x-nullable: true
        {{- end }}
        {{- if .Format }}
        format: {{ .Format }}
        {{- end }}
    {{- end }}

  {{ .Name }}Payload:
    type: object
    required:
    {{- range .PayloadProperties}}
      {{- if .Required }}
      - {{ .Name }}
      {{- end }}
    {{- end }}
    properties:
    {{- range .PayloadProperties}}
      {{ .Name }}:
        type: {{ .Type }}
        {{- if not .Required }}
        x-nullable: true
        {{- end }}
        {{- if .Format }}
        format: {{ .Format }}
        {{- end }}
    {{- end }}

  {{ .Name }}List:
    type: array
    items: 
      $ref: "#/definitions/{{ .Name }}"
`

	swaggerPathsTemplate = `swagger: "2.0"
info:
  title: ""
  version: 0.1.0
parameters:
  {{ .Name }}IdParam:
    type: string
    format: uuid4
    name: id
    description: ID of {{ .Name }}
    in: path
    required: true

paths:
  /api/v1/{{ .URLName }}s:
    get:
      security:
        - Bearer: []
      description: "Return a list of {{ .Name }}"
      tags:
        - {{ .Name }}
      summary: "Return a list of {{ .Name }}"
      operationId: Get{{ .Name }}ListRoute
      responses:
        "200":
          description: Success
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}List"

    post:
      security:
        - Bearer: []
      description: "Update the given {{ .Name }}"
      tags:
        - {{ .Name }}
      summary: "Update the given {{ .Name }}"
      operationId: Post{{ .Name }}Route
      parameters:
        - name: Payload
          in: body
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}Payload"
      responses:
        "200":
          description: Success
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}"

  /api/v1/{{ .URLName }}s/{id}:
    get:
      security:
        - Bearer: []
      description: "Return {{ .Name }} with ID"
      tags:
        - {{ .Name }}
      summary: "Return {{ .Name }} with ID"
      operationId: Get{{ .Name }}Route
      parameters:
        - $ref: "#/parameters/{{ .Name }}IdParam"
      responses:
        "200":
          description: Success
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}"

    put:
      security:
        - Bearer: []
      description: "Update the given {{ .Name }}"
      tags:
        - {{ .Name }}
      summary: "Update the given {{ .Name }}"
      operationId: Put{{ .Name }}Route
      parameters:
        - $ref: "#/parameters/{{ .Name }}IdParam"
        - name: Payload
          in: body
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}Payload"
      responses:
        "200":
          description: Success
          schema:
            $ref: "../definitions/{{ .URLName }}.yml#/definitions/{{ .Name }}"

    delete:
      security:
        - Bearer: []
      description: "Delete {{ .Name }} with ID"
      tags:
        - {{ .Name }}
      summary: "Delete {{ .Name }} with ID"
      operationId: Delete{{ .Name }}Route
      parameters:
        - $ref: "#/parameters/{{ .Name }}IdParam"
      responses:
        "204":
          description: Success
`

	getHandlerTemplate = `package {{ .Package }}

import (
    "net/http"
    "time"

    "{{ .Module }}/internal/api"
    "{{ .Module }}/internal/types"
    "{{ .Module }}/internal/util"
    "github.com/go-openapi/strfmt"
    "github.com/go-openapi/strfmt/conv"
    "github.com/go-openapi/swag"
    "github.com/labstack/echo/v4"
)

func Get{{ .Resource.Name }}Route(s *api.Server) *echo.Route {
    return s.Router.APIV1{{ .Resource.Name }}.GET("/:id", get{{ .Resource.Name }}Handler(s))
}

func get{{ .Resource.Name }}Handler(s *api.Server) echo.HandlerFunc {
    return func(c echo.Context) error {
        /* Uncomment for real implementation
        ctx := c.Request().Context()

        params := {{ .Package }}.NewGet{{ .Resource.Name }}RouteParams()
        err := util.BindAndValidatePathParams(c, &params)
        if err != nil {
            return err
        }
        id := params.ID.String()

        // TODO: Implement 
        */

        response := types.{{ .Resource.Name }}{
            {{- range .Resource.Fields }}
            {{ .Name }}: {{ .PlaceholderValue }},
            {{- end }}
        }

        return util.ValidateAndReturn(c, http.StatusOK, &response)
    }
}
`

	getListHandlerTemplate = `package {{ .Package }}

import (
    "net/http"
    "time"

    "{{ .Module }}/internal/api"
    "{{ .Module }}/internal/types"
    "{{ .Module }}/internal/util"
    "github.com/go-openapi/strfmt"
    "github.com/go-openapi/strfmt/conv"
    "github.com/go-openapi/swag"
    "github.com/labstack/echo/v4"
)

func Get{{ .Resource.Name }}ListRoute(s *api.Server) *echo.Route {
    return s.Router.APIV1{{ .Resource.Name }}.GET("", get{{ .Resource.Name }}ListHandler(s))
}

func get{{ .Resource.Name }}ListHandler(s *api.Server) echo.HandlerFunc {
    return func(c echo.Context) error {
        /* Uncomment for real implementation
        ctx := c.Request().Context()

        // TODO: Implement 
        */

        item := types.{{ .Resource.Name }}{
            {{- range .Resource.Fields }}
            {{ .Name }}: {{ .PlaceholderValue }},
            {{- end }}
        }
        response := types.{{ .Resource.Name }}List{&item}

        return util.ValidateAndReturn(c, http.StatusOK, &response)
    }
}
`
	postHandlerTemplate = `package {{ .Package }}

import (
    "net/http"
    "time"

    "{{ .Module }}/internal/api"
    "{{ .Module }}/internal/types"
    "{{ .Module }}/internal/util"
    "github.com/go-openapi/strfmt"
    "github.com/go-openapi/strfmt/conv"
    "github.com/go-openapi/swag"
    "github.com/labstack/echo/v4"
)

func Post{{ .Resource.Name }}Route(s *api.Server) *echo.Route {
    return s.Router.APIV1{{ .Resource.Name }}.POST("", post{{ .Resource.Name }}Handler(s))
}

func post{{ .Resource.Name }}Handler(s *api.Server) echo.HandlerFunc {
    return func(c echo.Context) error {
        /* Uncomment for real implementation
        ctx := c.Request().Context()

        var body types.{{ .Resource.Name }}Payload
		    err := util.BindAndValidateBody(c, &body)
		    if err != nil {
		    	  return err
		    }

        // TODO: Implement 
        */

        response := types.{{ .Resource.Name }}{
            {{- range .Resource.Fields }}
            {{ .Name }}: {{ .PlaceholderValue }},
            {{- end }}
        }

        return util.ValidateAndReturn(c, http.StatusOK, &response)
    }
}
`
	putHandlerTemplate = `package {{ .Package }}

import (
    "net/http"
    "time"

    "{{ .Module }}/internal/api"
    "{{ .Module }}/internal/types"
    "{{ .Module }}/internal/util"
    "github.com/go-openapi/strfmt"
    "github.com/go-openapi/strfmt/conv"
    "github.com/go-openapi/swag"
    "github.com/labstack/echo/v4"
)

func Put{{ .Resource.Name }}Route(s *api.Server) *echo.Route {
    return s.Router.APIV1{{ .Resource.Name }}.PUT("/:id", put{{ .Resource.Name }}Handler(s))
}

func put{{ .Resource.Name }}Handler(s *api.Server) echo.HandlerFunc {
    return func(c echo.Context) error {
        /* Uncomment for real implementation
        ctx := c.Request().Context()

        params := {{ .Package }}.NewGet{{ .Resource.Name }}RouteParams()
        err := util.BindAndValidatePathParams(c, &params)
        if err != nil {
            return err
        }
        id := params.ID.String()

        var body types.{{ .Resource.Name }}Payload
		    err = util.BindAndValidateBody(c, &body)
		    if err != nil {
		    	  return err
		    }

        // TODO: Implement 
        */

        response := types.{{ .Resource.Name }}{
            {{- range .Resource.Fields }}
            {{ .Name }}: {{ .PlaceholderValue }},
            {{- end }}
        }

        return util.ValidateAndReturn(c, http.StatusOK, &response)
    }
}
`
	deleteHandlerTemplate = `package {{ .Package }}

import (
    "net/http"

    "{{ .Module }}/internal/api"
    "github.com/labstack/echo/v4"
)

func Delete{{ .Resource.Name }}Route(s *api.Server) *echo.Route {
    return s.Router.APIV1{{ .Resource.Name }}.DELETE("/:id", delete{{ .Resource.Name }}Handler(s))
}

func delete{{ .Resource.Name }}Handler(s *api.Server) echo.HandlerFunc {
    return func(c echo.Context) error {
        /* Uncomment for real implementation
        ctx := c.Request().Context()

        params := {{ .Package }}.NewGet{{ .Resource.Name }}RouteParams()
        err := util.BindAndValidatePathParams(c, &params)
        if err != nil {
            return err
        }
        id := params.ID.String()

        // TODO: Implement 
        */

        return c.NoContent(http.StatusNoContent)
    }
}
`
)
