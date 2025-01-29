package driver

import (
	"html/template"
	"io"
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"barbecue/core"
)

type HtmlDriver struct {
	core.Driver
    Templates *template.Template
}

func NewHtmlDriver() *HtmlDriver {
	return &HtmlDriver {
	    Templates: template.Must(template.ParseGlob("html/*.html")),
	}
}

const (
	T_INDEX = "index.html"
	T_TASK = "task.html"
	T_FORM = "form.html"
	T_ERROR = "error.html"
)

func (d *HtmlDriver) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return d.Templates.ExecuteTemplate(w, name, data)
}

func (d *HtmlDriver) Ok(ctx echo.Context, code int, template string, data interface{}) error {
	return ctx.Render(code, template, data)
}

func (d *HtmlDriver) Ko(ctx echo.Context, code int) error {
	return ctx.Render(code, T_ERROR, struct { Code int ; Message string } {
		Code: code,
		Message: http.StatusText(code),
	})
}

func (d *HtmlDriver) Redirect(ctx echo.Context, uri string) error {
	if len(uri) == 0 || !strings.HasPrefix(uri, "/") { uri = "/" }
	return ctx.Redirect(http.StatusSeeOther, uri)
	return nil
}
