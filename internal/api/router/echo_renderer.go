package router

import (
	"fmt"
	"html/template"
	"io"

	"allaboutapps.dev/aw/go-starter/internal/api/router/templates"
	"github.com/labstack/echo/v4"
)

type echoRenderer struct {
	templates map[templates.ViewTemplate]*template.Template
}

func (t *echoRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmplHTML, ok := t.templates[templates.ViewTemplate(name)]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}

	return tmplHTML.Execute(w, data)
}
